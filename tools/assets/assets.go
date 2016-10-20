package assets

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"github.com/bughou-go/xiaomei/config"
	"strings"
)

func Assets(args []string) {
	root := config.Root
	if len(args) == 0 {
		args = []string{`release/public/reports/js`, `release/public/reports/css`}
	}
	paths := getFilesPath(root, args)

	assets_path := path.Join(root, `config/assets.json`)
	assets := getAssets(assets_path)
	changed := checkAndAddMd5(assets, paths)
	if changed == 0 {
		return
	}
	updateAssets(assets_path, assets)
	fmt.Printf("assets %d files changed.\n", changed)
}

func getFilesPath(root string, args []string) []string {
	paths := []string{}
	for _, arg := range args {
		file_path := path.Join(root, `..`, arg)
		file, err := os.Stat(file_path)
		if err != nil {
			panic(err)
		}
		if file.IsDir() {
			if path.Base(file_path) != `font` {
				paths = append(paths, getDirFiles(file_path)...)
			}
		} else {
			paths = append(paths, file_path)
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
		child_file := path.Join(dir, filename)
		if file.IsDir() {
			if path.Base(child_file) != `font` {
				paths = append(paths, getDirFiles(child_file)...)
			}
		} else {
			paths = append(paths, child_file)
		}
	}
	return paths
}

func getAssets(assets_path string) (assets map[string]string) {
	data, err := ioutil.ReadFile(assets_path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &assets)
	if err != nil {
		panic(err)
	}
	return
}

func updateAssets(assets_path string, assets map[string]string) {
	content, err := json.MarshalIndent(assets, ``, `  `)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(assets_path, content, os.ModeDir)
	if err != nil {
		panic(err)
	}
}

func checkAndAddMd5(assets map[string]string, paths []string) (changed int) {
	public_dir := path.Join(config.Root, `public/reports`)
	for _, file_path := range paths {
		asset_path := strings.Replace(file_path, public_dir, ``, -1)
		if ext := path.Ext(asset_path); ext != `.css` && ext != `.js` {
			continue
		}
		file_md5 := getMd5(file_path)
		if assets[asset_path] != file_md5 {
			assets[asset_path] = file_md5
			changed++
		}
	}
	return
}

func getMd5(file_path string) string {
	data, err := ioutil.ReadFile(file_path)
	if err != nil {
		panic(err)
	}
	m := md5.Sum(data)
	return hex.EncodeToString(m[:])
}
