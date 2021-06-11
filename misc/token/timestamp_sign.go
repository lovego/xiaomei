package token

import (
	"fmt"
	"time"

	"github.com/lovego/config/config"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func TimestampSignCmd() *cobra.Command {
	var secret string
	cmd := &cobra.Command{
		Use:   `timestamp-sign [env]`,
		Short: `Generate Timestamp and Sign headers for curl command.`,
		RunE: release.EnvCall(func(env string) error {
			ts := time.Now().Unix()
			if secret == "" {
				secret = release.EnvConfig(env).Secret
			}
			fmt.Printf("-H Timestamp:%d -H Sign:%s\n", ts, config.TimestampSign(ts, secret))
			return nil
		}),
	}
	cmd.Flags().StringVarP(&secret, `secret`, `s`, ``, `secret used to generate sign`)
	return cmd
}
