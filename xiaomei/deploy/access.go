package deploy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/cluster"
	"github.com/lovego/xiaomei/xiaomei/deploy/conf"
	"github.com/lovego/xiaomei/xiaomei/release"
)

func accessPrint(svcName string) error {
	accessConf, err := getAccessConf(svcName)
	if err != nil {
		return err
	}
	fmt.Print(accessConf)
	return nil
}

func accessSetup(svcName string) error {
	domain := getDomainName(svcName)
	script := fmt.Sprintf(`
	sudo tee /etc/nginx/sites-enabled/%s.conf > /dev/null &&
	sudo mkdir -p /var/log/nginx/%s &&
	sudo nginx -t &&
	sudo service nginx reload
	`, domain, domain,
	)
	accessConf, err := getAccessConf(svcName)
	if err != nil {
		return err
	}
	for _, node := range cluster.Nodes() {
		if node.Labels[`access`] == `true` {
			if _, err := node.Run(
				cmd.O{Stdin: strings.NewReader(accessConf)}, script,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func getAccessConf(svcName string) (string, error) {
	var confTmpl string
	if buf, err := ioutil.ReadFile(
		filepath.Join(release.Root(), `access`, `access.conf.tmpl`),
	); err == nil {
		confTmpl = string(buf)
	} else {
		return ``, err
	}
	tmpl := template.Must(template.New(``).Parse(confTmpl))
	configData, err := getConfigData(svcName)
	if err != nil {
		return ``, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, configData); err != nil {
		return ``, err
	}
	return buf.String(), nil
}

func getConfigData(svcName string) (interface{}, error) {
	data := struct {
		Env, SvcName, DomainName string
		BackendAddrs             []string
	}{
		Env:        release.Env(),
		SvcName:    svcName,
		DomainName: getDomainName(svcName),
	}
	if svcName == `` {
		svcName = getServiceToAccess()
	}
	if addrs, err := getDriver().Addrs(svcName); err == nil {
		data.BackendAddrs = addrs
		return data, nil
	} else {
		return nil, err
	}
}

func getDomainName(svcName string) string {
	domain := release.App().Domain()
	if svcName != `` {
		domain = getSvcDomain(domain, svcName)
	}
	return domain
}

func getSvcDomain(domain, svcName string) string {
	parts := strings.SplitN(domain, `.`, 2)
	if len(parts) == 2 {
		return parts[0] + `-` + svcName + `.` + parts[1]
	} else {
		return domain + `-` + svcName
	}
}

func getServiceToAccess() string {
	services := conf.ServiceNames()
	if services[`web`] {
		return `web`
	}
	if services[`app`] {
		return `app`
	}
	panic(`no backend service found in ` + conf.File() + `.`)
}
