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
	invender    = 1
	notinvender = 2
	depsTest    = 3
)

func depsCmd() *cobra.Command {
	var shortNameCmd int

	var cmd = &cobra.Command{
		Use:   `deps`,
		Short: "list dependence packages.",
	}

	cmd.Flags().IntVarP(&shortNameCmd, `in-vendor`, `v`, 1, `list dependences in vendor dir.`)
	cmd.Flags().IntVarP(&shortNameCmd, `exclude-test`, `t`, 3, `list dependences in test dir.`)
	fmt.Println("test111111", shortNameCmd)

	cmd.Run = func(c *cobra.Command, args []string) {
		fmt.Println("test2222222", shortNameCmd)
		deps := getDeps(shortNameCmd)
		fmt.Println(strings.Join(deps, "\n"))
	}
	fmt.Println("test ===---", cmd)

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
	fmt.Println("result ==", deps, path.Join(release.Root(), `../`))
	return
}
