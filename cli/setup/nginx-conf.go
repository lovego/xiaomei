package setup

const nginxConfig = `
# vim: set ft=nginx:

upstream {{ .DeployName }} {
{{- range  .DeployServers -}}
  {{- if .AppAddr }}
    server {{ .AppAddr }}:{{ $.AppPort }};
  {{- end -}}
{{ end }}
}

server {
  charset utf-8;

  server_name  {{ .Domain }};
  root {{ .AppRoot }}/public/;

  listen  80;
  include proxy_params;

  location / {
    proxy_pass   http://{{ .DeployName }};
  }

  location ~ ^/static/(.*)$ {
    try_files /$1 =404;
    expires max;
  }

  location = /favicon.ico {
  }

  {{ if .Nfs }} sendfile off; {{end}}
  access_log {{ .AppRoot }}/log/nginx.log;
  error_log  {{ .AppRoot }}/log/nginx.err;
}
`
