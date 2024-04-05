#!/usr/bin/env argsh
# shellcheck shell=bash disable=SC1091 disable=SC2034
set -euo pipefail

# import docker
# import hetzner
# import harbor

lok8s() {
  bash::version 4 3 || {
    echo "Bash 4.3 or later is required" >&2
    exit 2 
  }
  local -a usage=(
    'system-ensure' "Ensure system dependencies"
    'system-install' "Install system dependencies"
    'system-user' "Create system user"
    'system-docker' "Install Docker"
    'system-docker-compose' "Install Docker Compose"
    'system-certbot' "Install Certbot"
    'system-reboot' "Reboot system (if required)"
    'certbot-ensure' "Issue or renew SSL certificate"
    'docker-ensure' "Ensure docker"
    'docker-build' "Build runner docker image"
    'docker-start' "Start runner docker container"
    'docker-stop' "Stop runner docker container"
    'docker-cron' "Setup docker cron"
    'hetzner-install' "Install Hetzner image"
    'harbor-ensure' "Ensure Harbor"
  )
  :usage "Preparse linux system for runr." "${@}"
  # system-non-interactive
  # system-root
  "${usage[@]}"
}

[[ "${0}" != "${BASH_SOURCE[0]}" ]] || lok8s "${@}"