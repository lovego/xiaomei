package xiaomei

import (
	"encoding/json"
	"log"
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

func (res *Response) Data(data interface{}, err error) {
	res.DataWithKey(data, err, `data`)
}

func (res *Response) Result(data interface{}, err error) {
	res.DataWithKey(data, err, `result`)
}

func (res *Response) DataWithKey(data interface{}, err error, key string) {
	result := make(map[string]interface{})
	if err == nil {
		result[`code`] = `ok`
		result[`message`] = `success`
	} else {
		if e, ok := err.(interface {
			Code() string
		}); ok {
			if e2, ok := err.(interface {
				LogError() bool
			}); ok && e2.LogError() {
				res.LogError(err)
			}
			result[`code`] = e.Code()
		} else {
			res.WriteHeader(500)
			res.LogError(err)
			result[`code`] = `server-err`
		}

		if e, ok := err.(interface {
			Message() string
		}); ok {
			result[`message`] = e.Message()
		} else {
			result[`message`] = `Server Error.`
		}
	}

	if data != nil {
		result[key] = data
	} else if err != nil {
		if datas, ok := err.(interface {
			Data() interface{}
		}); ok {
			if data = datas.Data(); data != nil {
				result[key] = data
			}
		}
	}
	res.Json(result)
}

func (res *Response) LogError(err error) {
	if err == nil {
		return
	}
	res.Request.Error = err.Error()
	if stack, ok := err.(interface {
		Stack() string
	}); ok {
		res.Request.Stack = stack.Stack()
	}
}
