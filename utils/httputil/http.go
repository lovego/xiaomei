package httputil

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

func Get(url string, headers map[string]string, body interface{}) *Response {
	return Do(http.MethodGet, url, headers, body)
}

func Post(url string, headers map[string]string, body interface{}) *Response {
	return Do(http.MethodPost, url, headers, body)
}

func Head(url string, headers map[string]string, body interface{}) *Response {
	return Do(http.MethodHead, url, headers, body)
}

func Put(url string, headers map[string]string, body interface{}) *Response {
	return Do(http.MethodPut, url, headers, body)
}

func Delete(url string, headers map[string]string, body interface{}) *Response {
	return Do(http.MethodDelete, url, headers, body)
}

func Do(method, url string, headers map[string]string, body interface{}) *Response {
	req, err := http.NewRequest(method, url, makeBodyReader(body))
	if err != nil {
		panic(err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if resp, err := http.DefaultClient.Do(req); err != nil {
		panic(err)
	} else {
		return &Response{Response: resp}
	}
}

type Response struct {
	*http.Response
	body []byte
}

func (resp *Response) GetBody() []byte {
	if resp.body == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(`HTTP ` + resp.Request.Method + `: ` + resp.Request.URL.String() + "\n" +
				`Read Body: ` + err.Error(),
			)
		}
		resp.body = body
	}
	return resp.body
}

func (resp *Response) Ok() *Response {
	if resp.StatusCode != http.StatusOK {
		panic(`HTTP ` + resp.Request.Method + `: ` + resp.Request.URL.String() + "\n" +
			`Response Status: ` + resp.Status + "\n" + string(resp.GetBody()),
		)
	}
	return resp
}

func (resp *Response) Json(data interface{}) {
	if data == nil {
		resp.Body.Close()
		return
	}
	if err := json.Unmarshal(resp.GetBody(), &data); err != nil {
		panic(err)
	}
}

func makeBodyReader(data interface{}) (reader io.Reader) {
	if data == nil {
		return
	}
	switch body := data.(type) {
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
				panic(err)
			}
			reader = bytes.NewBuffer(buf)
		}
	}
	return
}
