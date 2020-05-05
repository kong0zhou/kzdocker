package log

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"kzdocker/base"
	"time"

	// postgres
	"github.com/lib/pq"
)

type logItem struct {
	CreateTime string `json:"T"`
	Level      string `json:"L"`
	Name       string `json:"N"`
	Caller     string `json:"C"`
	Message    string `json:"M"`
	Stacktrace string `json:"S"`
	UUID       string `json:"uuid"`
	GoID       int64  `json:"goid"`
}

type pgWriter struct {
	db  *sql.DB
	buf [][]byte

	preSyncTime  time.Time     //上一次同步时间
	syncInterval time.Duration //同步间隔
}

func newpgWriter() (core *pgWriter, err error) {
	if base.Config.ZapLog.Postgres.Dbname == `` ||
		base.Config.ZapLog.Postgres.Host == `` ||
		base.Config.ZapLog.Postgres.Password == `` ||
		base.Config.ZapLog.Postgres.Port == `` ||
		base.Config.ZapLog.Postgres.Schemas == `` ||
		base.Config.ZapLog.Postgres.Sslmode == `` ||
		base.Config.ZapLog.Postgres.Tablename == `` ||
		base.Config.ZapLog.Postgres.User == `` {
		err = fmt.Errorf(`newpgWriter() :配置文件的信息不齐全，请配置完所有信息再启动程序`)
		return nil, err
	}
	dbInfo := fmt.Sprintf(`host=%s port=%s user=%s dbname=%s sslmode=%s password=%s`,
		base.Config.ZapLog.Postgres.Host,
		base.Config.ZapLog.Postgres.Port,
		base.Config.ZapLog.Postgres.User,
		base.Config.ZapLog.Postgres.Dbname,
		base.Config.ZapLog.Postgres.Sslmode,
		base.Config.ZapLog.Postgres.Password)
	db, err := sql.Open(`postgres`, dbInfo)
	if err != nil {
		fmt.Println(`newpgWriter() error:`, err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		fmt.Println(`newpgWriter() error:`, err)
		return nil, err
	}
	stmt, err := db.Prepare(fmt.Sprintf(`create schema if not exists %s`,
		base.Config.ZapLog.Postgres.Schemas))
	if err != nil {
		fmt.Println(`newpgWriter() error:`, err)
		return nil, err
	}
	defer stmt.Close()
	result, err := stmt.Exec()
	if err != nil {
		fmt.Println(`newpgWriter() error:`, err)
		return nil, err
	}
	fmt.Println(`postgres result:`, result)
	stmt2, err := db.Prepare(fmt.Sprintf(`
		create table if not exists %s.%s(
			id SERIAL primary key,
			create_time varchar,
			level varchar,
			name varchar,
			caller varchar,
			message varchar,
			stacktrace varchar,
			uuid varchar,
			goid int8
		)
	`, base.Config.ZapLog.Postgres.Schemas, base.Config.ZapLog.Postgres.Tablename))
	if err != nil {
		fmt.Println(`newpgWriter() error:`, err)
		return nil, err
	}
	result, err = stmt2.Exec()
	if err != nil {
		fmt.Println(`newpgWriter() error:`, err)
		return nil, err
	}
	fmt.Println(`postgres result:`, result)
	return &pgWriter{
		db:           db,
		buf:          make([][]byte, 0),
		preSyncTime:  time.Now(),
		syncInterval: 15 * time.Second,
	}, nil
}

func (t *pgWriter) Write(d []byte) (n int, err error) {
	if t == nil {
		err = fmt.Errorf(`t *pgWriter is null`)
		fmt.Println(`Write() error: `, err)
		return 0, err
	}
	if d == nil || len(d) == 0 {
		err = fmt.Errorf(`d is null or empty`)
		fmt.Println(`Write() error: `, err)
		return 0, err
	}
	n = len(d)
	err = t.storebuf(d)
	if err != nil {
		return n, err
	}
	t.sync(false)
	if err != nil {
		return n, err
	}
	return n, nil
}

func (t *pgWriter) Sync() (err error) {
	err = t.sync(true)
	if err != nil {
		return err
	}
	return nil
}

func (t *pgWriter) Close() {
	t.db.Close()
}

func (t *pgWriter) storebuf(d []byte) (err error) {
	if t == nil {
		err = fmt.Errorf(`t *pgWriter is null`)
		fmt.Println(`storebuf() `, err)
		return err
	}
	if d == nil || len(d) == 0 {
		return nil
	}
	// fmt.Println(`storebuf info:`, `存储日志缓存`)
	data := make([]byte, len(d), len(d))
	copy(data, d)
	t.buf = append(t.buf, data)
	return nil
}

func (t *pgWriter) sync(force bool) (err error) {
	if t == nil {
		err = fmt.Errorf(`t *pgWriter is null`)
		fmt.Println(`storebuf() error:`, err)
		return err
	}
	if t.buf == nil || len(t.buf) == 0 {
		// fmt.Println(`sync() info:`, `buf is null`)
		return nil
	}
	if !force && time.Now().Sub(t.preSyncTime) < t.syncInterval {
		// fmt.Println(`sync() info:`, `时间未到`)
		return nil
	}

	// fmt.Println(`sync() info:`, `开始执行`)
	// 开始事务
	tx, err := t.db.Begin()
	if err != nil {
		fmt.Println(`(t *pgWriter) sync() error:`, err)
		return err
	}

	// copy命令
	stmt, err := tx.Prepare(pq.CopyInSchema(base.Config.ZapLog.Postgres.Schemas, base.Config.ZapLog.Postgres.Tablename,
		`create_time`, `level`, `name`, `caller`, `message`, `stacktrace`, `uuid`, `goid`))
	if err != nil {
		fmt.Println(`(t *pgWriter) sync() error:`, err)
		return err
	}
	defer stmt.Close()

	for _, v := range t.buf {
		var l logItem
		err = json.Unmarshal(v, &l)
		if err != nil {
			fmt.Println(`(t *pgWriter) sync() error:`, err)
			return err
		}
		_, err = stmt.Exec(l.CreateTime, l.Level, l.Name, l.Caller, l.Message, l.Stacktrace, l.UUID, l.GoID)
		if err != nil {
			fmt.Println(`(t *pgWriter) sync() error:`, err)
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(`(t *pgWriter) sync() error:`, err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println(`(t *pgWriter) sync() error:`, err)
		return err
	}
	t.buf = t.buf[:0]
	t.preSyncTime = time.Now()
	return nil
}
