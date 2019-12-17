package token

import (
	"fmt"
	"net/http"

	"github.com/lovego/sessions/cookiestore"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func genCmd() *cobra.Command {
	var domain string
	cmd := &cobra.Command{
		Use: `gen <env> <content> [flags]
<content> can be any string, but generally it may be a json encoded string`,
		DisableFlagsInUseLine: true,
		Short: `generate a token. (add a signature to a content and encode it)`,
		RunE: release.Env1Call(func(env, content string) error {
			ck := newCookie(release.AppConf(env).Cookie)
			if domain != `` {
				ck.Domain = domain
			}
			return generate(ck, release.AppConf(env).Secret, content)
		}),
	}
	cmd.Flags().StringVarP(&domain, `domain`, ``, ``, `specify the cookie domain.`)
	return cmd
}

func generate(ck *http.Cookie, secret, content string) error {
	encoded, err := cookiestore.New(secret).Encode(ck.Name, []byte(content))
	if err != nil {
		return err
	}

	ck.Value = string(encoded)

	fmt.Println(fmt.Sprintf(`document.cookie = "%s";`, ck.String()))
	return nil
}
