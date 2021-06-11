package token

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/lovego/sessions/cookiestore"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func parseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: `parse <env> <token>
<token> is a pure cookie value, without cookie name or other attributes.`,
		DisableFlagsInUseLine: true,
		Short:                 `parse a token. (decode a token and remove the signature)`,
		RunE: release.Env1Call(func(env, token string) error {
			if token == `` {
				return errors.New(`token cann't be empty.`)
			}
			ck := release.EnvConfig(env).HttpCookie()
			ck.Value = token
			return parse(&ck, release.EnvConfig(env).Secret)
		}),
	}
	return cmd
}

func parse(ck *http.Cookie, secret string) error {
	if data, err := cookiestore.New(secret).Decode(
		ck.Name, []byte(ck.Value), int64(ck.MaxAge),
	); err != nil {
		return err
	} else {
		fmt.Println(string(data))
	}
	return nil
}
