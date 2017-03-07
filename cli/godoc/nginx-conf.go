package godoc

const nginxConfig = `# vim: set ft=nginx:

upstream godoc {
{{- range .Addrs }}
  server {{ . }};
{{ end -}}
}

server {
  charset utf-8;
  listen  80;
  server_name {{ .Domain }};

  location / {
    proxy_pass   http://godoc;
    include proxy_params;
  }
}
`
