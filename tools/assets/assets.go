package assets

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/bughou-go/xiaomei/config"
	"gopkg.in/yaml.v2"
)

func Assets(args []string) {
	root := config.Root
	if len(args) == 0 {
		args = []string{`release/public/js`, `release/public/css`}
	}
	paths := getFilesPath(root, args)

	assetsPath := path.Join(root, `config/assets.yml`)
	assets := getAssets(assetsPath)
	changed := checkAndAddMd5(assets, paths)
	if changed == 0 {
		return
	}
	updateAssets(assetsPath, assets)
	fmt.Printf("assets %d files changed.\n", changed)
}

func getFilesPath(root string, args []string) []string {
	paths := []string{}
	for _, arg := range args {
		filePath := path.Join(root, `..`, arg)
		file, err := os.Stat(filePath)
		if err != nil {
			panic(err)
		}
		if file.IsDir() {
			if path.Base(filePath) != `font` {
				paths = append(paths, getDirFiles(filePath)...)
			}
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
		child := path.Join(dir, filename)
		if file.IsDir() {
			if path.Base(child) != `font` {
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
	publicDir := path.Join(config.Root, `public/reports`)
	for _, filePath := range paths {
		assetPath := strings.Replace(filePath, publicDir, ``, -1)
		if ext := path.Ext(assetPath); ext != `.css` && ext != `.js` {
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
