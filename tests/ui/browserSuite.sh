#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

browserDir="${OS_TEST_REPO_ROOT}/tests/ui/browser"
accountUsername="t${OS_TEST_RUN_ID}web"
accountPassword='abc123!abc'
browserLevel="${OS_TEST_LEVEL:-standard}"
browserArtifactsDir="${OS_TEST_ARTIFACTS_DIR}/browser/${browserLevel}"
browserToolingDir="${OS_TEST_ARTIFACTS_DIR}/browser-tooling"

echo "Feature ${OS_TEST_FEATURE}: browser suite (${browserLevel})"

otherAccountUsername="t${OS_TEST_RUN_ID}web2"

ensureBrowserAccount() {
	local username="$1"

	osCliCapture account get -n "${username}"
	if [[ "${cliStatus}" == "success" ]] && jq -e '.body.accounts | length > 0' >/dev/null 2>&1 <<<"${cliOutput}"; then
		recordAssertionPass "browser account ${username} exists"
		return 0
	fi

	osCliCapture account create -u "${username}" -p "${accountPassword}" --is-super-admin false
	assertCliStatus created "create browser account ${username}"
}

ensureBrowserAccount "${accountUsername}"
ensureBrowserAccount "${otherAccountUsername}"
osCliCapture account get -n "${otherAccountUsername}"
assertCliStatus success "read the other browser account"
OS_TEST_OTHER_ACCOUNT_ID="$(jq -r '.body.accounts[0].id' <<<"${cliOutput}")"
export OS_TEST_OTHER_ACCOUNT_ID

export OS_TEST_ACCOUNT_USERNAME="${accountUsername}"
export OS_TEST_ACCOUNT_PASSWORD="${accountPassword}"
export OS_TEST_BROWSER_ARTIFACTS_DIR="${browserArtifactsDir}"
mkdir -p "${browserArtifactsDir}" "${browserToolingDir}"

playwrightProjects=(--project=setup --project=lightpanda)

if [[ "${browserLevel}" == "exhaustive" ]]; then
	playwrightProjects=(--project=setup --project=chromium --project=firefox)
else
	loginResponse="$(curl -sk -X POST "${OS_TEST_API_URL}/api/v1/auth/login/" \
		-H 'Content-Type: application/json' \
		-d "{\"username\":\"${accountUsername}\",\"password\":\"${accountPassword}\"}")"
	accessToken="$(jq -r '.body.tokenStr' <<<"${loginResponse}")"

	if [[ -z "${accessToken}" || "${accessToken}" == "null" ]]; then
		echo "AccessTokenUnminted: ${loginResponse}" >&2
		exit 2
	fi

	cookieFile="${browserToolingDir}/lightpanda-cookies.txt"
	printf '127.0.0.1\tFALSE\t/\tTRUE\t0\tos-access-token\t%s\n' "${accessToken}" >"${cookieFile}"

	# Lightpanda's WebSocket client only trusts the app's self-signed certificate when the CA bundle and the app certificate are both loaded.
	caBundleFile="${browserToolingDir}/ca-bundle.crt"
	osBash "cat /etc/ssl/certs/ca-certificates.crt" >"${caBundleFile}"
	osCertFile="${browserToolingDir}/os.crt"
	osBash "cat /infinite/pki/os.crt" >"${osCertFile}"

	lightpandaPort=$((23000 + 10#${OS_TEST_RUN_ID: -4} % 1000))
	export OS_TEST_LIGHTPANDA_CDP_URL="ws://127.0.0.1:${lightpandaPort}"
	lightpandaLog="${browserToolingDir}/lightpanda.log"

	lightpanda serve --host 127.0.0.1 --port "${lightpandaPort}" \
		--cookie "${cookieFile}" \
		--ca-cert "${caBundleFile}" --ca-cert "${osCertFile}" \
		--insecure-disable-tls-host-verification >"${lightpandaLog}" 2>&1 &
	lightpandaPid=$!
	trap 'kill "${lightpandaPid}" 2>/dev/null || true' EXIT

	lightpandaReady=false
	for _ in $(seq 1 50); do
		if curl -sf "http://127.0.0.1:${lightpandaPort}/json/version" >/dev/null 2>&1; then
			lightpandaReady=true
			break
		fi
		sleep 0.2
	done

	if [[ "${lightpandaReady}" != "true" ]]; then
		echo "LightpandaServeFailed: see ${lightpandaLog}" >&2
		exit 2
	fi
fi

startedAtMs="$(( $(date +%s) * 1000 ))"

pushd "${browserDir}" >/dev/null || exit 1
if [[ ! -d node_modules ]]; then
	npm install --no-audit --no-fund --ignore-scripts
fi
if [[ "${browserLevel}" == "exhaustive" ]]; then
	if ! ./node_modules/.bin/playwright install chromium firefox; then
		echo "PlaywrightBrowserInstallFailed" >&2
		exit 2
	fi
fi
./node_modules/.bin/playwright test "${playwrightProjects[@]}"
suiteExitCode=$?
popd >/dev/null || true

elapsedMs="$(( $(($(date +%s) * 1000)) - startedAtMs ))"
recordTiming "ui.browser/${browserLevel}" "${elapsedMs}"

if ((suiteExitCode == 0)); then
	recordAssertionPass "browser suite"
else
	recordAssertionFail "browser suite (exit code ${suiteExitCode})"
fi

finishTests
