package token

import (
	"net/http"

	"github.com/lovego/config/conf"
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

func newCookie(cookie conf.Cookie) *http.Cookie {
	return &http.Cookie{
		Name:   cookie.Name,
		Domain: cookie.Domain,
		Path:   cookie.Path,
		MaxAge: cookie.MaxAge,
	}
}
