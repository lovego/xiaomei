#!/bin/bash

# 1. 使用VirtualBox安装Ubuntu16.04 初始GuestOS
# 2. 添加Host-only网卡
# 3. 设置22端口转发，使用终端登录GuestOS
# 4. 运行该脚本

set -e

main() {
  setup_sudo_no_password
  setup_vbox_hostonly_network
}

setup_sudo_no_password() {
  username=$(id -nu)
  local file="/etc/sudoers.d/$username"
  test -f "$file"  || echo "$username  ALL=NOPASSWD: ALL" | sudo tee --append "$file" > /dev/null
}

setup_vbox_hostonly_network() {
  file=/etc/network/interfaces.d/host-only
  test -f $file && return
  # ls /sys/class/net
  echo '
auto enp0s8
iface enp0s8 inet static
address 192.168.56.15
netmask 255.255.255.0
' | sudo tee --append $file > /dev/null
  sudo ifdown enp0s8 && sudo ifup enp0s8
}

setup_vbox_guest_addtions() {
  sudo apt-get install -y gcc make perl
  sudo rcvboxadd setup
  # 需要先挂载好VBoxGuestAdditions.iso
  sudo mount -t auto /dev/cdrom /media/cdrom
  sudo /media/cdrom/VBoxLinuxAdditions.run
}

setup_vbox_share_folder() {
  # 需要先设置好共享文件夹
  # auto mount by vbox
  sudo usermod -a -G vboxsf $(id -nu)
}

main
