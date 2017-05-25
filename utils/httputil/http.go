package httputil

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func Do(method, url string, headers map[string]string, body io.Reader) *Response {
	req, err := http.NewRequest(method, url, body)
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
	if err := json.Unmarshal(resp.GetBody(), &data); err != nil {
		panic(err)
	}
}
