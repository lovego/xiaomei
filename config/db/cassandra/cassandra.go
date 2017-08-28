package cassandra

import (
	"log"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/lovego/xiaomei/config"
)

var cassandraSessions = struct {
	sync.Mutex
	m map[string]*gocql.Session
}{m: make(map[string]*gocql.Session)}

func Session(name string) *gocql.Session {
	cassandraSessions.Lock()
	defer cassandraSessions.Unlock()
	session := cassandraSessions.m[name]
	if session == nil {
		session = newSession(name)
		cassandraSessions.m[name] = session
	}
	return session
}

func newSession(name string) *gocql.Session {
	conf := config.Get(`cassandra`).Get(name)
	cluster := gocql.NewCluster(conf.GetStringSlice(`hosts`)...)
	cluster.Keyspace = conf.GetString(`keyspace`)
	cluster.PageSize = 0
	cluster.Timeout = time.Second * 3
	cluster.Consistency = gocql.One // LocalOne
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 2}
	cluster.HostFilter = gocql.DataCentreHostFilter(conf.GetString(`dataCenter`))
	sess, err := cluster.CreateSession()
	if err != nil {
		log.Panic(err)
	}
	return sess
}

func Query(cass *gocql.Session, cql string, args ...interface{}) (
	data []map[string]interface{}, columns []gocql.ColumnInfo, err error,
) {
	var iter = cass.Query(cql, args...).Iter()
	data, err = iter.SliceMap()
	if err2 := iter.Close(); err2 != nil && err == nil {
		err = err2
	}
	columns = iter.Columns()
	return
}

func QueryPageSize(cass *gocql.Session, pageSize int, cql string, args ...interface{}) (
	data []map[string]interface{}, columns []gocql.ColumnInfo, err error,
) {
	var iter = cass.Query(cql, args...).PageSize(pageSize).Iter()
	data, err = iter.SliceMap()
	if err2 := iter.Close(); err2 != nil && err == nil {
		err = err2
	}
	columns = iter.Columns()
	return
}
