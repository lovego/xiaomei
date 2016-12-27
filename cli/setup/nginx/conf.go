package nginx

const defaultConfig = `# vim: set ft=nginx:

log_format {{ .Deploy.Name }} '$time_iso8601 $host'
	' $request_method $request_uri $content_length $remote_addr'
	' $status $body_bytes_sent $request_time'
	' $remote_addr "$http_referer" "$http_user_agent"';

upstream {{ .Deploy.Name }} {
{{- range .Servers.All -}}
  {{- if .HasTask "appserver" }}
    server {{ .AppAddr }};
  {{- end -}}
{{ end }}
}

server {
  charset utf-8;
  listen  80;
  server_name {{ .App.Domain }};
  root {{ .App.Root }}/public/;
  {{ if .Nfs }} sendfile off; {{ end }}

  location / {
    proxy_pass   http://{{ .Deploy.Name }};
		include proxy_params;
  }

  location ~ /.*\.html {
  }

  location ^~ /static/ {
		alias {{ .App.Root }}/public/;
    expires max;
  }

  location = /favicon.ico {
  }

  access_log {{ .App.Root }}/log/nginx.log {{ .Deploy.Name }};
  error_log  {{ .App.Root }}/log/nginx.err;
}
`
