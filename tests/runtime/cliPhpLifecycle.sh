#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

vhostHostname="t${OS_TEST_RUN_ID}.php.test"
serviceName="php-webserver"
phpVersion="8.3"
pageMarker="php-ok-${OS_TEST_RUN_ID}"
lsConfPath="/usr/local/lsws/conf/httpd_config.conf"

assertLsapiCapacityIsTuned() {
	local lsphpBlockCount
	lsphpBlockCount="$(osBash "grep -c '^extprocessor lsphp' ${lsConfPath}")"

	assertCommandSucceeds "httpd worker count is tuned to one on a two cpu container" \
		osBash "grep -q '^httpdworkers 1; # AUTO CALCULATED' ${lsConfPath}"
	assertCommandSucceeds "lsphp maxConns is tuned to ten on a two GiB container" \
		osBash "test \"\$(grep -c 'maxConns 10; # AUTO CALCULATED' ${lsConfPath})\" -eq ${lsphpBlockCount}"
	assertCommandSucceeds "lsphp child process capacity is tuned to ten on a two GiB container" \
		osBash "test \"\$(grep -c 'PHP_LSAPI_CHILDREN=10; # AUTO CALCULATED' ${lsConfPath})\" -eq ${lsphpBlockCount}"
}

echo "Feature ${OS_TEST_FEATURE}: php runtime lifecycle"

osCliCapture vhost create -n "${vhostHostname}" -t top-level
assertCliStatus created "create virtual host for php"

osCliCapture runtime php get -n "${vhostHostname}"
assertCliStatus serviceUnavailable "php config read fails while php-webserver is absent"
assertCliExitCode 69 "php config read exit code without php-webserver"

osCliCapture services create-installable -n "${serviceName}" -v "${phpVersion}"
assertCliStatus created "install php-webserver"
assertCliExitCode 0 "install php-webserver exit code"

waitForCliJq 900 "php-webserver reaches the running state" \
	".body.installedServices[] | select(.name == \"${serviceName}\" and .status == \"running\")" \
	services get -n "${serviceName}" -m 50

assertCommandSucceeds "restart the container to apply the boot capacity update" \
	"${OS_TEST_RUNTIME}" restart "${OS_TEST_CONTAINER}"
waitForCliJq 300 "php-webserver returns to the running state after the restart" \
	".body.installedServices[] | select(.name == \"${serviceName}\" and .status == \"running\")" \
	services get -n "${serviceName}" -m 50

assertLsapiCapacityIsTuned

osCliCaptureWithOperation "vhost.mapping.create-php" \
	vhost mapping create -n "${vhostHostname}" -p / -t service -v "${serviceName}"
assertCliStatus created "create mapping to php-webserver"
assertLsapiCapacityIsTuned

osCliCapture runtime php update -n "${vhostHostname}" -v "${phpVersion}"
assertCliStatus success "update php version"
assertCliExitCode 0 "update php version exit code"
assertLsapiCapacityIsTuned

osCliCapture runtime php update-setting -n "${vhostHostname}" -v "${phpVersion}" -N memory_limit -V 512M
assertCliStatus success "update php setting"
assertCliExitCode 0 "update php setting exit code"
assertLsapiCapacityIsTuned

osCliCapture runtime php get -n "${vhostHostname}"
assertCliStatus success "read php configs"
assertCliJq ".body.hostname == \"${vhostHostname}\"" "php configs report the virtual host"
assertCliJq ".body.version.value == \"${phpVersion}\"" "php configs report the selected version"
assertCliJq '.body.settings | map(select(.name == "memory_limit" and .value == "512M")) | length == 1' "updated php setting is persisted"

osBash "printf '<?php echo \"${pageMarker}\";' > /app/html/${vhostHostname}/index.php"
servedContent="$(curl -s --max-time 10 -H "Host: ${vhostHostname}" "http://127.0.0.1:${OS_TEST_HTTP_PORT}/")"
assertEquals "${pageMarker}" "${servedContent}" "php virtual host executes php and serves the result"

osCliCapture vhost mapping get -n "${vhostHostname}"
mappingId="$(jq -r '.body.virtualHostWithMappings[0].mappings[0].id' <<<"${cliOutput}")"

osCliCapture vhost mapping delete -i "${mappingId}"
assertCliStatus success "delete php mapping"
osCliCapture vhost delete -n "${vhostHostname}"
assertCliStatus success "delete php virtual host"

assertCommandSucceeds "deleted php virtual host block is removed" \
	osBash "! grep -q 'virtualhost ${vhostHostname}' ${lsConfPath}"
assertCommandSucceeds "deleted php virtual host map is removed" \
	osBash "! grep -q 'map.*${vhostHostname}' ${lsConfPath}"
assertCommandSucceeds "deleted php virtual host conf file is removed" \
	osBash "test ! -f /app/conf/php-webserver/${vhostHostname}.conf"
assertLsapiCapacityIsTuned

osCliCapture services delete -n "${serviceName}"
assertCliStatus success "delete php-webserver service"
waitForCliJq 300 "deleted php-webserver service is absent" \
	".body.installedServices | map(select(.name == \"${serviceName}\")) | length == 0" \
	services get -n "${serviceName}" -m 50

finishTests
