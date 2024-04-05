#!/usr/bin/env argsh
# shellcheck shell=bash disable=SC1090 disable=SC2034 disable=SC1091
set -euo pipefail

docker-cron() {
  :args "Setup docker system prune cron" "${@}"

  echo '%[1]s root /usr/bin/docker system prune -a --volumes --force' > /etc/cron.d/docker_system_prune
  /etc/init.d/cron reload
}

# shellcheck disable=SC2046
docker-stop() {
  :args "Stop runner container" "${@}"

  docker stop $(docker ps -qf "label=runr=runner") 2>/dev/null || :
  docker rm -f $(docker ps -qaf "label=runr=runner") 2>/dev/null || :
}

docker-build() {
  local dockerfile
  local -a args=(
    "dockerfile:~file" "Dockerfile"
  )
  :args "Build docker image" "${@}"

  docker build --no-cache -t github-runner:local - <"${dockerfile}"
}

docker-start() {
  local name orgrepo token labels path
  local -a args=(
    "name"    "Name of the runner"
    "orgrepo" "Organization or/and repository (org/repo)"
    "token"   "Runner token"
    "labels"  "Runner labels"
    "path"    "Runner work path"
  )
  :args "Run docker container" "${@}"

  path="${path}/${name}"
  local -a params
  params=(
    -v "${path}:${path}"
    -v /var/run/docker.sock:/var/run/docker.sock
    "-e" "RUNNER_NAME=${name}-$(hostname)"
    "-e" "RUNNER_TOKEN=${token}"
    "-e" "LABELS=${labels}"
    "-e" "RUNNER_WORKDIR=${path}"
    "-e" "DISABLE_AUTO_UPDATE=1"
    "-l" "runr=runner"
  )

  if echo "${orgrepo}" | grep -q '/'; then
    params+=(
      "-e" "RUNNER_SCOPE=repo"
      "-e" "REPO_URL=https://github.com/${orgrepo}"
    )
  else
    params+=(
      "-e" "RUNNER_SCOPE=org"
      "-e" "ORG_NAME=${orgrepo}"
    )
  fi

  mkdir -p "${path}"
  docker run -d --restart=always --name "${name}" \
    --add-host host.docker.internal:host-gateway \
    "${params[@]}" \
    github-runner:local
}