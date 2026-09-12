#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

vhostHostname="t${OS_TEST_RUN_ID}.ssl.test"
certificatePath="/tmp/os-test-ssl.crt"
privateKeyPath="/tmp/os-test-ssl.key"

echo "Feature ${OS_TEST_FEATURE}: ssl pair lifecycle"

osCliCapture vhost create -n "${vhostHostname}" -t top-level
assertCliStatus created "create virtual host for ssl"

osBash "openssl req -x509 -newkey rsa:2048 -nodes -keyout ${privateKeyPath} -out ${certificatePath} -days 1 -subj '/CN=${vhostHostname}' 2>/dev/null"
assertEquals 0 "$?" "generate a self-signed certificate pair"

osCliCapture ssl create -v "${vhostHostname}" -c "${certificatePath}" -k "${privateKeyPath}"
assertCliStatus created "create ssl pair"
assertCliExitCode 0 "create ssl pair exit code"

osCliCapture ssl get -n "${vhostHostname}"
assertCliStatus success "read ssl pairs"
assertCliJq ".body.sslPairs | map(select(.virtualHostHostname == \"${vhostHostname}\")) | length == 1" "created ssl pair is persisted"
sslPairId="$(jq -r --arg hostname "${vhostHostname}" '.body.sslPairs[] | select(.virtualHostHostname == $hostname) | .sslPairId' <<<"${cliOutput}")"
assertNotEmpty "${sslPairId}" "ssl pair exposes an id"

osCliCapture ssl delete -i "${sslPairId}"
assertCliStatus success "delete ssl pair"
assertCliExitCode 0 "delete ssl pair exit code"

osCliCapture ssl get -n "${vhostHostname}"
assertCliJq ".body.sslPairs | map(select(.sslPairId == \"${sslPairId}\")) | length == 0" "deleted ssl pair is absent"

osBash "printf 'not a certificate' > /tmp/os-test-bad.crt; printf 'not a key' > /tmp/os-test-bad.key"
osCliCapture ssl create -v "${vhostHostname}" -c /tmp/os-test-bad.crt -k /tmp/os-test-bad.key
assertCliStatus userError "create ssl pair with invalid certificate content fails"
assertCliExitCode 64 "invalid certificate exit code"

osCliCapture ssl delete -i '0000000000000000000000000000000000000000000000000000000000000000'
assertCliStatus infraError "delete nonexistent ssl pair fails"
assertCliExitCode 69 "delete nonexistent ssl pair exit code"

osCliCapture vhost delete -n "${vhostHostname}"
assertCliStatus success "delete virtual host after ssl tests"

finishTests
