package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/lovego/sessions/cookiestore"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func tokenParseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `parse <env> <value>`,
		Short: `parse token command.`,
		RunE: release.Env1Call(func(env, value string) error {
			if value == `` {
				return errors.New(`cookie value is required`)
			}
			ck := newCookie(release.AppConf(env).Cookie)
			return tokenParse(ck, release.AppConf(env).Secret, value)
		}),
	}
	return cmd
}

func tokenParse(ck *http.Cookie, secret, value string) error {
	data := make(map[string]interface{})
	ck.Value = value
	if err := cookiestore.New(secret).Get(ck, &data); err != nil {
		return err
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
