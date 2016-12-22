package db

import (
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/cmd"
)

func Redis(name string) {
	options := utils.RedisOptions(config.DB.Redis(name))
	command, options := sshOptions(`redis-cli`, options)
	cmd.Run(cmd.O{Panic: true}, command, options...)
}
