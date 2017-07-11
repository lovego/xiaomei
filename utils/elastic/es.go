package elastic

import (
	"math"
	"strings"

	"github.com/lovego/xiaomei/utils/httputil"
	"github.com/nu7hatch/gouuid"
)

type ES struct {
	counter   int
	BaseAddrs []string
}

// Ensure that at least one addr must be provided.
func New(addr string, otherAddrs ...string) *ES {
	addrs := []string{}
	addrs = append(addrs, addr)
	addrs = append(addrs, otherAddrs...)
	return &ES{BaseAddrs: addrs}
}

func (es *ES) Get(path string, bodyData, data interface{}) error {
	resp, err := httputil.Get(es.Uri(path), nil, bodyData)
	if err != nil {
		return err
	}
	if err := resp.Ok(); err != nil {
		return err
	}
	return resp.Json(data)
}

func (es *ES) Post(path string, bodyData, data interface{}) error {
	resp, err := httputil.Post(es.Uri(path), nil, bodyData)
	if err != nil {
		return err
	}
	if err := resp.Ok(); err != nil {
		return err
	}
	return resp.Json(data)
}

func (es *ES) Uri(path string) string {
	// Round-Robin elastic nodes
	if es.counter+1 >= math.MaxInt32 {
		es.counter = 0 // reset counter before overflow
	} else {
		es.counter++
	}
	i := es.counter % len(es.BaseAddrs)
	return es.BaseAddrs[i] + path
}

func GenUUID() (string, error) {
	if uid, err := uuid.NewV4(); err != nil {
		return ``, err
	} else {
		return strings.Replace(uid.String(), `-`, ``, -1), nil
	}
}
