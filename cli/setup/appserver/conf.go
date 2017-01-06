package appserver

const upstartConfig = `# vim: set ft=upstart:

env PATH={{.Path}}
chdir {{.App.Root}}
setuid {{.Deploy.User}}

start on {{.Servers.CurrentAppServer.AppStartOn}}
respawn
respawn limit 5 60

script
  set +e
  exec >>log/appserver.log 2>&1

  echo $(date -Iseconds) 'starting.'
  ./{{.App.Name}}
  exitStatus=$?
  echo $(date -Iseconds) "crashed with status: $exitStatus."
end script

post-start script
  set +e
  exec >>log/appserver.log 2>&1

  xiaomei setup wait-appserver
end script

post-stop script
  set +e
  exec >>log/appserver.log 2>&1

  echo $(date -Iseconds) 'stopped.
';
end script
`
