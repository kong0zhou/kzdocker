package utils

/**
用途： 用于controllers sse
描述：
	两个结构体：
	SseData: 用于构建sse发送的数据，直接填充即可
	Sse:sse发送器，
		暴露的属性：IsClose->客户端关闭时会发送一个信号
		暴露的函数：UpdateSse：将http.ResponseWriter升级为sse
				  Write :发送数据
*/

import (
	"fmt"
	"iron/log"
	"net/http"
	"time"
)

// SseData 填充sse发送的数据
type SseData struct {
	Event string
	ID    int64
	Retry uint
	Data  string
}

func (s SseData) convertText() (text string, err error) {
	sseText := "id:%d\nevent:%s\nretry:%d\ndata:%s\n\n"
	text = fmt.Sprintf(sseText, s.ID, s.Event, s.Retry, s.Data)
	return text, nil
}

// Sse 实现sse协议
type Sse struct {
	IsClosed  <-chan bool
	f         http.Flusher
	w         http.ResponseWriter
	touchTime time.Time
}

// UpdateSse 将http.ResponseWriter升级为sse
func UpdateSse(w http.ResponseWriter) (s *Sse, err error) {
	if w == nil {
		err = fmt.Errorf(`w is null`)
		log.Error(err.Error())
		return nil, err
	}
	fl, ok := w.(http.Flusher)
	if !ok {
		err = fmt.Errorf(`sse is unsupported`)
		log.Error(err.Error())
		return nil, err
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	w.Header().Set("Transfer-Encoding", "chunked")

	// 必须的执行这一步，sse链接是从服务端发送第一条信息之后才开始建立的
	fmt.Fprintf(w, "event:isconnected\ndata: connection is established\n\n")
	fl.Flush()

	// 将一个 sse close的信号量变成两个，一个用于停止sse心跳包
	isclosed := make(chan bool, 1)
	tickerClose := make(chan bool, 1)
	go func() {
		<-w.(http.CloseNotifier).CloseNotify()
		isclosed <- true
		tickerClose <- true
	}()

	sse := &Sse{
		// IsClosed: w.(http.CloseNotifier).CloseNotify(),
		IsClosed:  isclosed,
		f:         fl,
		w:         w,
		touchTime: time.Now(),
	}

	//心跳包
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for {
			select {
			case <-tickerClose:
				log.Info(`sse通道已经关闭，心跳功能停止`)
				return
			case <-ticker.C:
				now := time.Now()
				nop := "event:nop\ndata:nop\n\n"
				if now.Sub(sse.touchTime) < 30*time.Second {
					continue
				}
				log.Info("push keep alive msg")
				fmt.Fprint(w, nop)
				fl.Flush()
			}
		}
	}()

	return sse, nil
}

func (s *Sse) Write(data SseData) (err error) {
	if data.Data == `` {
		log.Warn(`see: data is null or empty`)
	}
	select {
	case <-s.IsClosed:
		err = fmt.Errorf(`sse is closed`)
		log.Error(err.Error())
		return err
	default:
		sseText, err := data.convertText()
		if err != nil {
			log.Error(err.Error())
			return err
		}
		_, err = fmt.Fprint(s.w, sseText)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		s.f.Flush()
		return nil
	}
}
