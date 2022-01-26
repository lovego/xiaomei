package misc

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lovego/email"
	"github.com/lovego/httputil"
	"github.com/lovego/tracer"
	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

func monitorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `monitor <config file>`,
		Short: `Monitor http/https urls.`,
		RunE: release.Arg1Call(``, func(configFile string) error {
			return monitor(configFile)
		}),
	}
	return cmd
}

type Config struct {
	Mailer  string   `yaml:"mailer"`
	Keepers []string `yaml:"keepers"`

	Loops        uint
	Addrs        []string `yaml:"addrs"`
	SlowResponse uint16   `yaml:"slowResponse"`
}

func monitor(configFile string) error {
	conf, err := setupConfig(configFile)
	if err != nil {
		return err
	}
	if conf.Loops == 0 {
		conf.Loops = 1
	}

	for i := uint(0); i < conf.Loops; i++ {
		for _, addr := range conf.Addrs {
			if err := monitorOneAddr(addr, conf.SlowResponse); err != nil {
				return err
			}
		}
	}
	return nil
}

func monitorOneAddr(addr string, slowTime uint16) error {
	ctx := tracer.Start(context.Background(), "")
	resp, err := httputil.GetCtx(ctx, "", addr, nil, nil)
	t := tracer.Get(ctx).Children[0]

	var problem string
	if err != nil {
		problem = fmt.Sprintf(`请求失败：%s`, err)
	} else if resp.StatusCode != http.StatusOK {
		problem = fmt.Sprintf(`状态异常：%d`, resp.StatusCode)
	} else if uint16(t.Duration) > slowTime {
		problem = fmt.Sprintf(`耗时超过：%v`, time.Duration(slowTime)*time.Millisecond)
	}

	if problem != "" || os.Getenv("debug") != "" {
		msg, err := makeMessage(addr, problem, t)
		if err != nil {
			return err
		}
		if problem != "" {
			return sendAlarm(msg)
		} else {
			fmt.Println(msg)
		}
	}
	return nil
}

func makeMessage(addr, problem string, t *tracer.Tracer) (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	if problem != "" {
		problem = "\n" + problem
	}
	return fmt.Sprintf(`%s
从%s访问
地址：%s%s
耗时：%v
%s
`, t.At.Format(time.RFC3339), hostname, addr, problem,
		time.Duration(t.Duration)*time.Millisecond, strings.Join(t.Logs, "\n"),
	), nil
}

var emailClient *email.Client
var keepers []string

func setupConfig(configFile string) (Config, error) {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}
	var conf Config
	if err := yaml.Unmarshal(content, &conf); err != nil {
		return Config{}, err
	}
	if client, err := email.NewClient(conf.Mailer); err != nil {
		return Config{}, err
	} else {
		emailClient = client
	}
	keepers = conf.Keepers
	return conf, nil
}

func sendAlarm(content string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := emailClient.Send(ctx, &email.Message{
		Headers: []email.Header{
			{Name: email.To, Values: keepers},
			{Name: email.Subject, Values: []string{"URL监控"}},
			{Name: email.ContentType, Values: []string{"text/plain; charset=utf-8"}},
		},
		Body: []byte(content),
	})
	if err != nil {
		fmt.Errorf("发送报警邮件失败：%v\n", err)
	}
	return nil
}
