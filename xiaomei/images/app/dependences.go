package app

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
)

func depsCmd() *cobra.Command {
	var inVendor bool
	var excludeTest bool
	cmd := &cobra.Command{
		Use:   `deps`,
		Short: "list dependence packages.",
		Run: func(c *cobra.Command, args []string) {
			deps := getDeps(inVendor, excludeTest)
			fmt.Println(strings.Join(deps, "\n"))
		},
	}
	cmd.Flags().BoolVarP(&inVendor, `in-vendor`, `v`, false, `list dependences in vendor dir.`)
	cmd.Flags().BoolVarP(&excludeTest, `exclude-test`, `e`, false, `list dependences exclude test.`)
	return cmd
}

func getDeps(inVendor, excludeTest bool) (deps []string) {
	result, err := cmd.Run(
		cmd.O{Output: true, Dir: path.Join(release.Root(), `../`)},
		`go`, `list`, `-e`, `-f`, `{{join .Deps "\n"}}`,
	)
	if err != nil {
		log.Panic(err)
	}

	pkgs := strings.Split(result, "\n")
	projectPath := release.Path()
	vendorPath := path.Join(projectPath, `vendor`)
	if inVendor {
		for _, pkg := range pkgs {
			if strings.HasPrefix(pkg, vendorPath) {
				deps = append(deps, pkg)
			}
		}
	} else {
		for _, pkg := range pkgs {
			if strings.Contains(pkg, `.`) && !strings.HasPrefix(pkg, projectPath) {
				deps = append(deps, pkg)
			}
		}
	}
	if !excludeTest {
		deps = append(deps, getTestDeps()...)
	}

	return
}

func getTestDeps() (deps []string) {
	result, err := cmd.Run(
		cmd.O{Output: true, Dir: path.Join(release.Root(), `../`)},
		`go`, `list`, `-e`, `-f`, `{{join .TestImports "\n"}}`, `./models/...`,
	)
	if err != nil {
		log.Panic(err)
	}
	pkgs := strings.Split(result, "\n")
	projectPath := release.Path()
	for _, pkg := range pkgs {
		if strings.Contains(pkg, `.`) && !strings.HasPrefix(pkg, projectPath) {
			deps = append(deps, pkg)
		}
	}
	return
}
