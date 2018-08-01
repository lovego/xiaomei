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

const (
	invender    = 0
	notinvender = 1
	depsTest    = 2
)

func depsCmd() *cobra.Command {
	var shortNameCmd int
	cmd := &cobra.Command{
		Use:   `deps`,
		Short: "list dependence packages.",
		Run: func(c *cobra.Command, args []string) {
			deps := getDeps(shortNameCmd)
			fmt.Println(strings.Join(deps, "\n"))
		},
	}
	fmt.Println("test ===---", cmd)

	cmd.Flags().IntVarP(&shortNameCmd, `in-vendor`, `v`, 0, `list dependences in vendor dir.`)
	cmd.Flags().IntVarP(&shortNameCmd, `exclude-test`, `t`, 2, `list dependences in test dir.`)
	return cmd
}

func getDeps(shortNameCmd int) []string {
	fmt.Println("test ===", shortNameCmd)
	switch shortNameCmd {
	case invender, notinvender:
		return invenderDeps(shortNameCmd)
	case depsTest:
		return testDeps()
	}
	return []string{}
}

func invenderDeps(shortNameCmd int) (deps []string) {
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
	if shortNameCmd == invender {
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
	return
}

func testDeps() (deps []string) {
	result, err := cmd.Run(
		cmd.O{Output: true, Dir: path.Join(release.Root(), `../`)},
		`go`, `list`, `-e`, `-f`, `{{ join .TestImports "\n" }}`, `./models/...`,
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
