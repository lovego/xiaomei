package db

import (
	"github.com/bughou-go/xiaomei/config"
	"gopkg.in/mgo.v2"
	"retail-tek.com/reports/utils/cmd"
	"strings"
)

type MongoConfig struct {
	User, Passwd, Host, Db string
}

var mongoConfig MongoConfig

func Mongo(key string) {
	options := getMongoOptions(key)
	command, options := sshOptions(`mongo`, options)
	cmd.Run(cmd.O{Panic: true}, command, options...)
}

func getMongoOptions(key string) []string {
	c := getMongoConfig(key)
	options := []string{}
	if c.User != `` {
		options = append(options, `-u`+c.User)
	}
	if c.Passwd != `` {
		options = append(options, `-p`+c.Passwd)
	}
	options = append(options, `--host`, c.Host, c.Db)
	return options
}

func getMongoConfig(key string) MongoConfig {
	if mongoConfig.User == `` {
		url := config.DB.Mongo(key)
		info, err := mgo.ParseURL(url)
		if err != nil {
			panic(err)
		}
		mongoConfig = MongoConfig{
			info.Username,
			info.Password,
			makeMongoHost(info.ReplicaSetName, info.Addrs),
			info.Database,
		}
	}
	return mongoConfig
}

func makeMongoHost(replicName string, addrs []string) string {
	ips := strings.Join(addrs, `:27017,`)
	if replicName != `` {
		return replicName + `/` + ips + `:27017`
	}
	return ips + `:27017`
}
