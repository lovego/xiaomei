package misc

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	cmdPkg "github.com/lovego/cmd"
	"github.com/lovego/xiaomei/misc/dbs"
	"github.com/lovego/xiaomei/misc/godoc"
	"github.com/lovego/xiaomei/misc/token"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func Cmds(rootCmd *cobra.Command) []*cobra.Command {
	return append(
		dbs.Cmds(), godoc.Cmd(), docCmd(), token.Cmd(), token.TimestampSignCmd(),
		specCmd(), coverCmd(), yamlCmd(), float32Cmd(), float64Cmd(),
		bashCompletionCmd(rootCmd),
	)
}

func docCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `doc`,
		Short: `Generate api documentation from goa router.`,
		RunE: func(_ *cobra.Command, args []string) error {
			_, err := cmdPkg.Run(cmdPkg.O{
				Dir: filepath.Dir(release.Root()),
				Env: []string{"GOA_DOC=1"},
			}, release.GoCmd(), "run", "main.go")
			return err
		},
	}
	return cmd
}

func coverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `cover [package] ...`,
		Short: `Show coverage details for packages.`,
		RunE: func(_ *cobra.Command, args []string) error {
			script := os.Expand(`
rm -f /tmp/cover.out && {
  $GoCmd test -p 1 --gcflags=-l -coverprofile /tmp/cover.out $packages
  test -f /tmp/cover.out && (($(wc -c </tmp/cover.out) > 10)) && {
    $GoCmd tool cover -func /tmp/cover.out | tail -n 1
    $GoCmd tool cover -html /tmp/cover.out
  }
}`, func(name string) string {
				switch name {
				case `GoCmd`:
					return release.GoCmd()
				case `packages`:
					return strings.Join(args, " ")
				default:
					return ``
				}
			})
			_, err := cmdPkg.Run(cmdPkg.O{}, "sh", "-c", script)
			return err
		},
	}
	return cmd
}

func yamlCmd() *cobra.Command {
	var goSyntax bool
	cmd := &cobra.Command{
		Use:   `yaml`,
		Short: `Parse yaml file.`,
		RunE: release.Arg1Call(``, func(p string) error {
			content, err := ioutil.ReadFile(p)
			if err != nil {
				return err
			}
			data := make(map[string]interface{})
			if err := yaml.Unmarshal(content, data); err != nil {
				return err
			}
			if goSyntax {
				fmt.Printf("%#v\n", data)
			} else {
				if buf, err := yaml.Marshal(data); err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("%s\n", buf)
				}
			}
			return nil
		}),
	}
	cmd.Flags().BoolVarP(&goSyntax, `go-syntax`, `g`, false, `print in golang syntax`)
	return cmd
}

func float32Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   `float32`,
		Short: `Print IEEE 754 binary bits of float32.`,
		RunE: release.Arg1Call(``, func(p string) error {
			f, err := strconv.ParseFloat(p, 32)
			if err != nil {
				return err
			}
			s := fmt.Sprintf("%032b", math.Float32bits(float32(f)))
			fmt.Printf("%s,%s,%s\n", s[:1], s[1:9], s[9:])
			return nil
		}),
	}
}

func float64Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   `float64`,
		Short: `Print IEEE 754 binary bits of float64.`,
		RunE: release.Arg1Call(``, func(p string) error {
			f, err := strconv.ParseFloat(p, 64)
			if err != nil {
				return err
			}
			s := fmt.Sprintf("%064b", math.Float64bits(f))
			fmt.Printf("%s,%s,%s\n", s[:1], s[1:12], s[11:])
			return nil
		}),
	}
}
