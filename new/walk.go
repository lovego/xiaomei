package new

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/lovego/fs"
)

func walk(tmplsDir, proDir string, config *Config, force bool) error {
	return filepath.Walk(tmplsDir, func(src string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		dst := strings.Replace(src, tmplsDir, proDir, 1)
		if info.IsDir() {
			if err := os.MkdirAll(dst, 0755); err == os.ErrExist {
				return nil
			} else {
				return err
			}
		} else {
			return copyFile(src, dst, info, config, force)
		}
	})
}

func copyFile(src, dst string, info os.FileInfo, config *Config, force bool) error {
	dir, file := filepath.Split(dst)
	isTmpl := strings.HasPrefix(file, `__`)
	if isTmpl {
		dst = filepath.Join(dir, strings.TrimPrefix(file, `__`))
	}

	if !force && fs.Exist(dst) {
		return fmt.Errorf(`%s: aready exists, use "-f" flag to override.`, dst)
	}

	if isTmpl {
		if content, err := renderTmpl(src, config); err == nil {
			return ioutil.WriteFile(dst, content, info.Mode())
		} else {
			return err
		}
	}
	return fs.Copy(src, dst)
}

func renderTmpl(tmplPath string, config *Config) ([]byte, error) {
	tmpl, err := template.New(filepath.Base(tmplPath)).Funcs(template.FuncMap{
		"genSecret": genSecret,
	}).ParseFiles(tmplPath)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
