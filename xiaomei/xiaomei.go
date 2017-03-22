package main

import (
	"fmt"
	"io/ioutil"
	"os"

	// "github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/images/access"
	"github.com/lovego/xiaomei/xiaomei/images/app"
	"github.com/lovego/xiaomei/xiaomei/images/web"
	"github.com/lovego/xiaomei/xiaomei/new"
	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/lovego/xiaomei/xiaomei/z"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func main() {
	cobra.EnableCommandSorting = false

	appCmd := app.Cmd()
	webCmd := web.Cmd()
	accessCmd := access.Cmd()
	clusterCmd := cluster.Cmd()
	appCmd.AddCommand(commonCmds(`app`)...)
	webCmd.AddCommand(commonCmds(`web`)...)
	accessCmd.AddCommand(commonCmds(`access`)...)

	root := rootCmd()
	root.AddCommand(appCmd, webCmd, accessCmd, clusterCmd)
	root.AddCommand(commonCmds(``)...)
	root.AddCommand(new.Cmd(), yamlCmd(), versionCmd())
	root.Execute()
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `xiaomei`,
		Short: `be small and beautiful.`,
	}

	if release.Arg1IsEnv() {
		cmd.SetArgs(os.Args[2:])
	}
	return cmd
}

func yamlCmd() *cobra.Command {
	var goSyntax bool
	cmd := &cobra.Command{
		Use:   `yaml`,
		Short: `parse yaml file.`,
		RunE: z.Arg1Call(``, func(p string) error {
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
				fmt.Println(data)
			}
			return nil
		}),
	}
	cmd.Flags().BoolVarP(&goSyntax, `go-syntax`, `g`, false, `print in go syntax`)
	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   `version`,
		Short: `show xiaomei version.`,
		RunE: z.NoArgCall(func() error {
			println(`xiaomei version 17.3.22`)
			return nil
		}),
	}
}
