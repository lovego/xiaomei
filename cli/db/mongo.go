package db

import (
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Mongo(key string) {
	command, options := sshOptions(`mongo`, []string{config.DB.Mongo(key)})
	cmd.Run(cmd.O{Panic: true}, command, options...)
}
