package token

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `token`,
		Short: `Generate or parse a token.`,
	}
	cmd.AddCommand(genCmd())
	cmd.AddCommand(parseCmd())
	return cmd
}
