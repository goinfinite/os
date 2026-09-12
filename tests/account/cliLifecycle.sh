#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

accountUsername="t${OS_TEST_RUN_ID}acc"
keeperUsername="t${OS_TEST_RUN_ID}keep"
accountPassword='abc123!abc'
updatedPassword='abc123!def'
publicKeyContent="$(cat "${OS_TEST_REPO_ROOT}/tests/fixtures/ed25519PublicKey.pub")"

echo "Feature ${OS_TEST_FEATURE}: account lifecycle"

osCliCapture account get -n "${accountUsername}"
assertCliStatus success "read accounts before creation"
assertCliJq '.body.accounts | length == 0' "account does not exist before creation"

osCliCapture account create -u "${accountUsername}" -p "${accountPassword}" --is-super-admin false
assertCliStatus created "create account"
assertCliExitCode 0 "create account exit code"

osCliCapture account get -n "${accountUsername}"
assertCliStatus success "read created account"
assertCliJq '.body.accounts | length == 1' "exactly one account matches the username"
assertCliJq ".body.accounts[0].username == \"${accountUsername}\"" "persisted username matches"
assertCliJq '.body.accounts[0].isSuperAdmin == false' "persisted isSuperAdmin matches"
accountId="$(jq -r '.body.accounts[0].id' <<<"${cliOutput}")"
assertNotEmpty "${accountId}" "created account exposes an id"

osCliCapture account update -u "${accountUsername}" -p "${updatedPassword}"
assertCliStatus success "update account password"
assertCliExitCode 0 "update account exit code"

osCliCapture auth login -u "${accountUsername}" -p "${updatedPassword}" -i 127.0.0.1
assertCliStatus success "login with the updated password"
assertCliJq '.body.tokenStr | length > 0' "login returns a session token"

osCliCapture account create-public-key -u "${accountId}" -n "os-test-key" -c "${publicKeyContent}"
assertCliStatus created "create secure access public key"

osCliCapture account get -n "${accountUsername}" -s true
assertCliJq '.body.accounts[0].secureAccessPublicKeys | length == 1' "public key is persisted on the account"
keyId="$(jq -r '.body.accounts[0].secureAccessPublicKeys[0].id' <<<"${cliOutput}")"
assertNotEmpty "${keyId}" "persisted public key exposes an id"

osCliCapture account delete-key -u "${accountId}" -i "${keyId}"
assertCliStatus success "delete secure access public key"
osCliCapture account get -n "${accountUsername}" -s true
assertCliJq '.body.accounts[0].secureAccessPublicKeys | length == 0' "deleted public key is absent"

osCliCapture account create -u "${keeperUsername}" -p "${accountPassword}"
assertCliStatus created "create keeper account"

osCliCapture account delete -i "${accountId}"
assertCliStatus success "delete account"
assertCliExitCode 0 "delete account exit code"

osCliCapture account get -n "${accountUsername}"
assertCliJq '.body.accounts | length == 0' "deleted account is absent"

osCliCapture account get -n "${keeperUsername}"
assertCliJq '.body.accounts | length == 1' "keeper account remains"

finishTests
