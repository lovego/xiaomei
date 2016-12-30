package db

import (
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/cmd"
	"github.com/bughou-go/xiaomei/utils/dsn"
)

func Mysql(key, serverFilter string) error {
	flags := dsn.Mysql(config.DB.Mysql(key)).Flags()
	return run(serverFilter, `mysql`, append([]string{`--pager=less -SX`}, flags...)...)
}

func MysqlDump(key, serverFilter string) error {
	flags := dsn.Mysql(config.DB.Mysql(key)).Flags()
	return run(serverFilter, `mysqldump`, append([]string{`-t`}, flags...)...)
}

func Mongo(key, serverFilter string) error {
	return run(serverFilter, `mongo`, config.DB.Mongo(key))
}

func Redis(key, serverFilter string) error {
	flags := dsn.Redis(config.DB.Redis(key)).Flags()
	return run(serverFilter, `redis-cli`, flags...)
}

func run(filter, command string, args ...string) error {
	if config.IsLocalEnv() {
		_, err := cmd.Run(cmd.O{}, command, args...)
		return err
	}
	servers := config.Servers.Matched2(filter, `appserver`)
	if len(servers) == 0 {
		return nil
	}
	_, err := cmd.Run(cmd.O{}, `ssh`,
		append([]string{`-t`, servers[0].SshAddr(), command}, args...)...,
	)
	return err
}
