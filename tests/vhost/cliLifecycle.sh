#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

vhostHostname="t${OS_TEST_RUN_ID}.vhost.test"
secondHostname="t${OS_TEST_RUN_ID}.second.test"
pageContent="os-test-${OS_TEST_RUN_ID}"

echo "Feature ${OS_TEST_FEATURE}: virtual host and mapping lifecycle"

osCliCapture vhost create -n "${vhostHostname}" -t top-level
assertCliStatus created "create virtual host"
assertCliExitCode 0 "create virtual host exit code"

osCliCapture vhost get -n "${vhostHostname}"
assertCliStatus success "read virtual host"
assertCliJq ".body.virtualHosts | length == 1" "exactly one virtual host matches the hostname"
assertCliJq ".body.virtualHosts[0].type == \"top-level\"" "persisted type matches"
assertCliJq ".body.virtualHosts[0].rootDirectory == \"/app/html/${vhostHostname}\"" "persisted root directory matches"

osCliCapture vhost mapping create -n "${vhostHostname}" -p / -t static-files
assertCliStatus created "create static-files mapping"
assertCliExitCode 0 "create mapping exit code"

osCliCapture vhost mapping get -n "${vhostHostname}"
assertCliStatus success "read mapping"
assertCliJq '.body.virtualHostWithMappings[0].mappings | length == 1' "exactly one mapping matches the hostname"
assertCliJq '.body.virtualHostWithMappings[0].mappings[0].path == "/"' "persisted mapping path matches"
assertCliJq '.body.virtualHostWithMappings[0].mappings[0].targetType == "static-files"' "persisted target type matches"
mappingId="$(jq -r '.body.virtualHostWithMappings[0].mappings[0].id' <<<"${cliOutput}")"
assertNotEmpty "${mappingId}" "mapping exposes an id"

osBash "mkdir -p /app/html/${vhostHostname} && printf '%s' '${pageContent}' > /app/html/${vhostHostname}/index.html"
servedContent="$(curl -s --max-time 5 -H "Host: ${vhostHostname}" "http://127.0.0.1:${OS_TEST_HTTP_PORT}/")"
assertEquals "${pageContent}" "${servedContent}" "mapping serves the virtual host content over HTTP"

osCliCapture vhost mapping delete -i "${mappingId}"
assertCliStatus success "delete mapping"
assertCliExitCode 0 "delete mapping exit code"

osCliCapture vhost mapping get -n "${vhostHostname}"
assertCliJq '.body.virtualHostWithMappings[0].mappings | length == 0' "deleted mapping is absent"

osCliCapture vhost delete -n "${vhostHostname}"
assertCliStatus success "delete virtual host"
assertCliExitCode 0 "delete virtual host exit code"

osCliCapture vhost get -n "${vhostHostname}"
assertCliJq '.body.virtualHosts | length == 0' "deleted virtual host is absent"

osCliCapture vhost create -n "${secondHostname}" -t top-level
assertCliStatus created "create second virtual host"
osCliCapture vhost create -n "${secondHostname}" -t top-level
assertCliJq '.status != "success"' "duplicate virtual host creation fails"
osCliCapture vhost get -n "${secondHostname}"
assertCliJq '.body.virtualHosts | length == 1' "duplicate creation did not persist a second host"
osCliCapture vhost delete -n "${secondHostname}" >/dev/null

osCliCapture vhost create -n 'not a hostname' -t top-level
assertCliStatus userError "create virtual host with malformed hostname fails"
assertCliExitCode 64 "malformed hostname exit code"

osCliCapture vhost mapping create -n 'missing.vhost.test' -p / -t static-files
assertCliExitCode 69 "mapping to a nonexistent virtual host fails"
osCliCapture vhost mapping delete -i 999999
assertCliExitCode 69 "deleting a nonexistent mapping fails"

finishTests
