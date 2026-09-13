#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

vhostHostname="t${OS_TEST_RUN_ID}.mkt.test"
adminEmail="os-test-${OS_TEST_RUN_ID}@example.com"
adminPassword='abc123!abc'

echo "Feature ${OS_TEST_FEATURE}: marketplace install lifecycle"

osCliCapture vhost create -n "${vhostHostname}" -t top-level
assertCliStatus created "create virtual host for marketplace install"

osCliCapture mktplace install -s pocketbase -n "${vhostHostname}" \
	-f "adminEmail:${adminEmail}" -f "adminPassword:${adminPassword}"
assertCliStatus created "install pocketbase from the marketplace"
assertCliExitCode 0 "install pocketbase exit code"

waitForCliJq 600 "installed marketplace item appears in the list" \
	".body.marketplaceInstalledItems[] | select(.hostname == \"${vhostHostname}\")" \
	mktplace list -m 50

osCliCapture mktplace list -m 50
assertCliJq ".body.marketplaceInstalledItems | map(select(.hostname == \"${vhostHostname}\" and .name == \"PocketBase\")) | length == 1" "installed item persists the hostname and name"
installedId="$(jq -r --arg hostname "${vhostHostname}" '.body.marketplaceInstalledItems[] | select(.hostname == $hostname) | .id' <<<"${cliOutput}")"
assertNotEmpty "${installedId}" "installed item exposes an id"

waitForCliJq 300 "pocketbase service reaches the running state" \
	'.body.installedServices[] | select((.name | startswith("pocketbase")) and .status == "running")' \
	services get -n pocketbase -m 50

httpCode="$(curl -s -o /dev/null -w '%{http_code}' --max-time 10 -H "Host: ${vhostHostname}" "http://127.0.0.1:${OS_TEST_HTTP_PORT}/api/health")"
assertEquals 200 "${httpCode}" "installed application answers over its mapping"

osCliCapture mktplace delete -i "${installedId}"
assertCliStatus success "delete installed marketplace item"
assertCliExitCode 0 "delete installed item exit code"

waitForCliJq 300 "deleted marketplace item is absent" \
	".body.marketplaceInstalledItems | map(select(.hostname == \"${vhostHostname}\")) | length == 0" \
	mktplace list -m 50

osCliCapture vhost delete -n "${vhostHostname}"
assertCliStatus success "delete marketplace virtual host"

osCliCapture mktplace install -s no-such-slug
assertCliJq '.status != "success"' "install unknown marketplace slug fails"

finishTests
