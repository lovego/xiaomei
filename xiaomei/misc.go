package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	cmdPkg "github.com/lovego/cmd"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func coverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `cover [package] ...`,
		Short: `show coverage details for packages.`,
		RunE: func(_ *cobra.Command, args []string) error {
			_, err := cmdPkg.Run(cmdPkg.O{}, "sh", "-c", fmt.Sprintf(`
rm -f /tmp/go-cover.out && {
  go test --gcflags=-l -coverprofile /tmp/go-cover.out %s
  test -f /tmp/go-cover.out && {
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

func autoCompleteCmd(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   `auto-complete`,
		Short: `setup shell auto completion.`,
		RunE: release.NoArgCall(func() error {
			var buf bytes.Buffer
			if err := rootCmd.GenBashCompletion(&buf); err != nil {
				return err
			}
			cmdPkg.SudoWriteFile(`/etc/bash_completion.d/xiaomei`, &buf)
			return nil
		}),
	}
}
