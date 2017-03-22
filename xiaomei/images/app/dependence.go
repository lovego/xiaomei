package app

import (
	"fmt"
	"path"
	"strings"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/utils/fs"
	"github.com/lovego/xiaomei/utils/slice"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func DepsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `deps`,
		Short: "list dependence packages.",
		Run: func(c *cobra.Command, args []string) {
			listDeps()
		},
	}
	return cmd
}

func listDeps() {
	deps := getAllDeps()
	vendorDeps := []string{}
	notVendorDeps := []string{}
	for _, dep := range deps {
		i := strings.Index(dep, `vendor/`)
		if i > -1 {
			vendorDeps = append(vendorDeps, dep[i+7:]) // 7 = len(`vendor/`)
		} else {
			notVendorDeps = append(notVendorDeps, dep)
		}
	}
	cmd.Run(cmd.O{}, `echo`, fmt.Sprintf("\ndependece in vendor:\n%s", strings.Join(vendorDeps, "\n")))
	cmd.Run(cmd.O{}, `echo`, fmt.Sprintf("\ndependece not in vendor:\n%s", strings.Join(notVendorDeps, "\n")))
}

func getAllDeps() []string {
	projectDir := path.Join(release.Root(), `../`)
	return getDirDeps(projectDir)
}

var already = make(map[string]bool)

func getDirDeps(dir string) []string {
	if already[dir] {
		return []string{}
	}
	goSrcPath, err := fs.GetGoSrcPath()
	if err != nil {
		panic(err)
	}
	result, _ := cmd.Run(cmd.O{Output: true, Dir: dir}, `go`, `list`, `-e`, `-f`, `{{join .Imports "\n"}}`)
	already[dir] = true
	deps := filterStandard(strings.Split(result, "\n"))
	for _, depPath := range deps {
		if strings.HasPrefix(depPath, release.Path()) {
			childDeps := filterStandard(getDirDeps(path.Join(goSrcPath, depPath)))
			for _, childDep := range childDeps {
				if !slice.ContainsString(deps, childDep) {
					deps = append(deps, childDep)
				}
			}
		}
	}
	return filterDeps(deps)
}

// 过滤xiaomei和项目内的包
func filterDeps(deps []string) []string {
	pkgs := []string{}
	for _, dep := range deps {
		if strings.HasPrefix(dep, path.Join(release.Path(), `vendor`)) ||
			!strings.HasPrefix(dep, `github.com/lovego/xiaomei`) &&
				!strings.HasPrefix(dep, release.Path()) {
			pkgs = append(pkgs, dep)
		}
	}
	return pkgs
}

// 过滤标准库的包
func filterStandard(deps []string) []string {
	pkgs := []string{}
	for _, dep := range deps {
		if strings.Contains(dep, `.`) {
			pkgs = append(pkgs, dep)
		}
	}
	return pkgs
}
