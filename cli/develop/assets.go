package develop

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bughou-go/xiaomei/config"
	"gopkg.in/yaml.v2"
)

func Assets(args []string) {
	if len(args) == 0 {
		args = []string{`release/public/js`, `release/public/css`}
	}
	paths := getFilesPath(args)

	assetsPath := filepath.Join(config.Root(), `config/assets.yml`)
	assets := getAssets(assetsPath)
	changed := checkAndAddMd5(assets, paths)
	if changed == 0 {
		return
	}
	updateAssets(assetsPath, assets)
	fmt.Printf("%d assets files changed.\n", changed)
}

func getFilesPath(args []string) []string {
	paths := []string{}
	for _, arg := range args {
		filePath := filepath.Join(config.Root(), `..`, arg)
		fi, err := os.Stat(filePath)
		if err != nil {
			panic(err)
		}
		if fi.IsDir() {
			paths = append(paths, getDirFiles(filePath)...)
		} else {
			paths = append(paths, filePath)
		}
	}
	return paths
}

func getDirFiles(dir string) (paths []string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		filename := file.Name()
		child := filepath.Join(dir, filename)
		if file.IsDir() {
			if filepath.Base(child) != `font` {
				paths = append(paths, getDirFiles(child)...)
			}
		} else {
			paths = append(paths, child)
		}
	}
	return paths
}

func getAssets(assetsPath string) (assets map[string]string) {
	data, err := ioutil.ReadFile(assetsPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &assets)
	if err != nil {
		panic(err)
	}
	return
}

func updateAssets(assetsPath string, assets map[string]string) {
	content, err := yaml.Marshal(assets)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(assetsPath, content, os.ModeDir)
	if err != nil {
		panic(err)
	}
}

func checkAndAddMd5(assets map[string]string, paths []string) (changed int) {
	publicDir := filepath.Join(config.Root(), `public/reports`)
	for _, filePath := range paths {
		assetPath := strings.Replace(filePath, publicDir, ``, -1)
		if ext := filepath.Ext(assetPath); ext != `.css` && ext != `.js` {
			continue
		}
		nowMd5 := getMd5(filePath)
		if assets[assetPath] != nowMd5 {
			assets[assetPath] = nowMd5
			changed++
		}
	}
	return
}

func getMd5(filePath string) string {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	m := md5.Sum(data)
	return hex.EncodeToString(m[:])
}
