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
		if erro, ok := err.(interface {
			Code() string
			Message() string
		}); ok && erro.Code() != "" {
			result[`code`] = erro.Code()
			result[`message`] = erro.Message()
			if e, ok := err.(interface {
				Err() error
			}); ok && e.Err() != nil {
				res.LogError(err)
			}
		} else {
			res.WriteHeader(500)
			result[`code`] = `server-err`
			result[`message`] = `Server Error.`
			res.LogError(err)
		}
	}

	if data != nil {
		result[key] = data
	} else if err != nil {
		if erro, ok := err.(interface {
			Data() interface{}
		}); ok && erro.Data() != nil {
			result[key] = erro.Data()
		}
	}
	res.Json(result)
}

func (res *Response) LogError(err error) {
	res.err = err
}

func (res *Response) GetError() error {
	return res.err
}
