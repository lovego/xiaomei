package new

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/lovego/fs"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ProPath        string `yaml:"-"`
	ProName        string `yaml:"-"`
	ProNameUrlSafe string `yaml:"-"`
	Registry       string `yaml:"registry"`
}

func getConfig(proDir string) (*Config, error) {
	proPath, err := getProjectPath(proDir)
	if err != nil {
		return nil, err
	}
	config, err := parseConfig()
	if err != nil {
		return nil, err
	}
	config.ProPath = proPath
	config.ProName = filepath.Base(proPath)
	config.ProNameUrlSafe = strings.Replace(config.ProName, `_`, `-`, -1)
	return config, nil
}

func getProjectPath(proDir string) (string, error) {
	if proDir == `` {
		return ``, errors.New(`project path can't be empty.`)
	}

	if !filepath.IsAbs(proDir) {
		var err error
		if proDir, err = filepath.Abs(proDir); err != nil {
			return ``, err
		}
	}

	srcPath := fs.GetGoSrcPath()
	proPath, err := filepath.Rel(srcPath, proDir)
	if err != nil {
		return ``, err
	}
	if proPath[0] == '.' {
		return ``, errors.New(`project dir must be under ` + srcPath + "\n")
	}
	return proPath, nil
}

func parseConfig() (*Config, error) {
	config := &Config{}
	if configPath := filepath.Join(os.Getenv(`HOME`), `.xiaomei.yml`); fs.Exist(configPath) {
		if content, err := ioutil.ReadFile(configPath); err != nil {
			return nil, err
		} else {
			if err := yaml.Unmarshal(content, config); err != nil {
				return nil, err
			}
		}
	}
	if config.Registry != `` && !strings.HasSuffix(config.Registry, `/`) {
		config.Registry += `/`
	}
	return config, nil
}

// 32 byte hex string
func genSecret() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
