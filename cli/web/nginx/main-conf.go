package nginx

const defaultMainConfig = `# vim: set ft=nginx:
user www-data;
pid /run/nginx.pid;

worker_processes auto;
worker_rlimit_nofile 100000;

events {
	worker_connections 100000;
  use epoll;
  multi_accept on;
}

http {
	# Basic Settings
	sendfile on;
	tcp_nopush on;
	tcp_nodelay on;
	keepalive_timeout 65;
	types_hash_max_size 2048;
  client_header_buffer_size 4k;

	include /etc/nginx/mime.types;
	default_type application/octet-stream;

	# Gzip Settings
	gzip on;
	gzip_disable "msie6";

	# Logging Settings
	log_format std '$time_iso8601 $host'
  ' $request_method $request_uri $content_length $server_protocol'
	' $status $body_bytes_sent $request_time'
	' $remote_addr "$http_referer" "$http_user_agent"';

	access_log /var/log/nginx/access.log;
	error_log /var/log/nginx/error.log;

	# Virtual Host Configs
	include /etc/nginx/conf.d/*.conf;
	include /etc/nginx/sites-enabled/*;
}
`
