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

	"github.com/lovego/xiaomei/utils/errs"
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
			return body, errs.Stack(err)
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
	return errs.Stack(errors.New(
		`HTTP ` + resp.Request.Method + ` ` + resp.Request.URL.String() + "\n" +
			`Unexpected Response Status: ` + resp.Status,
	))
}

func (resp *Response) Json(data interface{}) error {
	defer resp.Body.Close()
	if data == nil {
		return nil
	}
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	if err := decoder.Decode(&data); err != nil {
		return errs.Stack(err)
	}
	return nil
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
		return errs.Stack(err)
	}
	return nil
}

func makeBodyReader(data interface{}) (reader io.Reader, err error) {
	if data == nil {
		return
	}
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
		if !reflect.ValueOf(body).IsNil() {
			buf, err := json.Marshal(body)
			if err != nil {
				err = errs.Stack(err)
			}
			reader = bytes.NewBuffer(buf)
		}
	}
	return
}
