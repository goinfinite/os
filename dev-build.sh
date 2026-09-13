#!/bin/bash
#
# @description  Runs a disposable Infinite OS container and rebuilds its binary on every code change.
# @usage        bash dev-build.sh [http|http-unpriv|ols|no-cache] [--pid-only]
# @output       The container log stream on stdout and stderr.
# @requires     bash v4+, curl, make, podman, air
# @version      0.2.1
# @updated      2026-09-08
#
# No companion dev-build_test.sh on purpose: the script is a thin orchestrator
# over make, podman, air, and curl, and the repository runs its unit tests
# inside the container, never on the host.
#
set -euo pipefail

containerName="os"
imageName="os:latest"
devHostName="goinfinite.app"
devApiPort="1618"
devAccountUser="dev"
devAccountPassword="abc123!"
apiServeWaitSeconds=120
validBuildModes=(http http-unpriv ols no-cache)
buildMode=""
pidFilePath=""
watcherPid=""
attachPid=""
stopAlreadyRan=false

#
## CommandLine
#
cliArgumentsResolver() {
  local resolvedBuildMode=""
  local resolvedPidFilePath=""

  for arg in "$@"; do
    if [[ "${arg}" == "--pid-only" ]]; then
      resolvedPidFilePath="logs/dev-build.pid"
      continue
    fi

    argIsKnown="false"
    for mode in "${validBuildModes[@]}"; do
      if [[ "${arg}" == "${mode}" ]]; then
        argIsKnown="true"
      fi
    done

    if [[ "${argIsKnown}" == "false" ]]; then
      resolvedBuildMode="rejectedArg:${arg}"
      break
    fi

    if [[ -z "${resolvedBuildMode}" ]]; then
      resolvedBuildMode="${arg}"
    fi
  done

  echo "${resolvedBuildMode}|${resolvedPidFilePath}"
}

publishProcessId() {
  if [[ -z "${pidFilePath}" ]]; then
    return 0
  fi

  mkdir -p "$(dirname "${pidFilePath}")"
  echo $$ > "${pidFilePath}"
}

removePublishedProcessId() {
  if [[ -z "${pidFilePath}" ]]; then
    return 0
  fi

  rm -f "${pidFilePath}"
}

#
## PortPublishing
#
exposedPortsBuilder() {
  local ports="-p 127.0.0.1:${devApiPort}:${devApiPort}"

  case ${1} in
  http)
    ports="${ports} -p 127.0.0.1:80:80 -p 127.0.0.1:443:443"
    ;;
  http-unpriv)
    ports="${ports} -p 127.0.0.1:8080:80 -p 127.0.0.1:8443:443"
    ;;
  ols)
    ports="${ports} -p 127.0.0.1:7080:7080"
    ;;
  esac

  echo "${ports}"
}

enablePrivilegedPortBinding() {
  sudo sysctl net.ipv4.ip_unprivileged_port_start=80
}

#
## Output
#
reportDevBuildStarted() {
  cat <<EOF

<<>><<>><<>><<>><<>><<>><<>><<>><<>><<>><<>><<>><<>><<>><<>>

=> Starting the development build...
=> Any changes to the code trigger a rebuild automatically.
=> Dashboard: https://localhost:${devApiPort}

<<>><<>><<>><<>><<>><<>><<>><<>><<>><<>><<>><<>><<>><<>><<>>

EOF
}

#
## ContainerLifecycle
#
clearContainerImageCache() {
  podman image prune -a
  podman rmi "localhost/${imageName}" -f || true
}

buildDevImage() {
  make build
  podman build -t "${imageName}" --format docker .
}

startDevContainer() {
  # shellcheck disable=SC2086
  podman run --name "${containerName}" -d \
    --env "LOG_LEVEL=debug" --env "PRIMARY_VHOST=${devHostName}" \
    --env "DEV_MODE=true" \
    --hostname="${devHostName}" --cpus=2 --memory=2g --rm \
    --volume "$(pwd)/bin:/infinite/bin:Z,ro,bind,slave" \
    ${1} "${imageName}"
}

waitForApiToServe() {
  local deadline=$((SECONDS + apiServeWaitSeconds))

  while [[ ${SECONDS} -lt ${deadline} ]]; do
    if curl -ks -o /dev/null --max-time 2 "https://127.0.0.1:${devApiPort}/"; then
      return 0
    fi

    sleep 1
  done

  echo "=> The development API did not serve on port ${devApiPort} within ${apiServeWaitSeconds} seconds." >&2
  exit 1
}

replaceBinaryWithDevelopmentBuild() {
  podman exec "${containerName}" /bin/bash -c \
    'rm -f os && ln -s bin/os os && supervisorctl restart os-api'
}

createDevAccount() {
  podman exec "${containerName}" /bin/bash -c \
    "os account create -u ${devAccountUser} -p ${devAccountPassword} --is-super-admin false"
}

stopProcess() {
  if [[ -z "${1}" ]]; then
    return 0
  fi

  kill "${1}" 2>/dev/null || true
}

stopDevBuildSession() {
  if [[ "${stopAlreadyRan}" == "true" ]]; then
    return 0
  fi
  stopAlreadyRan=true

  stopProcess "${watcherPid}"
  stopProcess "${attachPid}"
  podman stop "${containerName}" &>/dev/null || true
  podman rm "${containerName}" &>/dev/null || true
  removePublishedProcessId
}

handleStopSignal() {
  stopDevBuildSession
  echo "=> Development build stopped."
  exit 0
}

#
## Watchers
#
# Air also starts the application on its own, which would run a second instance on
# the host. SILENT_EXIT_MODE keeps it as a rebuild trigger only.
startRebuildWatcher() {
  SILENT_EXIT_MODE=true air
}

streamContainerLogs() {
  podman attach "${containerName}"
}

#
## Runtime
#
if [[ "${BASH_SOURCE[0]}" != "$0" ]]; then
  return 0
fi

IFS='|' read -r buildMode pidFilePath <<< "$(cliArgumentsResolver "$@")"

if [[ "${buildMode}" == rejectedArg:* ]]; then
  echo "=> Unknown argument '${buildMode#rejectedArg:}'. Use ${validBuildModes[*]}, or --pid-only." >&2
  exit 1
fi

publishProcessId

if [[ "${buildMode}" == "no-cache" ]]; then
  clearContainerImageCache
fi

if [[ "${buildMode}" == "http" ]]; then
  enablePrivilegedPortBinding
fi

publishedPorts="$(exposedPortsBuilder "${buildMode}")"
buildDevImage
startDevContainer "${publishedPorts}"
trap handleStopSignal SIGINT SIGTERM
trap stopDevBuildSession EXIT
waitForApiToServe
replaceBinaryWithDevelopmentBuild
waitForApiToServe
createDevAccount

reportDevBuildStarted
startRebuildWatcher &
watcherPid=$!
streamContainerLogs &
attachPid=$!
wait "${watcherPid}"
