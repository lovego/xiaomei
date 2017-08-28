package mongo

import (
	"log"
	"sync"

	"github.com/lovego/xiaomei/config"
	"gopkg.in/mgo.v2"
)

var mongoSessions = struct {
	m map[string]Sess
	sync.Mutex
}{
	m: make(map[string]Sess),
}

func Session(name string) Sess {
	mongoSessions.Lock()
	defer mongoSessions.Unlock()
	sess, ok := mongoSessions.m[name]
	if !ok {
		session, err := mgo.Dial(config.Get(`mongo`).GetString(name))
		if err != nil {
			log.Panic(err)
		}
		sess = Sess{session}
		mongoSessions.m[name] = sess
	}
	return sess
}

type Sess struct {
	s *mgo.Session
}
type DB struct {
	db *mgo.Database
}
type Coll struct {
	c *mgo.Collection
}

func (s Sess) Session(work func(*mgo.Session)) {
	sess := s.s.Copy()
	defer sess.Close()
	work(sess)
}

func (db DB) Session(work func(*mgo.Database)) {
	sess := db.db.Session.Copy()
	defer sess.Close()
	work(db.db.With(sess))
}

func (c Coll) Session(work func(*mgo.Collection)) {
	sess := c.c.Database.Session.Copy()
	defer sess.Close()
	work(c.c.With(sess))
}

func (s Sess) DB(name string) DB {
	return DB{s.s.DB(name)}
}

func (db DB) C(name string) Coll {
	return Coll{db.db.C(name)}
}
