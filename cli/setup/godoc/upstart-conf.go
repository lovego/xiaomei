package godoc

const upstartConfig = `# vim: set ft=upstart:
start on runlevel [2345]
respawn
respawn limit 2 60

script
  GOPATH={{.GoPath}} {{.GodocBin}} -http={{.Addr}} {{if .IndexInterval -}}
	-index_interval={{.IndexInterval}}
	{{- end }}
end script
`
