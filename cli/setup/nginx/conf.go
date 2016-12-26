package nginx

const defaultConfig = `# vim: set ft=nginx:

log_format {{ .DeployName }} '$time_iso8601 $host'
	' $request_method $request_uri $content_length $remote_addr'
	' $status $body_bytes_sent $request_time'
	' $remote_addr "$http_referer" "$http_user_agent"';

upstream {{ .DeployName }} {
{{- range  .Servers -}}
  {{- if .HasTask "appserver" }}
    server {{ .AppAddr }};
  {{- end -}}
{{ end }}
}

server {
  charset utf-8;
  listen  80;
  server_name {{ .Domain }};
  root {{ .AppRoot }}/public/;
  {{ if .Nfs }} sendfile off; {{ end }}

  location / {
    proxy_pass   http://{{ .DeployName }};
		include proxy_params;
  }

  location ~ /.*\.html {
  }

  location ^~ /static/ {
		alias {{ .AppRoot }}/public/;
    expires max;
  }

  location = /favicon.ico {
  }

  access_log {{ .AppRoot }}/log/nginx.log {{ .DeployName }};
  error_log  {{ .AppRoot }}/log/nginx.err;
}
`
