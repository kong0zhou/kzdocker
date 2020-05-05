package log

import (
	"fmt"
	"kzdocker/base"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// initLogger 初始化logger
func initLogger() (err error) {
	if log != nil {
		return nil
	}
	timeFormat := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(`2006-01-02 15:04:05.000`))
	}
	// --------------Encoder config--------------
	textConfig := zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     timeFormat,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	jsonConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeFormat,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	// ------------ dst -------------
	dst := make([]zapcore.Core, 0)
	// ------------ console ------------
	if base.Config.ZapLog.Console.Enable {
		consoleEncoder := zapcore.NewConsoleEncoder(textConfig)
		consoleWriteSyncer := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(consoleEncoder, consoleWriteSyncer, zapcore.DebugLevel)
		dst = append(dst, core)
	}
	// ------------file------------
	if base.Config.ZapLog.File.Enable {
		if base.Config.ZapLog.File.RelPath == `` {
			err = fmt.Errorf(`Config.ZapLog.File.RelPath is empty`)
			return err
		}
		logPath := base.BasePath + base.Config.ZapLog.File.RelPath
		dir := filepath.Dir(logPath)
		fmt.Println(`log dir:`, dir)
		if dir != `.` && !isPathExist(dir) {
			fmt.Println(`mkdirAll`, dir)
			os.MkdirAll(dir, 0777)
		}
		file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println(`create file err:`, err)
			return err
		}
		fileEncoder := zapcore.NewJSONEncoder(jsonConfig)
		fileWriteSyncer := zapcore.Lock(file)
		core := zapcore.NewCore(fileEncoder, fileWriteSyncer, zapcore.InfoLevel)
		dst = append(dst, core)
	}
	// -------------postgres--------------
	if base.Config.ZapLog.Postgres.Enable {
		pgEncoder := zapcore.NewJSONEncoder(jsonConfig)
		writer, err := newpgWriter()
		if err != nil {
			return err
		}
		pgWriteSyncer := zapcore.Lock(writer)
		core := zapcore.NewCore(pgEncoder, pgWriteSyncer, zapcore.InfoLevel)
		dst = append(dst, core)
	}
	// -------------------------------
	if len(dst) == 0 {
		err = fmt.Errorf(`zap core is empty`)
		return err
	}
	allCore := zapcore.NewTee(dst...)
	log = zap.New(allCore,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		// 堆栈跟踪
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSampler(core, time.Second, 100, 100)
		}),
		zap.Fields(zap.String(`uuid`, base.AppUUID)),
	)
	defer log.Sync()
	fmt.Println(fmt.Sprintf(`***********************************
******     version: %s     ******
***********************************`, base.AppVersion))
	return nil
}
