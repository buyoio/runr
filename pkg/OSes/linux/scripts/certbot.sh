#!/usr/bin/env argsh
# shellcheck shell=bash disable=SC1091 disable=SC2034
set -euo pipefail

import system

certbot::ensure() {
  local domain
  local -a args=(
    "domain" "Domain name"
  )
  :args "Issue a certificate for a domain" "${@}"

  system::certbot
  certbot certonly \
    --standalone \
    --non-interactive \
    --no-redirect \
    --register-unsafely-without-email \
    --agree-tos \
    -d "${domain}"
}