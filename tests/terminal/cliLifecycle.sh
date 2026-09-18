#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

accountUsername="t${OS_TEST_RUN_ID}trm"
accountPassword='abc123!abc'

echo "Feature ${OS_TEST_FEATURE}: terminal lifecycle"

osCliCapture account create -u "${accountUsername}" -p "${accountPassword}" --is-super-admin false
assertCliStatus created "create the terminal account"

osCliCapture account get -n "${accountUsername}"
assertCliStatus success "read the terminal account"
accountId="$(jq -r '.body.accounts[0].id' <<<"${cliOutput}")"
assertNotEmpty "${accountId}" "terminal account exposes an id"

osCliCapture terminal create -a "${accountId}" -w /app
assertCliStatus created "create terminal session"
assertCliExitCode 0 "create terminal session exit code"
terminalSessionId="$(jq -r '.body.id' <<<"${cliOutput}")"
assertNotEmpty "${terminalSessionId}" "created terminal session exposes an id"

osCliCapture terminal list
assertCliStatus success "list terminal sessions"
assertCliJq ".body.pagination.itemsTotal >= 1" "list reports the total item count"
assertCliJq ".body.terminalSessions | map(select(.id == \"${terminalSessionId}\")) | length == 1" "created session is listed"
assertCliJq ".body.terminalSessions | map(select(.id == \"${terminalSessionId}\" and .accountUsername == \"${accountUsername}\")) | length == 1" "session owner matches"

assertCommandSucceeds "tmux session exists" \
	osBash "su - ${accountUsername} -c 'tmux list-sessions -F \"#{session_name}\"' | grep -q 'os-managed-${terminalSessionId}'"

osCliCapture terminal delete -i "${terminalSessionId}"
assertCliStatus success "delete terminal session"
assertCliExitCode 0 "delete terminal session exit code"

osCliCapture terminal list
assertCliJq ".body.terminalSessions | map(select(.id == \"${terminalSessionId}\")) | length == 0" "deleted session is absent"

assertCommandFails "tmux session is gone" \
	osBash "su - ${accountUsername} -c 'tmux list-sessions -F \"#{session_name}\"' | grep -q 'os-managed-${terminalSessionId}'"

osCliCapture terminal create -a "${accountId}" -w /nonexistent-directory
assertCliStatus userError "create rejects a missing working directory"

osCliCapture terminal create -a "${accountId}"
assertCliStatus created "create terminal session without a working directory"
assertCliJq '.body.workingDir == "/app"' "omitted working directory defaults to /app"
defaultSessionId="$(jq -r '.body.id' <<<"${cliOutput}")"
osCliCapture terminal delete -i "${defaultSessionId}"
assertCliStatus success "delete the defaulted terminal session"

for sessionIndex in $(seq 1 30); do
	osCliCapture terminal create -a "${accountId}"
	assertCliStatus created "create session ${sessionIndex} up to the per-account cap"
done

osCliCapture terminal create -a "${accountId}"
assertCliStatus userError "create rejects the session past the per-account cap"

osCliCapture terminal list
for sessionId in $(jq -r '.body.terminalSessions[].id' <<<"${cliOutput}"); do
	osCliCapture terminal delete -i "${sessionId}"
done

finishTests
