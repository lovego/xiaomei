package xiaomei

import (
	"encoding/json"
	"log"

	"github.com/lovego/xiaomei/utils/errs"
)

func (res *Response) Json(data interface{}) {
	bytes, err := json.Marshal(data)
	if err == nil {
		res.Header().Set(`Content-Type`, `application/json; charset=utf-8`)
		res.Write(bytes)
	} else {
		log.Panic(err)
	}
}

func (res *Response) Json2(data interface{}, err error) {
	if err != nil {
		res.LogError(err)
	}
	if bytes, err := json.Marshal(data); err == nil {
		res.Header().Set(`Content-Type`, `application/json; charset=utf-8`)
		res.Write(bytes)
	} else {
		panic(err)
	}
}

func (res Response) Message(err error) {
	type result struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err == nil {
		res.Json(result{Code: `ok`, Message: `success`})
	} else {
		if e, ok := err.(errs.CodeMessageErr); ok {
			res.Json(result{Code: e.Code(), Message: e.Message()})
		} else {
			res.LogError(err)
			res.Json(result{Code: `error`, Message: err.Error()})
		}
	}
}

func (res Response) Data(data interface{}, err error) {
	type dataT struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}
	if err == nil {
		res.Json(dataT{Code: `ok`, Message: `success`, Data: data})
	} else {
		if e, ok := err.(errs.CodeMessageErr); ok {
			res.Json(dataT{Code: e.Code(), Message: e.Message(), Data: data})
		} else {
			res.LogError(err)
			res.Json(dataT{Code: `error`, Message: err.Error(), Data: data})
		}
	}
}

func (res Response) Result(data interface{}, err error) {
	type dataT struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"result,omitempty"`
	}
	if err == nil {
		res.Json(dataT{Code: `ok`, Message: `success`, Data: data})
	} else {
		if e, ok := err.(errs.CodeMessageErr); ok {
			res.Json(dataT{Code: e.Code(), Message: e.Message(), Data: data})
		} else {
			res.LogError(err)
			res.Json(dataT{Code: `error`, Message: err.Error(), Data: data})
		}
	}
}

func (res Response) Model(data interface{}, err error) {
	type dataT struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"model,omitempty"`
	}
	if err == nil {
		res.Json(dataT{Code: `ok`, Message: `success`, Data: data})
	} else {
		if e, ok := err.(errs.CodeMessageErr); ok {
			res.Json(dataT{Code: e.Code(), Message: e.Message(), Data: data})
		} else {
			res.LogError(err)
			res.Json(dataT{Code: `error`, Message: err.Error(), Data: data})
		}
	}
}

func (res *Response) LogError(err error) {
	if err == nil {
		return
	}
	log := map[string]interface{}{`err`: err}
	if stack, ok := err.(interface {
		Stack() string
	}); ok {
		log[`stack`] = stack.Stack()
	}
	res.Log(log)
}
