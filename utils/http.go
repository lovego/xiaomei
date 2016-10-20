package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func HttpGet(url string, data interface{}) []byte {
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic(`HTTP GET: ` + url + "\n" + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(`HTTP GET: ` + url + "\n" + err.Error())
	}

	if err := json.Unmarshal(body, data); err != nil {
		panic(err)
	}
	return body
}

func HttpPost(url, typ string, body []byte, data interface{}) []byte {
	resp, err := http.Post(url, typ, bytes.NewBuffer(body))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic(`HTTP POST: ` + url + "\n" + resp.Status)
	}

	res_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(`HTTP POST: ` + url + "\n" + err.Error())
	}

	if err := json.Unmarshal(res_body, &data); err != nil {
		panic(err)
	}
	return body
}
