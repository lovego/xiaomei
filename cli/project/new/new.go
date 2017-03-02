package new

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bughou-go/xiaomei/utils/cmd"
)

func New(dir string) error {
	proPath, proDir, err := checkProjectDir(dir)
	if err != nil {
		return err
	}
	if err := copyTemplates(proDir); err != nil {
		return err
	}

	return processTemplates(proPath, proDir)
}

func copyTemplates(proDir string) error {
	exampleDir, err := getExampleDir()
	if err != nil {
		return err
	}
	if !cmd.Ok(cmd.O{}, `cp`, `-rT`, exampleDir, proDir) {
		return errors.New(`cp templates failed.`)
	}
	return nil
}

func processTemplates(proPath, proDir string) error {
	proName := filepath.Base(proPath)
	script := fmt.Sprintf(`
	sed -i'' 's/example/%s/g' .gitignore $(fgrep -rl example release/config)
	sed -i'' 's/%s/%s/g' main.go
	sed -i'' 's/secret-string/%s/g' release/config/envs/production.yml
	`, proDir, proName,
		strings.Replace(examplePath, `/`, `\/`, -1),
		strings.Replace(proPath, `/`, `\/`, -1),
		generateSecret(),
	)
	if !cmd.Ok(cmd.O{Dir: proDir}, `sh`, `-c`, script) {
		return errors.New(`process templates failed.`)
	}
	return nil
}

// 32 byte hex string
func generateSecret() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
