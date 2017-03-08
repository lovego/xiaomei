package new

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bughou-go/xiaomei/utils/fs"
	"github.com/bughou-go/xiaomei/xiaomei/z"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   `new <project-path>`,
		Short: `create a new project.`,
		RunE: z.Arg1Call(``, func(dir string) error {
			return New(dir)
		}),
	}
}

func New(proDir string) error {
	proPath, err := getProjectPath(proDir)
	if err != nil {
		return err
	}
	exampleDir, err := getExampleDir()
	if err != nil {
		return err
	}
	return execTemplates(exampleDir, proDir, proPath)
}

func execTemplates(exampleDir, proDir, proPath string) error {
	proName := filepath.Base(proPath)
	domainMapping, domainList := domainNames(proName)
	data := struct {
		ProPath, ProName, Secret, DomainList string
		DomainMapping                        map[string]string
	}{
		proPath, proName, genSecret(), domainList, domainMapping,
	}

	return filepath.Walk(exampleDir, func(src string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		dst := strings.Replace(src, exampleDir, proDir, 1)
		if info.IsDir() {
			return os.Mkdir(dst, 0755)
		} else {
			return copyFile(src, dst, info, data)
		}
	})
}

func domainNames(proName string) (map[string]string, string) {
	mapping := map[string]string{}
	lists := []string{}
	for _, env := range []string{`dev`, `test`, `qa`, `production`} {
		var domain string
		if env == `production` {
			domain = proName + `.com`
		} else {
			domain = proName + `.` + env
		}
		mapping[env] = domain
		lists = append(lists, domain)
	}
	return mapping, strings.Join(lists, ` `)
}

func copyFile(src, dst string, info os.FileInfo, data interface{}) error {
	dir, file := filepath.Split(dst)
	if !strings.HasPrefix(file, `tmpl.`) {
		return fs.Copy(src, dst)
	}
	dst = filepath.Join(dir, strings.TrimPrefix(file, `tmpl.`))
	if content, err := renderTmpl(src, data); err == nil {
		return ioutil.WriteFile(dst, content, info.Mode())
	} else {
		return err
	}
}

func renderTmpl(tmplPath string, data interface{}) ([]byte, error) {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 32 byte hex string
func genSecret() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
