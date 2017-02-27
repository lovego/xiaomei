package deps

import (
	"github.com/bughou-go/xiaomei/utils/cmd"
	"strings"
)

func List() {
	deps, _ := cmd.Run(cmd.O{Output: true}, `go`, `list`, `-e`, `-f`, `'{{join .Deps "\n" }}'`)
	pkgs := []string{}
	for _, dep := range strings.Split(deps, "\n") {
		if strings.Contains(dep, `.`) &&
			!strings.HasPrefix(dep, `github.com/bughou-go/xiaomei`) {
			pkgs = append(pkgs, dep)
		}
	}
	cmd.Run(cmd.O{}, `echo`, strings.Join(pkgs, "\n"))
}
