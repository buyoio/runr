#!/usr/bin/env argsh
# shellcheck shell=bash disable=SC1091 disable=SC2034
set -euo pipefail

import github
import system

harbor::ensure() {
  local domain path="/opt"
  local -a args=(
    "domain" "Domain name"
    "path"   "Path to install"
  )
  :args "Install or update Harbor" "${@}"

  system::certbot
  system::docker
  system::docker-compose
  certbot::ensure "${domain}"

  mkdir -p "${path}"
  pushd "${path}"
  trap 'popd' EXIT RETURN

  local -r v="$(github::latest goharbor/harbor)"
  ! grep -qw "${v}" .harbor/version &>/dev/null || return 0

  curl -sL -o harbor.tgz "https://github.com/goharbor/harbor/releases/download/${v}/harbor-online-installer-${v}.tgz"
  tar xzf --overwrite harbor.tgz
  rm harbor.tgz

  local -r passwd_admin="$(string::random)"
  local -r passwd_db="$(string::random)"
  [ ! -f db ] || {
    passwd_db="$(cat db)"
  }
  echo "${passwd_db}" > db

  pushd harbor
  trap 'popd; popd' EXIT RETURN

  cp harbor.yml.tmpl harbor.yml
  tld="$(echo "${domain}" | rev | cut -d'.' -f1,2 | rev)"
  sed -i \
    -e '/hostname:/ s%: .*%: '"${domain}"'%' \
    -e '/  certificate:/ s%: .*%: /etc/letsencrypt/live/'"${tld}"'/fullchain.pem%' \
    -e '/  private_key:/ s%: .*%: /etc/letsencrypt/live/'"${tld}"'/privkey.pem%' \
    -e '/harbor_admin_password:/ s%: .*%: '"${passwd_admin}"'%' \
    -e '/  password:/ s%: .*%: '"${passwd_db}"'%' \
    harbor.yml

  echo "${v}" > version
  ./install.sh
}