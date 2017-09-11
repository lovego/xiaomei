package registry

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func Tags(imgName string) (tags []string) {
	data := getResponseJson(http.MethodGet, imgName+`/tags/list`)
	if data != nil && data[`tags`] != nil {
		for _, tag := range data[`tags`].([]interface{}) {
			tags = append(tags, tag.(string))
		}
	}
	return
}

func Remove(imgName, digest string) {
	resp := getResponse(http.MethodDelete, imgName+`/manifests/`+digest)
	if resp.StatusCode != http.StatusAccepted {
		content, err := ioutil.ReadAll(resp.Body)
		log.Panicf("unexpected response: %s\n%s\n%v", resp.Status, content, err)
	}
}

func Digest(imgName, tag string) string {
	resp := getResponse(http.MethodHead, imgName+`/manifests/`+tag)
	digest := resp.Header.Get(`Docker-Content-Digest`)
	if digest == `` {
		log.Panicf("get image digest faild for: %s:%s ", imgName, tag)
	}
	return digest
}

var httpClient = http.Client{Timeout: 5 * time.Second}

func getResponse(method, resource string) *http.Response {
	// TODO: https or http check.
	// TODO: https://registry.hub.docker.com/v2/
	uri, err := url.Parse(`http://` + resource)
	if err != nil {
		log.Panic(err)
	}
	uri.Path = `/v2` + uri.Path
	req, err := http.NewRequest(method, uri.String(), nil)
	if err != nil {
		log.Panic(err)
	}
	// schemaVersion: 2
	req.Header.Set(`Accept`, `application/vnd.docker.distribution.manifest.v2+json`)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Panic(err)
	}
	return resp
}

func getResponseJson(method, resource string) map[string]interface{} {
	resp := getResponse(method, resource)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(body, &m); err != nil {
		log.Panic(err)
	}
	return m
}
