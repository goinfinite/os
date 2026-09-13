#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

serviceName="os-test-svc-${OS_TEST_RUN_ID}"

echo "Feature ${OS_TEST_FEATURE}: custom service lifecycle"

osCliCapture services create-custom -n "${serviceName}" -t other -c '/bin/sleep infinity' -w /tmp -a false
assertCliStatus created "create custom service"
assertCliExitCode 0 "create custom service exit code"

waitForCliJq 60 "custom service reaches the running state" \
	".body.installedServices[] | select(.name == \"${serviceName}\" and .status == \"running\")" \
	services get -n "${serviceName}" -m 50

osCliCapture services get -n "${serviceName}" -m 50
assertCliJq ".body.installedServices | map(select(.name == \"${serviceName}\" and .type == \"other\")) | length == 1" "persisted service type matches"

osCliCapture services update -n "${serviceName}" -s restart
assertCliStatus success "restart custom service"
assertCliExitCode 0 "restart custom service exit code"
waitForCliJq 60 "restarted service returns to the running state" \
	".body.installedServices[] | select(.name == \"${serviceName}\" and .status == \"running\")" \
	services get -n "${serviceName}" -m 50

osCliCapture services delete -n "${serviceName}"
assertCliStatus success "delete custom service"
assertCliExitCode 0 "delete custom service exit code"
waitForCliJq 60 "deleted service is absent" \
	".body.installedServices | map(select(.name == \"${serviceName}\")) | length == 0" \
	services get -n "${serviceName}" -m 50

osCliCapture services create-custom -n "${serviceName}" -t other -c '/bin/sleep infinity' -w /tmp -a false
assertCliStatus created "recreate custom service"
waitForCliJq 60 "recreated service reaches the running state" \
	".body.installedServices[] | select(.name == \"${serviceName}\" and .status == \"running\")" \
	services get -n "${serviceName}" -m 50
osCliCapture services create-custom -n "${serviceName}" -t other -c '/bin/sleep infinity' -w /tmp -a false
assertCliStatus infraError "duplicate custom service creation fails"
assertCliExitCode 69 "duplicate custom service exit code"
osCliCapture services delete -n "${serviceName}"
assertCliStatus success "delete recreated custom service"

osCliCapture services create-custom -n 'Invalid Service Name' -t other -c '/bin/sleep infinity' -w /tmp -a false
assertCliStatus userError "create custom service with invalid name fails"
assertCliExitCode 64 "invalid service name exit code"

osCliCapture services update -n 'os-test-missing' -s restart
assertCliStatus infraError "restart nonexistent service fails"
assertCliExitCode 69 "restart nonexistent service exit code"

osCliCapture services delete -n 'os-test-missing'
assertCliStatus infraError "delete nonexistent service fails"
assertCliExitCode 69 "delete nonexistent service exit code"

osCliCapture services get-installables -n redis
assertCliStatus success "read installable services"
assertCliJq '.body.installableServices | length == 1' "installable filter returns exactly redis"
assertCliJq '.body.installableServices[0].versions | length > 0' "redis exposes installable versions"

finishTests
