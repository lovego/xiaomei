package deploy

const updateCodeShell = `# vim: set ft=sh:
main () {
  clone_code_if_should

  cd {{.DeployPath}}/release || exit 1

  reset_code_to_target_point

  mkdir -p bins
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

reset_code_to_target_point() {
  git fetch origin -u --tags {{.GitBranch}}:{{.GitBranch}} || exit 1
  git checkout -q {{.GitBranch}} || exit 1

  if git merge-base --is-ancestor {{.GitTag}} {{.GitBranch}}; then
    git reset --hard {{.GitTag}} || exit 1
  else
    echo tag {{.GitTag}} is not in branch {{.GitBranch}}; exit 1
  fi
}

main
`

const setupShell = `# vim: set ft=sh:
main () {
  cd {{.DeployPath}}/release || exit 1

  ln -sfT envs/{{.Env}}.yml config/env.yml || exit 1
  ln -sfT bins/{{ .GitTag }} {{ .AppName }} || exit 1
  ./{{ .AppName }} setup {{ .Tasks }} || exit 1

  clear_local_obsolete_deploy_tags
}

clear_local_obsolete_deploy_tags() {
  for local_tag in $(git tag -l {{.Env}}'*'); do
    git ls-remote --tags --exit-code origin "$local_tag" >/dev/null || git tag -d "$local_tag"
  done
}

main
`
