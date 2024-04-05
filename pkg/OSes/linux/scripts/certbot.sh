#!/usr/bin/env argsh
# shellcheck shell=bash disable=SC1091 disable=SC2034
set -euo pipefail

import system

certbot-ensure() {
  local domain email
  local -a args=(
    "domain" "Domain name"
    "email|" "Email for Let's Encrypt certificate"
  )
  :args "Issue a certificate for a domain" "${@}"

  local -a e=()
  if [[ -n "${email}" ]]; then
    e+=("--email" "${email}")
  else
    e+=("--register-unsafely-without-email")
  fi
  system-certbot

  if [[ -f "/etc/letsencrypt/live/${domain}/fullchain.pem" ]]; then
    certbot renew
    return
  fi

  certbot certonly \
    --standalone \
    --non-interactive \
    --no-redirect \
    --agree-tos \
    -d "${domain}" "${e[@]}"
}