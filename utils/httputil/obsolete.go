package httputil

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func Http(method, url string, headers map[string]string, body io.Reader, data interface{}) []byte {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		panic(err)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(`HTTP POST: ` + url + "\n" + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		panic(`HTTP ` + method + `: ` + url + "\n" + `Response Status: ` + resp.Status + "\n" + string(content))
	}

	if err := json.Unmarshal(content, &data); err != nil {
		panic(err)
	}
	return content
}

func HttpStatus(method, url string, headers map[string]string, body io.Reader) (int, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}
