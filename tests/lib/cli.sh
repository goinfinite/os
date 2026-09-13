#!/usr/bin/env bash

if ! declare -F recordTiming >/dev/null; then
	source "$(dirname "${BASH_SOURCE[0]}")/timing.sh"
fi

cliOutput=""
cliErrorOutput=""
cliExitCode=0
cliStatus=""
cliBody=""
cliIsJson=false

readCurrentTimeMillis() {
	if [[ -n "${EPOCHREALTIME:-}" ]]; then
		local epochRealtime="${EPOCHREALTIME/./}"
		printf '%s' "$((10#${epochRealtime:0:13}))"
		return 0
	fi

	printf '%s' "$(($(date +%s) * 1000))"
}

resolveCliOperationName() {
	local operationName="" token tokenCount=0

	for token in "$@"; do
		if [[ "${token}" == -* ]]; then
			break
		fi

		operationName="${operationName:+${operationName}.}${token}"
		tokenCount=$((tokenCount + 1))
		if ((tokenCount == 3)); then
			break
		fi
	done

	printf '%s' "${operationName}"
}

osCliCapture() {
	local outputFile errorFile startedAtMs
	outputFile="$(mktemp)"
	errorFile="$(mktemp)"
	startedAtMs="$(readCurrentTimeMillis)"

	"${OS_TEST_RUNTIME}" exec "${OS_TEST_CONTAINER}" os "$@" >"${outputFile}" 2>"${errorFile}"
	cliExitCode=$?

	cliOutput="$(cat "${outputFile}")"
	cliErrorOutput="$(cat "${errorFile}")"
	rm -f "${outputFile}" "${errorFile}"

	cliIsJson=false
	cliStatus=""
	cliBody=""

	if jq -e . >/dev/null 2>&1 <<<"${cliOutput}"; then
		cliIsJson=true
		cliStatus="$(jq -r '.status' <<<"${cliOutput}")"
		cliBody="$(jq -c '.body' <<<"${cliOutput}")"
	fi

	recordTiming "$(resolveCliOperationName "$@")" "$(( $(readCurrentTimeMillis) - startedAtMs ))"
	return 0
}

osCliRun() {
	"${OS_TEST_RUNTIME}" exec "${OS_TEST_CONTAINER}" os "$@"
}

osBash() {
	"${OS_TEST_RUNTIME}" exec "${OS_TEST_CONTAINER}" bash -c "$1"
}

assertCliStatus() {
	local expected="$1" description="$2"

	if [[ "${expected}" == "${cliStatus}" ]]; then
		recordAssertionPass "${description}"
		return 0
	fi

	local responseSummary="${cliOutput:0:300}"
	if [[ -z "${responseSummary}" ]]; then
		responseSummary="${cliErrorOutput:0:300}"
	fi

	recordAssertionFail "${description} (expected '${expected}', got '${cliStatus}' response '${responseSummary}')"
	return 1
}

assertCliExitCode() {
	assertEquals "$1" "${cliExitCode}" "$2"
}

assertCliJq() {
	local jqFilter="$1" description="$2"

	if [[ "${cliIsJson}" != "true" ]]; then
		recordAssertionFail "${description} (response is not JSON: '${cliOutput}')"
		return 1
	fi

	if jq -e "${jqFilter}" >/dev/null 2>&1 <<<"${cliOutput}"; then
		recordAssertionPass "${description}"
		return 0
	fi

	recordAssertionFail "${description} (filter '${jqFilter}' failed on '${cliOutput}')"
	return 1
}

waitForCliJq() {
	local timeoutSeconds="$1" description="$2" jqFilter="$3"
	shift 3
	local deadline=$((SECONDS + timeoutSeconds))

	while ((SECONDS < deadline)); do
		osCliCapture "$@"

		if [[ "${cliStatus}" == "success" ]] && jq -e "${jqFilter}" >/dev/null 2>&1 <<<"${cliOutput}"; then
			recordAssertionPass "${description}"
			return 0
		fi

		sleep 1
	done

	recordAssertionFail "${description} (timeout after ${timeoutSeconds}s)"
	return 1
}
