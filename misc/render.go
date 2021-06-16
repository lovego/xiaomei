package misc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func renderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `render <env> <template file> [output file]`,
		Short: `render template with envionment config.`,
		RunE: release.Env2Call(func(env string, tmplFile, outputFile string) error {
			return RenderFileWithEnvConfig(env, tmplFile, outputFile)
		}),
	}
	return cmd
}

func RenderFileWithEnvConfig(env, tmplFile, outputFile string) error {
	if outputFile != "" {
		return RenderFileTo(tmplFile, nil, release.EnvConfig(env), outputFile)
	}
	if output, err := RenderFile(tmplFile, nil, release.EnvConfig(env)); err != nil {
		return err
	} else {
		fmt.Println(output.String())
		return nil
	}
}

func RenderFileTo(tmplFile string, funcs template.FuncMap, data interface{}, outputFile string) error {
	output, err := RenderFile(tmplFile, funcs, data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outputFile, output.Bytes(), 0644)
}

func RenderFile(tmplFile string, funcs template.FuncMap, data interface{}) (bytes.Buffer, error) {
	var buf bytes.Buffer
	content, err := ioutil.ReadFile(tmplFile)
	if err != nil {
		return buf, err
	}
	tmpl := template.New(``)
	if funcs != nil {
		tmpl.Funcs(funcs)
	}
	if _, err := tmpl.Parse(string(content)); err != nil {
		return buf, err
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return buf, err
	}
	return buf, nil
}
