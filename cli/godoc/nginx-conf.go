package godoc

const nginxConfig = `# vim: set ft=nginx:

upstream {{ .Deploy.Name }}_godoc {
{{- range  .Servers -}}
  {{- if .HasTask "app" }}
    server {{ .ListenAddr }}:{{ $.Godoc.Port }};
  {{- end -}}
{{ end }}
}

server {
  charset utf-8;
  listen  80;
  server_name {{ .Godoc.Domain }}.{{.App.Domain}};

  location / {
    proxy_pass   http://{{ .DeployName }}_godoc;
		include proxy_params;
  }
}
`
