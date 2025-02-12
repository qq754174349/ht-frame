package result

import "encoding/json"

type Result struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	TraceId string      `json:"traceId"`
}

func NewResult(code int, msg string, traceId string, data interface{}) *Result {
	return &Result{Code: code, Msg: msg, TraceId: traceId, Data: data}
}

func NewBaseSuccessResult(traceId string) *Result {
	return NewTemplateResult(SUCCESS, traceId)
}

func NewSuccessResult(traceId string, data interface{}) *Result {
	return NewResult(SUCCESS.Code, SUCCESS.Msg, traceId, data)
}

func NewFailResult(msg string, traceId string) *Result {
	return NewResult(FAILURE.Code, msg, traceId, nil)
}

func NewTemplateResult(template Template, traceId string) *Result {
	return NewResult(template.Code, template.Msg, traceId, nil)
}

func (Result *Result) ToString() string {
	marshal, _ := json.Marshal(Result)
	return string(marshal)
}
