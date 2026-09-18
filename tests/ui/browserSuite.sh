#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

browserDir="${OS_TEST_REPO_ROOT}/tests/ui/browser"
accountUsername="t${OS_TEST_RUN_ID}web"
accountPassword='abc123!abc'
browserArtifactsDir="${OS_TEST_ARTIFACTS_DIR}/browser"

echo "Feature ${OS_TEST_FEATURE}: browser suite"

osCliCapture account create -u "${accountUsername}" -p "${accountPassword}" --is-super-admin false
assertCliStatus created "create the browser account"

otherAccountUsername="t${OS_TEST_RUN_ID}web2"
osCliCapture account create -u "${otherAccountUsername}" -p "${accountPassword}" --is-super-admin false
assertCliStatus created "create the other browser account"
osCliCapture account get -n "${otherAccountUsername}"
assertCliStatus success "read the other browser account"
OS_TEST_OTHER_ACCOUNT_ID="$(jq -r '.body.accounts[0].id' <<<"${cliOutput}")"
export OS_TEST_OTHER_ACCOUNT_ID

export OS_TEST_ACCOUNT_USERNAME="${accountUsername}"
export OS_TEST_ACCOUNT_PASSWORD="${accountPassword}"
export OS_TEST_BROWSER_ARTIFACTS_DIR="${browserArtifactsDir}"
mkdir -p "${browserArtifactsDir}"

startedAtMs="$(( $(date +%s) * 1000 ))"

pushd "${browserDir}" >/dev/null || exit 1
if [[ ! -d node_modules ]]; then
	npm install --no-audit --no-fund
fi
npx playwright install chromium
npx playwright test
suiteExitCode=$?
popd >/dev/null || true

elapsedMs="$(( $(($(date +%s) * 1000)) - startedAtMs ))"
recordTiming "ui.browser/standard" "${elapsedMs}"

if ((suiteExitCode == 0)); then
	recordAssertionPass "browser suite"
else
	recordAssertionFail "browser suite (exit code ${suiteExitCode})"
fi

finishTests
