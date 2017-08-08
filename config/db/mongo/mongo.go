package mongo

import (
	"sync"

	"github.com/lovego/xiaomei/config"
	"gopkg.in/mgo.v2"
)

var mongoSessions = struct {
	m map[string]MongoSession
	sync.Mutex
}{
	m: make(map[string]MongoSession),
}

func Session(name string) MongoSession {
	mongoSessions.Lock()
	defer mongoSessions.Unlock()
	sess, ok := mongoSessions.m[name]
	if !ok {
		session, err := mgo.Dial(config.DataSource(`mongo`, name))
		if err != nil {
			panic(err)
		}
		sess = MongoSession{session}
		mongoSessions.m[name] = sess
	}
	return sess
}

type MongoSession struct {
	s *mgo.Session
}
type MongoDB struct {
	db *mgo.Database
}
type MongoColl struct {
	c *mgo.Collection
}

func (s MongoSession) Session(work func(*mgo.Session)) {
	sess := s.s.Copy()
	defer sess.Close()
	work(sess)
}

func (db MongoDB) Session(work func(*mgo.Database)) {
	sess := db.db.Session.Copy()
	defer sess.Close()
	work(db.db.With(sess))
}

func (c MongoColl) Session(work func(*mgo.Collection)) {
	sess := c.c.Database.Session.Copy()
	defer sess.Close()
	work(c.c.With(sess))
}

func (s MongoSession) DB(name string) MongoDB {
	return MongoDB{s.s.DB(name)}
}

func (db MongoDB) C(name string) MongoColl {
	return MongoColl{db.db.C(name)}
}
