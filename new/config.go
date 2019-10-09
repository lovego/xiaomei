package new

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"path/filepath"
	"strings"

	"github.com/lovego/fs"
)

type Config struct {
	ProPath        string
	ProName        string
	ProNameUrlSafe string
	Domain         string
	Registry       string
}

func getConfig(typ, dir, registry, domain string) (*Config, error) {
	var proPath string
	if typ == `app` || typ == `full` {
		var err error
		if proPath, err = getProjectPath(dir); err != nil {
			return nil, err
		}
	}
	proName := filepath.Base(dir)
	return &Config{
		ProPath:        proPath,
		ProName:        proName,
		ProNameUrlSafe: strings.Replace(proName, `_`, `-`, -1),
		Domain:         domain,
		Registry:       registry,
	}, nil
}

func getProjectPath(dir string) (string, error) {
	srcPath := fs.GetGoSrcPath()
	proPath, err := filepath.Rel(srcPath, dir)
	if err != nil {
		return ``, err
	}
	if proPath[0] == '.' {
		return ``, errors.New(`project dir must be under ` + srcPath + "\n")
	}
	return proPath, nil
}

// 32 byte hex string
func genSecret() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
