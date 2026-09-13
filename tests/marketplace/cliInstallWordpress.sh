#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

vhostHostname="t${OS_TEST_RUN_ID}.wp.test"
adminUsername="os-test-${OS_TEST_RUN_ID}"
adminPassword='abc123!abc'
adminMailAddress="os-test-${OS_TEST_RUN_ID}@example.com"

echo "Feature ${OS_TEST_FEATURE}: wordpress marketplace install lifecycle"

osCliCapture vhost create -n "${vhostHostname}" -t top-level
assertCliStatus created "create virtual host for wordpress"
assertCliExitCode 0 "create wordpress virtual host exit code"

osCliCapture mktplace install -s wordpress -n "${vhostHostname}" \
	-f "adminUsername:${adminUsername}" \
	-f "adminPassword:${adminPassword}" \
	-f "adminMailAddress:${adminMailAddress}"
assertCliStatus created "install wordpress from the marketplace"
assertCliExitCode 0 "install wordpress exit code"

waitForCliJq 900 "installed wordpress item appears in the list" \
	".body.marketplaceInstalledItems[] | select(.hostname == \"${vhostHostname}\")" \
	mktplace list -m 50

osCliCapture mktplace list -m 50
assertCliJq ".body.marketplaceInstalledItems | map(select(.hostname == \"${vhostHostname}\" and .name == \"WordPress\")) | length == 1" "installed item persists the hostname and name"
installedId="$(jq -r --arg hostname "${vhostHostname}" '.body.marketplaceInstalledItems[] | select(.hostname == $hostname) | .id' <<<"${cliOutput}")"
assertNotEmpty "${installedId}" "installed wordpress item exposes an id"

waitForCliJq 300 "php-webserver service reaches the running state" \
	'.body.installedServices[] | select((.name | startswith("php-webserver")) and .status == "running")' \
	services get -n php-webserver -m 50

httpCode="$(curl -sL -o /dev/null -w '%{http_code}' --max-time 15 \
	-H "Host: ${vhostHostname}" \
	"http://127.0.0.1:${OS_TEST_HTTP_PORT}/wp-login.php")"
assertEquals 200 "${httpCode}" "wordpress answers over HTTP"

servedContent="$(curl -sL --max-time 15 -H "Host: ${vhostHostname}" \
	"http://127.0.0.1:${OS_TEST_HTTP_PORT}/wp-login.php")"
assertContains "${servedContent}" "WordPress" "wordpress content over HTTP"

httpsCode="$(curl -sLk -o /dev/null -w '%{http_code}' --max-time 15 \
	--proto '=https' --proto-redir '=https' \
	--resolve "${vhostHostname}:${OS_TEST_HTTPS_PORT}:127.0.0.1" \
	"https://${vhostHostname}:${OS_TEST_HTTPS_PORT}/wp-login.php")"
assertEquals 200 "${httpsCode}" "wordpress answers over HTTPS"

servedContentHttps="$(curl -sLk --max-time 15 \
	--proto '=https' --proto-redir '=https' \
	--resolve "${vhostHostname}:${OS_TEST_HTTPS_PORT}:127.0.0.1" \
	"https://${vhostHostname}:${OS_TEST_HTTPS_PORT}/wp-login.php")"
assertContains "${servedContentHttps}" "WordPress" "wordpress content over HTTPS"

osCliCapture mktplace delete -i "${installedId}"
assertCliStatus success "delete installed wordpress item"
assertCliExitCode 0 "delete installed wordpress item exit code"

waitForCliJq 300 "deleted wordpress item is absent" \
	".body.marketplaceInstalledItems | map(select(.hostname == \"${vhostHostname}\")) | length == 0" \
	mktplace list -m 50

osCliCapture vhost delete -n "${vhostHostname}"
assertCliStatus success "delete wordpress virtual host"
assertCliExitCode 0 "delete wordpress virtual host exit code"

finishTests
