package elastic

import (
	"strings"

	"github.com/lovego/xiaomei/utils/httputil"
	"github.com/nu7hatch/gouuid"
)

type ES struct {
	BaseAddr string
}

func New(addr string) *ES {
	return &ES{BaseAddr: strings.TrimSuffix(addr, `/`)}
}

func (es *ES) Get(path string, bodyData, data interface{}) {
	httputil.Get(es.Uri(path), nil, bodyData).Ok().Json(data)
}

func (es *ES) Post(path string, bodyData, data interface{}) {
	httputil.Post(es.Uri(path), nil, bodyData).Ok().Json(data)
}

func (es *ES) Uri(path string) string {
	if strings.HasPrefix(path, `/`) {
		return es.BaseAddr + path
	} else {
		return es.BaseAddr + `/` + path
	}
}

func GenUUID() string {
	if uid, err := uuid.NewV4(); err != nil {
		panic(err)
	} else {
		return strings.Replace(uid.String(), `-`, ``, -1)
	}
}
