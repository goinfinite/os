#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

serviceName="os-test-svc-${OS_TEST_RUN_ID}"

echo "Feature ${OS_TEST_FEATURE}: installable service lifecycle"

osCliCapture services get-installables -n python
assertCliStatus success "read installable catalog"
assertCliJq '.body.installableServices | map(select(.name == "python" and (.versions | index("3.12")))) | length == 1' "catalog lists python 3.12"

osCliCapture services create-installable -n python -v 3.12
assertCliStatus created "install python 3.12"
assertCliExitCode 0 "install python exit code"

waitForCliJq 300 "installed python service appears" \
	'.body.installedServices | length == 1' \
	services get -n python -m 50

osCliCapture services get -n python -m 50
assertCliJq '.body.installedServices | length == 1' "exactly one python service is installed"
assertCliJq '.body.installedServices[0].version == "3.12"' "installed python version matches"
installedName="$(jq -r '.body.installedServices[0].name' <<<"${cliOutput}")"
assertNotEmpty "${installedName}" "installed python service exposes a name"

osCliCapture services create-installable -n python -v 9.99
assertCliJq '.status != "success"' "installing a python version outside the catalog fails"

osCliCapture services create-installable -n no-such-service -v 1
assertCliJq '.status != "success"' "installing an unknown service fails"

osCliCapture services delete -n "${installedName}"
assertCliStatus success "delete installed python service"
assertCliExitCode 0 "delete installed python exit code"

waitForCliJq 120 "deleted python service is absent" \
	'.body.installedServices | length == 0' \
	services get -n python -m 50

echo "Feature ${OS_TEST_FEATURE}: custom service lifecycle"

osCliCapture services create-custom -n "${serviceName}" -t other -c '/bin/sleep infinity' -w /tmp -a false
assertCliStatus created "create custom service"
assertCliExitCode 0 "create custom service exit code"

waitForCliJq 60 "custom service reaches the running state" \
	".body.installedServices[] | select(.name == \"${serviceName}\" and .status == \"running\")" \
	services get -n "${serviceName}" -m 50

osCliCapture services update -n "${serviceName}" -s restart
assertCliStatus success "restart custom service"
waitForCliJq 60 "restarted service returns to the running state" \
	".body.installedServices[] | select(.name == \"${serviceName}\" and .status == \"running\")" \
	services get -n "${serviceName}" -m 50

osCliCapture services delete -n "${serviceName}"
assertCliStatus success "delete custom service"
waitForCliJq 60 "deleted custom service is absent" \
	".body.installedServices | map(select(.name == \"${serviceName}\")) | length == 0" \
	services get -n "${serviceName}" -m 50

finishTests
