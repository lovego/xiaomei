package token

import (
	"net/http"

	"github.com/lovego/config/conf"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `token`,
		Short: `Generate/parse token command.`,
	}
	cmd.AddCommand(tokenGenCmd())
	cmd.AddCommand(tokenParseCmd())
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
