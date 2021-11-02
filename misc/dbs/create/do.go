package create

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/lovego/xiaomei/release"
)

func Do(env, typ, key string, dropDB bool) error {
	if dropDB {
		switch env {
		case `production`, `staging`, `preview`:
			return errors.New("dropDB is forbidden under environment: " + env)
		}
	}
	keys, err := getKeys(key)
	if err != nil {
		return err
	}
	for _, key := range keys {
		if err := create(env, typ, key, dropDB); err != nil {
			return err
		}
	}
	return nil
}

func create(env, typ, key string, dropDB bool) error {
	content, err := ioutil.ReadFile(filepath.Join(release.ProjectRoot(), `sqls`, key+`.sql`))
	if err != nil {
		return err
	}
	return Creation{env: env, typ: typ, key: key, dropDB: dropDB, sql: string(content)}.do()
}

func getKeys(key string) ([]string, error) {
	var keys []string
	if key != "" {
		keys = []string{key}
	} else {
		files, err := filepath.Glob(filepath.Join(release.ProjectRoot(), `sqls`, `*.sql`))
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			keys = append(keys, strings.TrimSuffix(filepath.Base(file), `.sql`))
		}
	}
	return keys, nil
}
