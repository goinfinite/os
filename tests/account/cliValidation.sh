#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

accountUsername="t${OS_TEST_RUN_ID}val"
missingUsername="t${OS_TEST_RUN_ID}missing"
accountPassword='abc123!abc'

echo "Feature ${OS_TEST_FEATURE}: account validation"

osCliCapture account create -u 'Invalid-User' -p "${accountPassword}"
assertCliStatus userError "create account with invalid username fails"
assertCliExitCode 64 "invalid username exit code"

osCliCapture account create -u "${accountUsername}" -p "${accountPassword}"
assertCliStatus created "create valid account"

osCliCapture account create -u "${accountUsername}" -p "${accountPassword}"
assertCliStatus infraError "duplicate account creation fails"
assertCliExitCode 69 "duplicate account exit code"

osCliCapture account get -n "${accountUsername}"
assertCliJq '.body.accounts | length == 1' "duplicate creation did not persist a second account"

osCliCapture account update -u "${missingUsername}" -p "${accountPassword}"
assertCliStatus infraError "update nonexistent account fails"
assertCliExitCode 69 "update nonexistent account exit code"

osCliCapture account delete -i 999999
assertCliStatus infraError "delete nonexistent account fails"
assertCliExitCode 69 "delete nonexistent account exit code"

osCliCapture account create -u "${accountUsername}"
assertCliExitCode 1 "create account without the required password flag fails"
assertEquals false "${cliIsJson}" "missing required flag does not return a JSON response"

osCliCapture account delete
assertCliExitCode 1 "delete account without the required id flag fails"

osCliCapture account delete-key -u 999999
assertCliExitCode 1 "delete key without the required key id flag fails"

finishTests
