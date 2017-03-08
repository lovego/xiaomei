package web

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `web`,
		Short: `the web server.`,
	}
	return cmd
}
