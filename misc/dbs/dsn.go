package dbs

import (
	"net/url"
	"strings"
)

// Data Source Name
type DSN struct {
	Host, Port, User, Password, DB string
}

func ParseDSN(urlStr string) DSN {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic("invalid data source: " + urlStr)
	}
	parts := strings.Split(u.Path, `/`)
	if len(parts) != 2 {
		panic("invalid db in data source: " + urlStr)
	}

	password, _ := u.User.Password()
	dsn := DSN{
		Host:     u.Hostname(),
		Port:     u.Port(),
		User:     u.User.Username(),
		Password: password,
		DB:       strings.Split(u.Path, `/`)[1],
	}
	return dsn
}

func (dsn DSN) MysqlFlags() []string {
	return []string{`-h` + dsn.Host, `-P` + dsn.Port, `-u` + dsn.User, `-p` + dsn.Password, dsn.DB}
}

func (dsn DSN) RedisFlags() []string {
	flags := []string{`-h`, dsn.Host, `-p`, dsn.Port, `-n`, dsn.DB}
	if dsn.User != `` {
		flags = append(flags, `--user`, dsn.User)
	}
	if dsn.Password != `` {
		flags = append(flags, `-a`, dsn.Password)
	}
	return flags
}
