main () {
  clone_code_if_should

  cd {{.DeployPath}} || exit 1

  reset_code_to_target_point

  killall -q godoc
  which godoc > /dev/null || install golang-go.tools
  GOPATH={{ .GOPATH }} nohup godoc -http=:1234 >/dev/null 2>&1 &
}


reset_code_to_target_point() {  
  git fetch origin -u {{.GitBranch}}:{{.GitBranch}} || exit 1
  git checkout -q {{.GitBranch}} || exit 1

  if git merge-base --is-ancestor {{.GitCommit}} {{.GitBranch}}; then
    git reset --hard {{.GitCommit}} || exit 1
  else
    echo commit {{.GitCommit}} is not in branch {{.GitBranch}}; exit 1
  fi
}

clone_code_if_should() {
  test -d {{.DeployPath}} ||
    sudo mkdir -p {{.DeployPath}} &&
    sudo chown -R $(id -un):$(id -gn) {{.DeployPath}} || exit 1

  if test ! -d {{.DeployPath}}/.git; then
    ssh-keygen -F {{.GitHost}} > /dev/null || ssh-keyscan -H {{.GitHost}} >> ~/.ssh/known_hosts
    git clone --depth=1 {{.GitAddr}} {{.DeployPath}} || exit 1
  fi
}

install() {
  # 超过3天没更新源
  if test -n "`find /var/lib/apt/periodic/update-success-stamp -mtime +2`"; then
    sudo apt-get update --fix-missing
  fi
  sudo apt-get install -y "$1"
}

main

