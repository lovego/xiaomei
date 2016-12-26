package godoc

const nginxConfig = `# vim: set ft=nginx:

upstream {{ .Deploy.Name }}_godoc {
{{- range .Servers.All -}}
  {{- if .HasTask "appserver" }}
    server {{ .GodocAddr }};
  {{- end -}}
{{ end }}
}

server {
  charset utf-8;
  listen  80;
  server_name {{ .Godoc.Domain }}.{{ .App.Domain }};

  location / {
    proxy_pass   http://{{ .Deploy.Name }}_godoc;
		include proxy_params;
  }
}
`
