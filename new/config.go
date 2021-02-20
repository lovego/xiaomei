package new

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"path/filepath"
	"strings"
)

type Config struct {
	ProName        string
	ProNameUrlSafe string
	Registry       string
	RepoPrefix     string
	ModulePath     string
	Domain         string
}

func getConfig(typ, projectDir, registry, domain string) (*Config, error) {
	var proName = filepath.Base(projectDir)
	if proName == `` {
		return nil, errors.New(`project name can't be empty.`)
	}
	if registry != "" && registry[len(registry)-1] != '/' {
		registry += "/"
	}

	config := &Config{
		ProName:        proName,
		ProNameUrlSafe: strings.Replace(proName, `_`, `-`, -1),
		Registry:       registry,
	}
	if i := strings.IndexByte(registry, '/'); i > 0 {
		config.RepoPrefix = strings.TrimSuffix(registry[i:], "/")
	}
	switch typ {
	case `app`:
		config.ModulePath = projectDir
	case `logc`:
		config.Domain = domain
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
