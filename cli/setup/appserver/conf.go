package appserver

const upstartConfig = `# vim: set ft=upstart:

setuid {{.UserName}}
chdir {{.AppRoot}}

start on {{.AppStartOn}}
respawn
respawn limit 5 60

script
  set +e
  exec >>log/appserver.log 2>&1

  echo $(date -Iseconds) 'starting.'
  ./{{.AppName}}
  exitStatus=$?
  echo $(date -Iseconds) "crashed with status: $exitStatus."
end script

post-start script
  set +e
  exec >>log/appserver.log 2>&1

  ./{{.AppName}} wait-appserver
end script

post-stop script
  set +e
  exec >>log/appserver.log 2>&1

  echo $(date -Iseconds) 'stopped.
';
end script
`
