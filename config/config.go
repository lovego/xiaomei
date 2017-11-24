package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lovego/xiaomei/config/conf"
	"github.com/lovego/fs"
	"github.com/lovego/mailer"
)

var theConf = conf.Get(Root()).Get(Env())
var theData = conf.Data(Root(), Env())
var theMailer = getMailer()

var theRoot string

func Root() string {
	if theRoot == `` {
		program, err := filepath.Abs(os.Args[0])
		if err != nil {
			log.Panic(err)
		}
		if strings.HasSuffix(program, `.test`) /* go test ... */ ||
			strings.HasPrefix(program, `/tmp/`) /* go run ... */ {
			cwd, err := os.Getwd()
			if err != nil {
				log.Panic(err)
			}
			projectDir := fs.DetectDir(cwd, `release/deploy.yml`)
			if projectDir != `` {
				theRoot = filepath.Join(projectDir, `release/img-app`)
			} else {
				log.Panic(`app root not found.`)
			}
		} else { // project binary file
			theRoot = filepath.Dir(program)
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
