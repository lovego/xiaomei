package setup

const upstartConfig = `# vim: set ft=upstart:

setuid {{.UserName}}
chdir {{.AppRoot}}

start on {{.AppStartOn}}
respawn
respawn limit 5 60

# stop the job if completely impossible to start.
# pre-start is the only process that can stop the job with exit code or stop.
pre-start script
  exec >>log/appserver.log 2>&1
  test -x {{.AppName}} || {
    echo $(date -Iseconds) 'starting canceled: no executable {{.AppName}}.'
    exit 1
  }
  test -f config/config.yml || {
    echo $(date -Iseconds) 'starting canceled: no config/config.yml.'
    exit 1
  }
  test -f config/env.yml || {
    echo $(date -Iseconds) 'starting canceled: no config/env.yml.'
    exit 1
  }
end script

script
  set +e
  exec >>log/appserver.log 2>&1
  echo $(date -Iseconds) 'starting.'
  ./{{.AppName}}
  exitStatus=$?
  echo $(date -Iseconds) "crashed with status: $exitStatus."
end script

post-start script
  exec >>log/appserver.log 2>&1

  # wait until the AppPort has been bound.
  i=0;
  while [ $i -lt {{.StartTimeout}} ]; do
    case $(status) in
    *' start/post-start, '* )
      lsof -itcp@{{.AppAddrPort}} > /dev/null && {
				echo $(date -Iseconds) 'started. ({{.AppAddrPort}})'; exit
      } ;;
    * )
      echo $(date -Iseconds) 'starting failed.'; exit ;;
    esac
    sleep 1
    i=$(( i + 1 ));
  done

  # kill process if timeout
  s=$(status)
  case "$s" in
  *' start/post-start, process '* )
    s=${s#* start/post-start, process *}
    pid=${s%%[!0-9]*}
    test -n "$pid" && kill -- -$pid # kill process group
    ;;
  esac
  echo $(date -Iseconds) 'starting timeout.';
end script

post-stop script
  exec >>log/appserver.log 2>&1
  echo $(date -Iseconds) 'stopped.
';
end script
`
