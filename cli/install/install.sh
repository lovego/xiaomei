#!/bin/bash

main() {
  setup_sudo_no_password
  which git          || install git
  which nginx        || install nginx-core
  which redis-server || install redis-server
  install_mysql_server
  install_mongodb_shell
  install_docker
}

setup_sudo_no_password() {
  local line='ubuntu  ALL=NOPASSWD: ALL'
  local regexp="$(echo "$line" | sed -r 's/[\t ]+/[[:space:]]+/g')"
  sudo egrep -q "^$regexp$" /etc/sudoers || echo "$line" | sudo tee --append /etc/sudoers > /dev/null
}


install_mysql_server() {
  which mysqld && return
  sudo debconf-set-selections <<< "mysql-server-5.6 mysql-server/root_password password root"
  sudo debconf-set-selections <<< "mysql-server-5.6 mysql-server/root_password_again password root"
  install mysql-server-5.6
}

install_mongodb_shell() {
  which mongo && return
  sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv EA312927
  echo "deb http://repo.mongodb.org/apt/ubuntu trusty/mongodb-org/3.2 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-3.2.list
  sudo apt-get update
  sudo apt-get install -y mongodb-org-shell
}

install_docker() {
  which docker && return
  sudo apt-key adv --keyserver hkp://ha.pool.sks-keyservers.net:80 \
    --recv-keys 58118E89F3A912897C070ADBF76221572C52609D

  echo "deb https://apt.dockerproject.org/repo ubuntu-trusty main" | sudo tee /etc/apt/sources.list.d/docker.list
  sudo apt-get update
  sudo apt-get install -y linux-image-extra-$(uname -r) linux-image-extra-virtual docker-engine 
  sudo usermod -aG docker $(id -un)
}

install() {
  # 超过3天没更新源
  if test -n "`find /var/lib/apt/periodic/update-success-stamp -mtime +2`"; then
    sudo apt-get update --fix-missing
  fi
  sudo apt-get install -y "$1"
}

main
