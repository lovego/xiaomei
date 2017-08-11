package httputil

import (
	"net/http"
	"time"
)

var DefaultClient = &Client{Client: &http.Client{Timeout: 10 * time.Second}}

func Get(url string, headers map[string]string, body interface{}) (*Response, error) {
	return DefaultClient.Do(http.MethodGet, url, headers, body)
}

func Post(url string, headers map[string]string, body interface{}) (*Response, error) {
	return DefaultClient.Do(http.MethodPost, url, headers, body)
}

func Head(url string, headers map[string]string, body interface{}) (*Response, error) {
	return DefaultClient.Do(http.MethodHead, url, headers, body)
}

func Put(url string, headers map[string]string, body interface{}) (*Response, error) {
	return DefaultClient.Do(http.MethodPut, url, headers, body)
}

func Delete(url string, headers map[string]string, body interface{}) (*Response, error) {
	return DefaultClient.Do(http.MethodDelete, url, headers, body)
}

func GetJson(url string, headers map[string]string, body, data interface{}) error {
	return DefaultClient.DoJson(http.MethodGet, url, headers, body, data)
}

func PostJson(url string, headers map[string]string, body, data interface{}) error {
	return DefaultClient.DoJson(http.MethodPost, url, headers, body, data)
}

func HeadJson(url string, headers map[string]string, body, data interface{}) error {
	return DefaultClient.DoJson(http.MethodHead, url, headers, body, data)
}

func PutJson(url string, headers map[string]string, body, data interface{}) error {
	return DefaultClient.DoJson(http.MethodPut, url, headers, body, data)
}

func DeleteJson(url string, headers map[string]string, body, data interface{}) error {
	return DefaultClient.DoJson(http.MethodDelete, url, headers, body, data)
}
