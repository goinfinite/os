#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

vhostHostname="t${OS_TEST_RUN_ID}.passbolt.test"
adminMailAddress="os-test-${OS_TEST_RUN_ID}@example.com"
phpVersion="8.3"

echo "Feature ${OS_TEST_FEATURE}: passbolt marketplace install lifecycle"

osCliCapture vhost create -n "${vhostHostname}" -t top-level
assertCliStatus created "create virtual host for passbolt"
assertCliExitCode 0 "create passbolt virtual host exit code"

osCliCapture mktplace install -s passbolt -n "${vhostHostname}" \
	-f "adminFirstName:OS" \
	-f "adminLastName:Test" \
	-f "adminMailAddress:${adminMailAddress}" \
	-f "smtpFrom:${adminMailAddress}" \
	-f "smtpHost:smtp.example.com" \
	-f "smtpPort:587" \
	-f "smtpUsername:os-test" \
	-f "smtpPassword:abc123!abc"
assertCliStatus created "install passbolt from the marketplace"
assertCliExitCode 0 "install passbolt exit code"

waitForCliJq 900 "installed passbolt item appears in the list" \
	".body.marketplaceInstalledItems[] | select(.hostname == \"${vhostHostname}\")" \
	mktplace list -m 50

osCliCapture mktplace list -m 50
assertCliJq ".body.marketplaceInstalledItems | map(select(.hostname == \"${vhostHostname}\" and .name == \"Passbolt\")) | length == 1" "installed item persists the hostname and name"
installedId="$(jq -r --arg hostname "${vhostHostname}" '.body.marketplaceInstalledItems[] | select(.hostname == $hostname) | .id' <<<"${cliOutput}")"
assertNotEmpty "${installedId}" "installed passbolt item exposes an id"

waitForCliJq 300 "php-webserver service reaches the running state" \
	'.body.installedServices[] | select((.name | startswith("php-webserver")) and .status == "running")' \
	services get -n php-webserver -m 50

osCliCapture runtime php get -n "${vhostHostname}"
assertCliStatus success "read php configs after passbolt install"
assertCliJq ".body.modules | map(select(.name == \"pear\" and .status == true)) | length == 1" "passbolt install leaves the pear tool module enabled"

httpsCode="$(curl -sk -o /dev/null -w '%{http_code}' --max-time 20 \
	--resolve "${vhostHostname}:${OS_TEST_HTTPS_PORT}:127.0.0.1" \
	"https://${vhostHostname}:${OS_TEST_HTTPS_PORT}/auth/login")"
assertEquals 200 "${httpsCode}" "passbolt answers over HTTPS"

servedContent="$(curl -sk --max-time 20 \
	--resolve "${vhostHostname}:${OS_TEST_HTTPS_PORT}:127.0.0.1" \
	"https://${vhostHostname}:${OS_TEST_HTTPS_PORT}/auth/login")"
assertContains "${servedContent}" "Passbolt" "passbolt page identifies the application"

osCliCapture runtime php update-module -n "${vhostHostname}" -N pear -V false -v "${phpVersion}"
assertCliStatus success "disable the pear tool module"
assertCliJq '.body.modulesSuccessfullyUpdated | map(select(.name == "pear" and .status == false)) | length == 1' "disable reports pear as updated"
assertCliJq '.body.failedModulesWithReason | length == 0' "disable reports no failed modules"

osCliCapture runtime php get -n "${vhostHostname}"
assertCliJq ".body.modules | map(select(.name == \"pear\" and .status == false)) | length == 1" "disabled pear status is reflected"

osCliCapture runtime php update-module -n "${vhostHostname}" -N pear -V true -v "${phpVersion}"
assertCliStatus success "enable the pear tool module"
assertCliJq '.body.modulesSuccessfullyUpdated | map(select(.name == "pear" and .status == true)) | length == 1' "enable reports pear as updated"
assertCliJq '.body.failedModulesWithReason | length == 0' "enable reports no failed modules"

osCliCapture runtime php get -n "${vhostHostname}"
assertCliJq ".body.modules | map(select(.name == \"pear\" and .status == true)) | length == 1" "re-enabled pear status is reflected"

osCliCapture mktplace delete -i "${installedId}"
assertCliStatus success "delete installed passbolt item"
assertCliExitCode 0 "delete installed passbolt item exit code"

waitForCliJq 300 "deleted passbolt item is absent" \
	".body.marketplaceInstalledItems | map(select(.hostname == \"${vhostHostname}\")) | length == 0" \
	mktplace list -m 50

osCliCapture vhost delete -n "${vhostHostname}"
assertCliStatus success "delete passbolt virtual host"
assertCliExitCode 0 "delete passbolt virtual host exit code"

finishTests
