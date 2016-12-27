package godoc

const upstartConfig = `# vim: set ft=upstart:
start on runlevel [2345]
respawn
respawn limit 2 60

script
  GOPATH={{.GoPath}} {{.GodocBin}} -http={{.Servers.CurrentAppServer.GodocAddr}} -index_interval={{.Godoc.IndexInterval}}
end script
`
