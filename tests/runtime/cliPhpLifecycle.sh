#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

vhostHostname="t${OS_TEST_RUN_ID}.php.test"
serviceName="php-webserver"
phpVersion="8.3"
pageMarker="php-ok-${OS_TEST_RUN_ID}"

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

osCliCapture vhost mapping create -n "${vhostHostname}" -p / -t service -v "${serviceName}"
assertCliStatus created "create mapping to php-webserver"

osCliCapture runtime php update -n "${vhostHostname}" -v "${phpVersion}"
assertCliStatus success "update php version"
assertCliExitCode 0 "update php version exit code"

osCliCapture runtime php update-setting -n "${vhostHostname}" -v "${phpVersion}" -N memory_limit -V 512M
assertCliStatus success "update php setting"
assertCliExitCode 0 "update php setting exit code"

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

osCliCapture services delete -n "${serviceName}"
assertCliStatus success "delete php-webserver service"
waitForCliJq 300 "deleted php-webserver service is absent" \
	".body.installedServices | map(select(.name == \"${serviceName}\")) | length == 0" \
	services get -n "${serviceName}" -m 50

finishTests
