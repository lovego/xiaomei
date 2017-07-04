package httputil

import (
	"net/http"
)

func Get(url string, headers map[string]string, body interface{}) (*Response, error) {
	return Do(http.MethodGet, url, headers, body)
}

func Post(url string, headers map[string]string, body interface{}) (*Response, error) {
	return Do(http.MethodPost, url, headers, body)
}

func Head(url string, headers map[string]string, body interface{}) (*Response, error) {
	return Do(http.MethodHead, url, headers, body)
}

func Put(url string, headers map[string]string, body interface{}) (*Response, error) {
	return Do(http.MethodPut, url, headers, body)
}

func Delete(url string, headers map[string]string, body interface{}) (*Response, error) {
	return Do(http.MethodDelete, url, headers, body)
}

func GetJson(url string, headers map[string]string, body, data interface{}) error {
	return DoJson(http.MethodGet, url, headers, body, data)
}

func PostJson(url string, headers map[string]string, body, data interface{}) error {
	return DoJson(http.MethodPost, url, headers, body, data)
}

func HeadJson(url string, headers map[string]string, body, data interface{}) error {
	return DoJson(http.MethodHead, url, headers, body, data)
}

func PutJson(url string, headers map[string]string, body, data interface{}) error {
	return DoJson(http.MethodPut, url, headers, body, data)
}

func DeleteJson(url string, headers map[string]string, body, data interface{}) error {
	return DoJson(http.MethodDelete, url, headers, body, data)
}

func Do(method, url string, headers map[string]string, body interface{}) (*Response, error) {
	bodyReader, err := makeBodyReader(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if resp, err := http.DefaultClient.Do(req); err != nil {
		return nil, err
	} else {
		return &Response{Response: resp}, nil
	}
}

func DoJson(method, url string, headers map[string]string, body, data interface{}) error {
	resp, err := Do(method, url, headers, body)
	if err != nil {
		return err
	}
	if err := resp.Ok(); err != nil {
		return err
	}
	return resp.Json(data)
}
