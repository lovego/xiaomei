package main

import (
	"fmt"
	"io/ioutil"

	"github.com/lovego/xiaomei/xiaomei/release"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

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
		RunE: release.NoArgCall(func() error {
			println(`xiaomei version 17.4.27`)
			return nil
		}),
	}
}
