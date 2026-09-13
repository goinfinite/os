#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

vhostHostname="t${OS_TEST_RUN_ID}.node.test"

echo "Feature ${OS_TEST_FEATURE}: node marketplace install lifecycle"

osCliCapture vhost create -n "${vhostHostname}" -t top-level
assertCliStatus created "create virtual host for node app"
assertCliExitCode 0 "create node app virtual host exit code"

osCliCapture mktplace install -s uptime-kuma -n "${vhostHostname}"
assertCliStatus created "install uptime kuma from the marketplace"
assertCliExitCode 0 "install uptime kuma exit code"

waitForCliJq 900 "installed uptime kuma item appears in the list" \
	".body.marketplaceInstalledItems[] | select(.hostname == \"${vhostHostname}\")" \
	mktplace list -m 50

osCliCapture mktplace list -m 50
assertCliJq ".body.marketplaceInstalledItems | map(select(.hostname == \"${vhostHostname}\" and .name == \"Uptime Kuma\")) | length == 1" "installed item persists the hostname and name"
installedId="$(jq -r --arg hostname "${vhostHostname}" '.body.marketplaceInstalledItems[] | select(.hostname == $hostname) | .id' <<<"${cliOutput}")"
assertNotEmpty "${installedId}" "installed uptime kuma item exposes an id"

waitForCliJq 300 "node service reaches the running state" \
	'.body.installedServices[] | select((.name | startswith("node")) and .status == "running")' \
	services get -n node -m 50

httpCode="$(curl -sL -o /dev/null -w '%{http_code}' --max-time 15 -H "Host: ${vhostHostname}" "http://127.0.0.1:${OS_TEST_HTTP_PORT}/")"
assertEquals 200 "${httpCode}" "uptime kuma answers over its mapping"

servedContent="$(curl -sL --max-time 15 -H "Host: ${vhostHostname}" "http://127.0.0.1:${OS_TEST_HTTP_PORT}/")"
assertContains "${servedContent}" "Uptime Kuma" "uptime kuma page identifies the application"

osCliCapture mktplace delete -i "${installedId}"
assertCliStatus success "delete installed uptime kuma item"
assertCliExitCode 0 "delete installed uptime kuma item exit code"

waitForCliJq 300 "deleted uptime kuma item is absent" \
	".body.marketplaceInstalledItems | map(select(.hostname == \"${vhostHostname}\")) | length == 0" \
	mktplace list -m 50

osCliCapture vhost delete -n "${vhostHostname}"
assertCliStatus success "delete uptime kuma virtual host"
assertCliExitCode 0 "delete uptime kuma virtual host exit code"

finishTests
