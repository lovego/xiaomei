package main

import (
	"testing"
)

func TestMakeConf(t *testing.T) {
	got := string(makeConf(`z_test.tmpl`, configData{
		ListenPort:  `8001`,
		SendfileOff: true,
	}))

	expect := `# vim: set ft=nginx:

server {
  listen {{ .ListenPort }} default_server;
  root /var/www/example;

  index index.html;

  location ~ \.html {
    add_header Cache-Control "must-revalidate";
  }

  location ~ \.(js|css|png|gif|jpg|svg|ico|woff|woff2|ttf|eot|map|json)$ {
    expires max;
  }
  sendfile off;
  access_log /var/log/nginx/example/web.log std;
  error_log  /var/log/nginx/example/web.err;
}
`
	if got != expect {
		t.Errorf("\nexpect:\n%sgot:\n%s", expect, got)
	}
}
