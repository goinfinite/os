#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

vhostHostname="t${OS_TEST_RUN_ID}.security.test"
ruleName="os-test-rule-${OS_TEST_RUN_ID}"

echo "Feature ${OS_TEST_FEATURE}: mapping security rule lifecycle"

osCliCapture vhost create -n "${vhostHostname}" -t top-level
assertCliStatus created "create virtual host for security rule"

osCliCapture vhost mapping security create -n "${ruleName}" -a '10.0.0.0/8' -S 100 -H 200
assertCliStatus created "create mapping security rule"
assertCliExitCode 0 "create security rule exit code"

osCliCapture vhost mapping security get -n "${ruleName}"
assertCliStatus success "read mapping security rules"
assertCliJq ".body.mappingSecurityRules | map(select(.name == \"${ruleName}\" and .rpsSoftLimitPerIp == 100 and .rpsHardLimitPerIp == 200)) | length == 1" "created security rule is persisted"
ruleId="$(jq -r --arg name "${ruleName}" '.body.mappingSecurityRules[] | select(.name == $name) | .id' <<<"${cliOutput}")"
assertNotEmpty "${ruleId}" "security rule exposes an id"

osCliCapture vhost mapping create -n "${vhostHostname}" -p / -t static-files -s "${ruleId}"
assertCliStatus created "create mapping bound to the security rule"

osCliCapture vhost mapping get -n "${vhostHostname}"
assertCliJq ".body.virtualHostWithMappings[0].mappings[0].mappingSecurityRuleId == ${ruleId}" "mapping persists the security rule reference"
mappingId="$(jq -r '.body.virtualHostWithMappings[0].mappings[0].id' <<<"${cliOutput}")"

osCliCapture vhost mapping security update -i "${ruleId}" -n "${ruleName}-updated" -S 50
assertCliStatus success "update mapping security rule"

osCliCapture vhost mapping security get -i "${ruleId}"
assertCliJq ".body.mappingSecurityRules[0].name == \"${ruleName}-updated\"" "updated security rule name is persisted"
assertCliJq '.body.mappingSecurityRules[0].rpsSoftLimitPerIp == 50' "updated security rule limit is persisted"

osCliCapture vhost mapping delete -i "${mappingId}"
assertCliStatus success "delete mapping before the rule"
osCliCapture vhost mapping security delete -i "${ruleId}"
assertCliStatus success "delete mapping security rule"

osCliCapture vhost mapping security get -i "${ruleId}"
assertCliJq '.body.mappingSecurityRules | length == 0' "deleted security rule is absent"

osCliCapture vhost delete -n "${vhostHostname}"
assertCliStatus success "delete virtual host after security tests"

finishTests
