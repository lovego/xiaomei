package elastic

import (
	"log"
	"strings"

	"github.com/lovego/xiaomei/utils/httputil"
	"github.com/nu7hatch/gouuid"
)

type ES struct {
	BaseAddrs []string
	i         int
	client    *httputil.Client
}

func New(addrs ...string) *ES {
	if len(addrs) == 0 {
		log.Panic(`empty elastic addrs`)
	}
	return &ES{BaseAddrs: addrs, client: httputil.DefaultClient}
}

func New2(client *httputil.Client, addrs ...string) *ES {
	if len(addrs) == 0 {
		log.Panic(`empty elastic addrs`)
	}
	return &ES{BaseAddrs: addrs, client: client}
}

func (es *ES) Get(path string, bodyData, data interface{}) error {
	resp, err := es.client.Get(es.Uri(path), nil, bodyData)
	if err != nil {
		return err
	}
	if err := resp.Ok(); err != nil {
		return err
	}
	return resp.Json(data)
}

func (es *ES) Post(path string, bodyData, data interface{}) error {
	resp, err := es.client.Post(es.Uri(path), nil, bodyData)
	if err != nil {
		return err
	}
	if err := resp.Ok(); err != nil {
		return err
	}
	return resp.Json(data)
}

func (es *ES) Uri(path string) string {
	uri := es.BaseAddrs[es.i] + path
	if len(es.BaseAddrs) > 1 { // Round-Robin elastic nodes
		es.i++
		if es.i >= len(es.BaseAddrs) {
			es.i = 0
		}
	}
	return uri
}

func GenUUID() (string, error) {
	if uid, err := uuid.NewV4(); err != nil {
		return ``, err
	} else {
		return strings.Replace(uid.String(), `-`, ``, -1), nil
	}
}
