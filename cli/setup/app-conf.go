package setup

const upstartConfig = `
# vim: set ft=upstart:

setuid {{.UserName}}
chdir {{.AppRoot}}

start on {{.AppStartOn}}
respawn
respawn limit 3 10

pre-start script
  test -f ./bin/appserver || { stop; exit 1; }
end script

script 
  exec ./bin/appserver >>log/app.log 2>>log/app.err
end script

post-start script
  bash >>log/app.err 2>&1 <<'EOF'
    waited=0
    while true; do
      if (( waited >= 120 )); then
        echo 'app server starting timeout'; stop; exit 1
      elif status $JOB | fgrep 'stop/post-start,'; then
        echo 'app server starting failed'; exit 1
      elif lsof -i4:{{.AppPort}} > /dev/null; then
        break
      fi
      (( w = waited > 0 ? waited : 1 ))
      sleep $w
      (( waited += w ))
    done
EOF
end script
`
