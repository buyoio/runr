#!/usr/bin/env argsh
# shellcheck shell=bash disable=SC1091
set -euo pipefail

hetzner::rescue-mode() {
  [ -x /root/.oldroot/nfs/install/installimage ] || {
      return 127
  }
}

hetzner::install() {
  local installimage
  # shellcheck disable=SC2034
  local -a args=(
    "installimage:~stdin" "config (use - for stdin)"
  )
  :args "Prevision a server with a Hetzner installimage" "${@}"

  hetzner::rescue-mode || {
    echo "Server is not in rescue mode" >&2
    return 1 
  }
  [ "${installimage}" != "-" ] || installimage="$(cat)"

  echo "${installimage}" > /tmp/installimage
  /root/.oldroot/nfs/install/installimage -a -c /tmp/installimage
  reboot
}