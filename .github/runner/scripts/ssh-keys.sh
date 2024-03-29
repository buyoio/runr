#!/usr/bin/env bash
set -euo pipefail

github::ssh-keys() {
  local owner="${1}"
  local team="${2}"

  # short version in `gh` cli
  # gh api "/orgs/${OWNER}/teams/${TEAM_NAME}/members" --jq '.[].login' |
  #   xargs -I@ gh api /users/@/keys --jq '.[].key'
  github::team-keys "${owner}" "${team}" | 
    github::user-keys

}

github::team-keys() {
  local owner="${1}"
  local team="${2}"
  curl -sXGET \
    -H "Accept: application/vnd.github.v3+json" \
    -H "Authorization: Bearer ${GITHUB_TOKEN}" \
    "https://api.github.com/orgs/${owner}/teams/${team}/members" |
    jq -r '.[].login'
}

# shellcheck disable=SC2120
github::user-keys() {
  local user="${1:-$(cat)}"
  curl -sXGET \
      -H "Accept: application/vnd.github.v3+json" \
      "https://api.github.com/users/${user}/keys" |
      jq -r '.[].key'
}

main() {
  github::ssh-keys "${OWNER?}" "${GITHUB_TEAM?}" |
    jq -Rrc '[.,inputs]'
}

[[ "${0}" != "${BASH_SOURCE[0]}" ]] || main "${@}"