#!/bin/bash

main() {
  setup_sudo_no_password
  which git          || install git
  which nginx        || install nginx-core
  which redis-server || install redis-server
  install_mysql_server
  install_mongodb_shell
}

setup_sudo_no_password() {
  local line='ubuntu  ALL=NOPASSWD: ALL'
  local regexp="$(echo "$line" | sed -r 's/[\t ]+/[[:space:]]+/g')"
  sudo egrep -q "^$regexp$" /etc/sudoers || echo "$line" | sudo tee --append /etc/sudoers > /dev/null
}


install_mysql_server() {
  if ! which mysqld; then
    sudo debconf-set-selections <<< "mysql-server-5.6 mysql-server/root_password password root"
    sudo debconf-set-selections <<< "mysql-server-5.6 mysql-server/root_password_again password root"
    install mysql-server-5.6
  fi
}

install_mongodb_shell() {
  if ! which mongo; then
    sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv EA312927
    echo "deb http://repo.mongodb.org/apt/ubuntu trusty/mongodb-org/3.2 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-3.2.list
    sudo apt-get update
    sudo apt-get install -y mongodb-org-shell
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
