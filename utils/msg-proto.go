package utils

/**
描述：与前端的数据协议
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kzdocker/log"
	"net/http"
)

// ReplyProto 后端响应数据通信协议
type ReplyProto struct {
	Status   int         `json:"status"` //状态 0正常，小于0出错，大于0可能有问题
	Msg      string      `json:"msg"`    //状态信息
	Data     interface{} `json:"data"`
	API      string      `json:"API"`      //api接口
	Method   string      `json:"method"`   //post,put,get,delete
	SN       string      `json:"SN"`       //标识符
	RowCount int         `json:"rowCount"` //Data若是数组，算其长度

	write http.ResponseWriter
	sse   *Sse
}

// NewReplyProto 工厂函数，新建一个ReplyProto对象
func NewReplyProto(r *http.Request, w http.ResponseWriter) (*ReplyProto, error) {
	if r == nil {
		err := fmt.Errorf(`*http.Request is nil`)
		log.Error(err.Error())
		return nil, err
	}
	if w == nil {
		err := fmt.Errorf(`*http.ResponseWriter is nil`)
		log.Error(err.Error())
		return nil, err
	}
	t := ReplyProto{}
	t.Method = r.Method
	t.API = r.RequestURI
	t.write = w
	return &t, nil
}

// NewSseReplyProto 工厂函数，新建一个ReplyProto对象(基于sse)
func NewSseReplyProto(r *http.Request, sse *Sse) (*ReplyProto, error) {
	if r == nil {
		err := fmt.Errorf(`*http.Request is nil`)
		log.Error(err.Error())
		return nil, err
	}
	if sse == nil {
		err := fmt.Errorf(`*Sse is nil`)
		log.Error(err.Error())
		return nil, err
	}
	t := ReplyProto{}
	t.Method = r.Method
	t.API = r.RequestURI
	t.sse = sse
	return &t, nil
}

// ErrorResp 返回错误信息，status为-1
func (t *ReplyProto) ErrorResp(errMsg string) (err error) {
	if t == nil {
		err = fmt.Errorf(`replyProto is null`)
		log.Error(err.Error())
		return
	}
	if t.write == nil {
		err = fmt.Errorf("arguments can not be a nil value")
		log.Error(err.Error())
		return err
	}
	t.Status = -1
	t.Msg = errMsg
	t.RowCount = 0
	response, err := json.Marshal(&t)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	_, err = t.write.Write(response)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// SuccessResp 正常返回数据，status 为 0
func (t *ReplyProto) SuccessResp(data interface{}) (err error) {
	if t == nil {
		err = fmt.Errorf(`replyProto is null`)
		log.Error(err.Error())
		return
	}
	if t.write == nil {
		err = fmt.Errorf("arguments can not be a nil value")
		log.Error(err.Error())
		return err
	}
	t.Status = 0
	t.Data = data
	t.RowCount = 0
	response, err := json.Marshal(&t)
	if err != nil {
		return err
	}
	_, err = t.write.Write(response)
	if err != nil {
		return err
	}
	return nil
}

// DefinedResp 自定义返回数据
func (t *ReplyProto) DefinedResp(status int, msg string, data interface{}, SN string, rowCount int) (err error) {
	if t.write == nil {
		err = fmt.Errorf("arguments can not be a nil value")
		log.Error(err.Error())
		return err
	}
	if t == nil {
		err = fmt.Errorf(`replyProto is null`)
		log.Error(err.Error())
		return
	}
	t.Status = status
	t.Msg = msg
	t.Data = data
	t.RowCount = rowCount
	t.SN = SN
	response, err := json.Marshal(&t)
	if err != nil {
		return err
	}
	_, err = t.write.Write(response)
	if err != nil {
		return err
	}
	return nil
}

// SseError sse错误数据发送
func (t *ReplyProto) SseError(errMsg string) (err error) {
	if t == nil {
		err = fmt.Errorf(`replyProto is null`)
		log.Error(err.Error())
		return
	}
	if t.sse == nil {
		err = fmt.Errorf("sse is nil")
		log.Error(err.Error())
		return
	}
	t.Status = -1
	t.Msg = errMsg
	t.RowCount = 0
	response, err := json.Marshal(&t)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	err = t.sse.Write(SseData{Event: "error", Data: string(response)})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// SseSuccess sse成功发送
func (t *ReplyProto) SseSuccess(data interface{}) (err error) {
	if t == nil {
		err = fmt.Errorf(`replyProto is null`)
		log.Error(err.Error())
		return
	}

	if t.sse == nil {
		err = fmt.Errorf("sse is nil")
		log.Error(err.Error())
		return
	}
	t.Status = 0
	t.Data = data
	t.RowCount = 0
	response, err := json.Marshal(&t)
	err = t.sse.Write(SseData{Event: "message", Data: string(response)})
	if err != nil {
		return err
	}
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

//ReqProto 前端请求数据通讯协议
type ReqProto struct {
	Action   string          `json:"action"` //请求类型GET/POST/PUT/DELETE
	Data     json.RawMessage `json:"data"`   //请求数据
	Sets     []string        `json:"sets"`
	OrderBy  string          `json:"orderBy"`  //排序要求
	Filter   string          `json:"filter"`   //筛选条件
	Page     int             `json:"page"`     //分页
	PageSize int             `json:"pageSize"` //分页大小
}

//DecodeBody 解析r.Body里面的数据，前提是使用json格式，前端使用相同的数据协议
func DecodeBody(r *http.Request) (*ReqProto, error) {
	var err error
	if r == nil {
		err = fmt.Errorf(`Parameter is nil`)
		log.Error(err.Error())
		return nil, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var req ReqProto
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return &req, nil
}

//DecodeURL 解析url中的数据
func DecodeURL(r *http.Request, key string) (*ReqProto, error) {
	var err error
	if r == nil {
		err = fmt.Errorf(`request is nil`)
		log.Error(err.Error())
		return nil, err
	}
	if key == `` {
		err = fmt.Errorf(`key is empty`)
		log.Error(err.Error())
		return nil, err
	}
	data := r.URL.Query().Get(key)
	// logs.Info(data)
	// var filePath string
	var req ReqProto
	err = json.Unmarshal([]byte(data), &req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return &req, nil
}
