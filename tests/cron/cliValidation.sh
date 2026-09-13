#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

missingComment="os-test-missing-${OS_TEST_RUN_ID}"

echo "Feature ${OS_TEST_FEATURE}: cron validation"

osCliCapture cron create -s 'not a cron schedule' -c 'echo os-test' -d "${missingComment}"
assertCliStatus userError "create cron with invalid schedule fails"
assertCliExitCode 64 "invalid schedule exit code"

osCliCapture cron update -i 999999 -s '* * * * *'
assertCliStatus infraError "update nonexistent cron fails"
assertCliExitCode 69 "update nonexistent cron exit code"

osCliCapture cron delete -d "${missingComment}"
assertCliStatus infraError "delete nonexistent cron by comment fails"
assertCliExitCode 69 "delete nonexistent cron exit code"

osCliCapture cron delete -i 999999
assertCliStatus infraError "delete nonexistent cron by id fails"
assertCliExitCode 69 "delete nonexistent cron by id exit code"

osCliCapture cron create -s '* * * * *'
assertCliExitCode 1 "create cron without the required command flag fails"
assertEquals false "${cliIsJson}" "missing required flag does not return a JSON response"

osCliCapture cron update -s '* * * * *'
assertCliExitCode 1 "update cron without the required id flag fails"

finishTests
