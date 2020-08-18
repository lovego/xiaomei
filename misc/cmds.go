package misc

import (
	"fmt"
	"io/ioutil"
	"math"
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
		dbs.Cmds(), godoc.Cmd(), token.Cmd(), token.TimestampSignCmd(),
		specCmd(), coverCmd(), yamlCmd(), float32Cmd(), float64Cmd(),
		bashCompletionCmd(rootCmd),
	)
}

func coverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `cover [package] ...`,
		Short: `Show coverage details for packages.`,
		RunE: func(_ *cobra.Command, args []string) error {
			_, err := cmdPkg.Run(cmdPkg.O{}, "sh", "-c", fmt.Sprintf(`
rm -f /tmp/go-cover.out && {
  go test -p 1 --gcflags=-l -coverprofile /tmp/go-cover.out %s
  test -f /tmp/go-cover.out && (($(wc -c </tmp/go-cover.out) > 10)) && {
    go tool cover -func /tmp/go-cover.out | tail -n 1
    go tool cover -html /tmp/go-cover.out
  }
}`, strings.Join(args, " ")))
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
	cmd.Flags().BoolVarP(&goSyntax, `go-syntax`, `g`, false, `print in go syntax`)
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
