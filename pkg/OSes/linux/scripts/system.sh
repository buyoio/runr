#!/usr/bin/env argsh
# shellcheck shell=bash disable=SC1091 disable=SC2190 disable=SC2034
set -euo pipefail

import github

system::non-interactive() {
  export DEBIAN_FRONTEND=noninteractive
  echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections
}

system::root() {
  [ "${EUID}" = 0 ] || sudo -i
}

system::reboot() {
  [ ! -e /var/run/reboot-required ] || reboot
}

system::ensure() {
  local -a packages
  local -A args=(
    'packages' "System dependencies"
  )
  :args "Ensure system dependencies" "${@}"

  apt-get update
  apt-get -o Dpkg::Options::="--force-confold" upgrade -q -y
  apt-get install -y "${packages[@]}"

  system::reboot
}

system::docker() {
  local user force
  local -a args=(
    'user'      "User to add to docker group"
    'force|f:+' "Force install"
  )
  :args "Install Docker" "${@}"
  
  ! command -v docker || (( force )) || return 0
  sh -c "$(curl -sL https://get.docker.com)"
  usermod -aG docker "${user}"
}

system::docker-compose() {
  local -a args=(
    'force|f:+' "Force install"
  )
  :args "Install Docker Compose" "${@}"
  
  ! command -v docker-compose || (( force )) || return 0
  local -r v="$(github::latest docker/compose)"
  curl -s \
    -L "https://github.com/docker/compose/releases/download/${v}/docker-compose-$(uname -s)-$(uname -m)" \
    -o /usr/local/bin/docker-compose
  chmod +x /usr/local/bin/docker-compose
}

system::certbot() {
  :args "Install Certbot" "${@}"
  ! command -v certbot || return 0

  system::ensure curl zip python3 python3-pip
  python3 -m pip install certbot
}