#!/usr/bin/env bash

assertionPassCount=0
assertionFailCount=0

recordAssertionPass() {
	printf 'PASS  %s\n' "$1"
	assertionPassCount=$((assertionPassCount + 1))
}

recordAssertionFail() {
	printf 'FAIL  %s\n' "$1"
	assertionFailCount=$((assertionFailCount + 1))
}

assertEquals() {
	local expected="$1" actual="$2" description="$3"

	if [[ "${expected}" == "${actual}" ]]; then
		recordAssertionPass "${description}"
		return 0
	fi

	recordAssertionFail "${description} (expected '${expected}', got '${actual}')"
	return 1
}

assertNotEmpty() {
	local value="$1" description="$2"

	if [[ -n "${value}" ]]; then
		recordAssertionPass "${description}"
		return 0
	fi

	recordAssertionFail "${description} (value is empty)"
	return 1
}

assertContains() {
	local content="$1" expectedSubstring="$2" description="$3"

	if [[ "${content}" == *"${expectedSubstring}"* ]]; then
		recordAssertionPass "${description}"
		return 0
	fi

	recordAssertionFail "${description} (missing '${expectedSubstring}' in '${content:0:200}')"
	return 1
}

assertNumberGreaterThan() {
	local threshold="$1" actual="$2" description="$3"

	if ((actual > threshold)); then
		recordAssertionPass "${description}"
		return 0
	fi

	recordAssertionFail "${description} (expected a value greater than ${threshold}, got ${actual})"
	return 1
}

assertCommandSucceeds() {
	local description="$1"
	shift

	if "$@"; then
		recordAssertionPass "${description}"
		return 0
	fi

	recordAssertionFail "${description}"
	return 1
}

assertCommandFails() {
	local description="$1"
	shift

	if "$@"; then
		recordAssertionFail "${description}"
		return 1
	fi

	recordAssertionPass "${description}"
	return 0
}

finishTests() {
	printf '\n%s passed, %s failed\n' "${assertionPassCount}" "${assertionFailCount}"

	if ((assertionFailCount > 0)); then
		exit 1
	fi

	exit 0
}
