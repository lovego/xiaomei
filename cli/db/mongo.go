package db

import (
	"github.com/bughou-go/xiaomei/config"
	"retail-tek.com/reports/utils/cmd"
)

type MongoConfig struct {
	User, Passwd, Host, Db string
}

var mongoConfig MongoConfig

func Mongo(key string) {
	command, options := sshOptions(`mongo`, []string{config.DB.Mongo(key)})
	cmd.Run(cmd.O{Panic: true}, command, options...)
}
