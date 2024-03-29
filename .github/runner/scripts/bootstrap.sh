#!/usr/bin/env bash
set -euo pipefail

direnv::install() {
  echo "> install direnv"
  direnv_path="/usr/local/bin/direnv"

  curl -o "${direnv_path}" -sfL "$(
    curl -sfL https://api.github.com/repos/direnv/direnv/releases/latest |
      grep browser_download_url |
      cut -d '"' -f 4 |
      grep "direnv.linux.amd64"
  )"
  chmod +x "${direnv_path}"
}

direnv::hook() {
  # shellcheck disable=SC2016
  echo 'eval "$(direnv hook bash)"' >>/home/user/.bashrc
}

main() {
  direnv::install
  direnv::hook
}

[[ "${0}" != "${BASH_SOURCE[0]}" ]] || main "${@}"