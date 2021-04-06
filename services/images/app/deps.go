package app

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/lovego/cmd"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func depsCmd() *cobra.Command {
	var inVendor bool
	var excludeTest bool
	cmd := &cobra.Command{
		Use:   `deps`,
		Short: "List dependence packages.",
		RunE: func(c *cobra.Command, args []string) error {
			deps, err := getDeps(inVendor, excludeTest)
			if err == nil {
				fmt.Println(strings.Join(deps, "\n"))
			}
			return err
		},
	}
	cmd.Flags().BoolVarP(&inVendor, `in-vendor`, `v`, false, `list dependences in vendor dir.`)
	cmd.Flags().BoolVarP(&excludeTest, `exclude-test`, `e`, false, ` exclude test dependences.`)
	return cmd
}

func getDeps(inVendor, excludeTest bool) (deps []string, err error) {
	modulePath, err := release.ModulePath()
	if err != nil {
		return nil, err
	}

	pkgs := getDepPkgs(excludeTest)

	vendorPath := path.Join(modulePath, `vendor`)

	if inVendor {
		for _, pkg := range pkgs {
			if strings.HasPrefix(pkg, vendorPath) {
				deps = append(deps, pkg)
			}
		}
	} else {
		for _, pkg := range pkgs {
			if strings.Contains(pkg, `.`) && !strings.HasPrefix(pkg, modulePath) {
				deps = append(deps, pkg)
			}
		}
	}

	return
}

func getDepPkgs(excludeTest bool) []string {
	o := cmd.O{Output: true, Dir: path.Join(release.Root(), `../`)}
	result, err := cmd.Run(
		o, release.GoCmd(), `list`, `-f`, `{{join .Deps "\n"}}`,
	)
	if err != nil {
		log.Panic(err)
	}
	pkgs := strings.Split(result, "\n")

	if !excludeTest {
		result, err := cmd.Run(
			o, release.GoCmd(), `list`, `-f`, `{{join .TestImports "\n"}}`, `./models/...`,
		)
		if err != nil {
			log.Panic(err)
		}
		pkgs = appendIfNotExists(pkgs, strings.Split(result, "\n"))
	}
	return pkgs
}

func appendIfNotExists(a, b []string) []string {
	m := make(map[string]bool)
	for _, elem := range a {
		m[elem] = true
	}
	for _, elem := range b {
		if !m[elem] {
			a = append(a, elem)
			m[elem] = true
		}
	}
	return a
}
