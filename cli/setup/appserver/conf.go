package appserver

const upstartConfig = `# vim: set ft=upstart:

setuid {{.Deploy.User}}
chdir {{.App.Root}}

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

  ./{{.App.Name}} setup wait-appserver
end script

post-stop script
  set +e
  exec >>log/appserver.log 2>&1

  echo $(date -Iseconds) 'stopped.
';
end script
`
