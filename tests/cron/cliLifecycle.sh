#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

cronComment="os-test-${OS_TEST_RUN_ID}"
updatedComment="os-test-${OS_TEST_RUN_ID}-updated"

echo "Feature ${OS_TEST_FEATURE}: cron lifecycle"

osCliCapture cron create -s '* * * * *' -c 'echo os-test' -d "${cronComment}"
assertCliStatus created "create cron"
assertCliExitCode 0 "create cron exit code"

osCliCapture cron get
assertCliStatus success "read crons"
assertCliJq ".body.crons | map(select(.comment == \"${cronComment}\")) | length == 1" "created cron is persisted"
assertCliJq ".body.crons | map(select(.comment == \"${cronComment}\" and .schedule == \"* * * * *\")) | length == 1" "created cron schedule matches"
cronId="$(jq -r --arg comment "${cronComment}" '.body.crons[] | select(.comment == $comment) | .id' <<<"${cliOutput}")"
assertNotEmpty "${cronId}" "created cron exposes an id"

osCliCapture cron update -i "${cronId}" -s '*/10 * * * *' -c 'echo os-test-updated' -d "${updatedComment}"
assertCliStatus success "update cron"
assertCliExitCode 0 "update cron exit code"

osCliCapture cron get
assertCliJq ".body.crons | map(select(.id == ${cronId} and .schedule == \"*/10 * * * *\" and .comment == \"${updatedComment}\")) | length == 1" "updated cron state is persisted"

osCliCapture cron delete -i "${cronId}"
assertCliStatus success "delete cron by id"
assertCliExitCode 0 "delete cron exit code"

osCliCapture cron get
assertCliJq ".body.crons | map(select(.id == ${cronId})) | length == 0" "deleted cron is absent"

finishTests
