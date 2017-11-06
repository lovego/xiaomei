package httputil

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type Response struct {
	*http.Response
	body []byte
}

func (resp *Response) GetBody() ([]byte, error) {
	if resp.body == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		resp.body = body
		if err != nil {
			return body, err
		}
	}
	return resp.body, nil
}

func (resp *Response) Ok() error {
	if resp.StatusCode != http.StatusOK {
		return resp.CodeError()
	}
	return nil
}

func (resp *Response) Check(codes ...int) error {
	for _, code := range codes {
		if resp.StatusCode == code {
			return nil
		}
	}
	return resp.CodeError()
}

func (resp *Response) CodeError() error {
	msg := `HTTP ` + resp.Request.Method + ` ` + resp.Request.URL.String() + "\n" +
		`Unexpected Response: ` + resp.Status
	if body, err := resp.GetBody(); err == nil {
		msg += "\n" + string(body)
	} else {
		msg += "\n(GetBody error: " + err.Error() + `)`
	}
	return errors.New(msg)
}

func (resp *Response) Json(data interface{}) error {
	if data == nil {
		defer resp.Body.Close()
		return nil
	}
	body, err := resp.GetBody()
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewBuffer(body))
	decoder.UseNumber()
	return decoder.Decode(&data)
}

func (resp *Response) Json2(data interface{}) error {
	if data == nil {
		resp.Body.Close()
		return nil
	}
	body, err := resp.GetBody()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}
	return nil
}

func makeBodyReader(data interface{}) (io.Reader, error) {
	if data == nil {
		return nil, nil
	}
	var reader io.Reader
	switch body := data.(type) {
	case io.Reader:
		reader = body
	case string:
		if len(body) > 0 {
			reader = strings.NewReader(body)
		}
	case []byte:
		if len(body) > 0 {
			reader = bytes.NewBuffer(body)
		}
	default:
		if !isNil(body) {
			buf, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			reader = bytes.NewBuffer(buf)
		}
	}
	return reader, nil
}

func isNil(data interface{}) bool {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		return false
	}
}
