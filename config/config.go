package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lovego/fs"
	"github.com/lovego/mailer"
	"github.com/lovego/xiaomei/config/conf"
)

var theConf = conf.Get(filepath.Join(Root(), `config/config.yml`)).Get(Env())
var theData = conf.Data(filepath.Join(Root(), `config/envs/`+Env()+`.yml`))
var theMailer = getMailer()

var theRoot string

func Root() string {
	if theRoot == `` {
		program, err := filepath.Abs(os.Args[0])
		if err != nil {
			log.Panic(err)
		}
		if dir := filepath.Dir(program); fs.Exist(filepath.Join(dir, `config/config.yml`)) {
			theRoot = dir
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				log.Panic(err)
			}
			projectDir := fs.DetectDir(cwd, `release/img-app/config/config.yml`)
			if projectDir != `` {
				theRoot = filepath.Join(projectDir, `release/img-app`)
			} else {
				log.Panic(`app root not found.`)
			}
		}
	}
	return theRoot
}

var theEnv string

func Env() string {
	if theEnv == `` {
		theEnv = os.Getenv(`GOENV`)
		if theEnv == `` {
			if strings.HasSuffix(os.Args[0], `.test`) {
				theEnv = `test`
			} else {
				theEnv = `dev`
			}
		}
	}
	return theEnv
}

func getMailer() *mailer.Mailer {
	m, err := mailer.New(theConf.Mailer)
	if err != nil {
		panic(err)
	}
	return m
}
