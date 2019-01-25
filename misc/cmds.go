package misc

import (
	"fmt"
	"io/ioutil"
	"strings"

	cmdPkg "github.com/lovego/cmd"
	"github.com/lovego/config/conf"
	"github.com/lovego/xiaomei/misc/dbs"
	"github.com/lovego/xiaomei/misc/godoc"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func Cmds(rootCmd *cobra.Command) []*cobra.Command {
	return append(
		dbs.Cmds(),
		godoc.Cmd(), timestampSignCmd(), specCmd(), coverCmd(), yamlCmd(), bashCompletionCmd(rootCmd),
	)
}

func timestampSignCmd() *cobra.Command {
	var secret string
	cmd := &cobra.Command{
		Use:   `timestamp-sign [<env>]`,
		Short: `generate Timestamp and Sign headers for curl command.`,
		RunE: release.EnvCall(func(env string) error {
			var ts int64
			var sign string
			if secret != "" {
				ts, sign = conf.TimestampSign(secret)
			} else {
				ts, sign = release.AppConf(env).TimestampSign()
			}
			fmt.Printf("-H Timestamp:%d -H Sign:%s\n", ts, sign)
			return nil
		}),
	}
	cmd.Flags().StringVarP(&secret, `secret`, `s`, ``, `secret used to generate sign`)
	return cmd
}

func coverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `cover [package] ...`,
		Short: `show coverage details for packages.`,
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
		Short: `parse yaml file.`,
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
