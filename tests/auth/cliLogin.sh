#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

accountUsername="t${OS_TEST_RUN_ID}auth"
accountPassword='abc123!abc'

echo "Feature ${OS_TEST_FEATURE}: login"

osCliCapture account create -u "${accountUsername}" -p "${accountPassword}"
assertCliStatus created "account created for login tests"

osCliCapture auth login -u "${accountUsername}" -p "${accountPassword}" -i 127.0.0.1
assertCliStatus success "login with valid credentials"
assertCliExitCode 0 "valid login exit code"
assertCliJq '.body.type == "sessionToken"' "login returns a session token type"
assertCliJq '.body.tokenStr | length > 0' "login returns a token string"

tokenExpiry="$(jq -r '.body.expiresIn' <<<"${cliOutput}")"
assertNumberGreaterThan "$(date +%s)" "${tokenExpiry}" "token expiry is in the future"

osCliCapture auth login -u "${accountUsername}" -p 'wrong-pass-1!' -i 127.0.0.1
assertCliStatus infraError "login with wrong password fails"
assertCliExitCode 69 "wrong password exit code"

osCliCapture auth login -u "${accountUsername}" -p "${accountPassword}" -i 'not-an-ip'
assertCliStatus userError "login with malformed ip address fails"
assertCliExitCode 64 "malformed ip exit code"

osCliCapture auth login -u "${accountUsername}" -p "${accountPassword}"
assertCliExitCode 1 "login without the required ip flag fails"
assertEquals false "${cliIsJson}" "missing required flag does not return a JSON response"

finishTests
