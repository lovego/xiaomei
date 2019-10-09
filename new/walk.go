package new

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/lovego/fs"
)

func walk(tmplsDir, proDir string, config *Config) error {
	return filepath.Walk(tmplsDir, func(src string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		dst := strings.Replace(src, tmplsDir, proDir, 1)
		if info.IsDir() {
			return os.MkdirAll(dst, 0755)
		} else {
			return copyFile(src, dst, info, config)
		}
	})
}

func copyFile(src, dst string, info os.FileInfo, config *Config) error {
	dir, file := filepath.Split(dst)
	if !strings.HasPrefix(file, `_`) {
		return fs.Copy(src, dst)
	}
	dst = filepath.Join(dir, strings.TrimPrefix(file, `_`))
	if content, err := renderTmpl(src, config); err == nil {
		return ioutil.WriteFile(dst, content, info.Mode())
	} else {
		return err
	}
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
