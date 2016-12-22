package db

import (
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/utils/dsn"
)

func Redis(name string) {
	options := dsn.RedisDSN(config.DB.Redis(name)).Options()
	command, options := sshOptions(`redis-cli`, options)
	cmd.Run(cmd.O{Panic: true}, command, options...)
}
