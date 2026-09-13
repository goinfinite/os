#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

accountUsername="t${OS_TEST_RUN_ID}race"
accountPassword='abc123!abc'
firstPassword='abc123!first'
secondPassword='abc123!second'
concurrency=5
exitCodesFile="$(mktemp)"

echo "Feature ${OS_TEST_FEATURE}: concurrent account operations"

for _ in $(seq 1 "${concurrency}"); do
	(
		osCliRun account create -u "${accountUsername}" -p "${accountPassword}" >/dev/null 2>&1
		printf '%s\n' "$?" >>"${exitCodesFile}"
	) &
done
wait

successCount="$(grep -c '^0$' "${exitCodesFile}" || true)"
assertEquals 1 "${successCount}" "exactly one of ${concurrency} concurrent creates succeeds"

osCliCapture account get -n "${accountUsername}"
assertCliStatus success "read account after concurrent creates"
assertCliJq '.body.accounts | length == 1' "exactly one account persisted after concurrent creates"

for password in "${firstPassword}" "${secondPassword}"; do
	(
		osCliRun account update -u "${accountUsername}" -p "${password}" >/dev/null 2>&1
	) &
done
wait

successfulLogins=0
for password in "${firstPassword}" "${secondPassword}"; do
	osCliCapture auth login -u "${accountUsername}" -p "${password}" -i 127.0.0.1
	if [[ "${cliStatus}" == "success" ]]; then
		successfulLogins=$((successfulLogins + 1))
	fi
done
assertEquals 1 "${successfulLogins}" "exactly one concurrent password update wins"

osCliCapture account get -n "${accountUsername}"
assertCliJq '.body.accounts | length == 1' "account survives concurrent updates"

rm -f "${exitCodesFile}"
finishTests
