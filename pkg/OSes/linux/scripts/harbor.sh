#!/usr/bin/env argsh
# shellcheck shell=bash disable=SC1091 disable=SC2034
set -euo pipefail

import github
import system

harbor-ensure() {
  local domain email path="/opt"
  local -a args=(
    "domain" "Domain name"
    "path|"  "Path to install"
    "email|" "Email for Let's Encrypt certificate"
  )
  :args "Install or update Harbor" "${@}"

  system-certbot
  system-docker-compose
  certbot-ensure "${domain}" --email "${email}"

  mkdir -p "${path}"
  cd "${path}"

  local -r v="$(github::latest goharbor/harbor)"
  ! grep -qw "${v}" harbor/version &>/dev/null || return 0

  curl -sL -o harbor.tgz "https://github.com/goharbor/harbor/releases/download/${v}/harbor-online-installer-${v}.tgz"
  tar -xzf harbor.tgz
  rm harbor.tgz
  cd harbor

  local passwd_admin passwd_db
  if [[ -f db ]]; then
    passwd_db="$(cat ./db)"
    passwd_admin="$(cat ./admin)"
  else
    passwd_db="$(string::random)"
    passwd_admin="$(string::random)"
  fi
  echo -n "${passwd_db}" > ./db
  echo -n "${passwd_admin}" > ./admin

  cp harbor.yml.tmpl harbor.yml
  sed -i \
    -e '/hostname:/ s%: .*%: '"${domain}"'%' \
    -e '/  certificate:/ s%: .*%: /etc/letsencrypt/live/'"${domain}"'/fullchain.pem%' \
    -e '/  private_key:/ s%: .*%: /etc/letsencrypt/live/'"${domain}"'/privkey.pem%' \
    -e '/harbor_admin_password:/ s%: .*%: '"${passwd_admin}"'%' \
    -e '/  password:/ s%: .*%: '"${passwd_db}"'%' \
    harbor.yml

  echo "${v}" > version
  ./install.sh
}