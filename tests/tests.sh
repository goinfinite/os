#!/usr/bin/env bash
#
## SuiteConventions
#
# This suite intentionally deviates from .agents/rules/code/shell.md.
# The runner and the test scripts use `set -uo pipefail` without `-e`.
# Assertions must continue after a failure so the summary is complete, and finishTests computes the exit code.
# The scripts omit the @tag headers because tests/registry.yaml is the single source of test metadata.
#
set -uo pipefail

testsDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repoRoot="$(cd "${testsDir}/.." && pwd)"
registryPath="${testsDir}/registry.yaml"
goldenPath="${testsDir}/golden.yaml"
source "${testsDir}/lib/container.sh"
source "${testsDir}/lib/timing.sh"

validLevels=(fast standard exhaustive)
requestedLevel="standard"
requestedFeatures=()
requestedScopes=()
jobs=1
listOnly=false
cleanOnly=false
OS_TEST_REBUILD=false
runId="$(date +%s)$$"
artifactsDir="${testsDir}/artifacts/${runId}"
resultsFile="${artifactsDir}/results.tsv"
timingsFile="${artifactsDir}/timings.tsv"
OS_TEST_RUNTIME=""

printUsage() {
	cat <<'EOF'
Usage:
  tests/tests.sh [--level=<fast|standard|exhaustive>] [--feature=<feature>[,<feature>...]]
                 [--scope=<scope>[,<scope>...]] [--jobs=<n>] [--rebuild]
  tests/tests.sh list [--feature=<feature>[,<feature>...]] [--scope=<scope>[,<scope>...]]
  tests/tests.sh clean

Flags:
  --level    How deeply to test. Values: fast, standard, exhaustive. Default: standard.
  --feature  Which product feature to test. Comma-separated. Default: all registered features.
  --scope    Which testing dimension to run. Comma-separated. Default: every scope included by the level.
  --jobs     How many features to run in parallel, each with its own container. Default: 1.
  --rebuild  Rebuild the runtime and unit test images before running.

Subcommands:
  list       Print the registered test tree after applying any filters. Does not execute tests.
  clean      Stop and remove every container named os-test-*.
EOF
}

appendCommaSeparatedItems() {
	local rawList="$1"
	local -n targetList="$2"
	local item

	IFS=',' read -ra items <<<"${rawList}"
	for item in "${items[@]}"; do
		[[ -n "${item}" ]] || continue
		targetList+=("${item}")
	done
}

parseArguments() {
	local args=("$@")
	local index=0

	while ((index < ${#args[@]})); do
		local arg="${args[index]}"

		case "${arg}" in
		list) listOnly=true ;;
		clean) cleanOnly=true ;;
		--level=*) requestedLevel="${arg#*=}" ;;
		--feature=*) appendCommaSeparatedItems "${arg#*=}" requestedFeatures ;;
		--scope=*) appendCommaSeparatedItems "${arg#*=}" requestedScopes ;;
		--jobs=*) jobs="${arg#*=}" ;;
		--rebuild) OS_TEST_REBUILD=true ;;
		--level | --feature | --scope | --jobs)
			index=$((index + 1))
			if ((index >= ${#args[@]})); then
				echo "Missing value for ${arg}" >&2
				exit 2
			fi

			case "${arg}" in
			--level) requestedLevel="${args[index]}" ;;
			--feature) appendCommaSeparatedItems "${args[index]}" requestedFeatures ;;
			--scope) appendCommaSeparatedItems "${args[index]}" requestedScopes ;;
			--jobs) jobs="${args[index]}" ;;
			esac
			;;
		--help | -h)
			printUsage
			exit 0
			;;
		*)
			echo "Unknown argument '${arg}'" >&2
			printUsage >&2
			exit 2
			;;
		esac

		index=$((index + 1))
	done
}

levelRank() {
	case "$1" in
	fast) printf '1' ;;
	standard) printf '2' ;;
	exhaustive) printf '3' ;;
	*) printf '0' ;;
	esac
}

requireHostTools() {
	local commandName
	for commandName in yq jq timeout curl; do
		if command -v "${commandName}" >/dev/null 2>&1; then
			continue
		fi

		echo "MissingRequiredTool: ${commandName} (run 'mise install')" >&2
		exit 2
	done
}

registryFeatures() {
	yq -r '.features | keys | .[]' "${registryPath}"
}

registryEntriesJson() {
	yq -o=json -I=0 ".features[\"$1\"]" "${registryPath}"
}

registryFeatureExists() {
	registryFeatures | grep -x "$1" >/dev/null
}

registryScopeExists() {
	local scope="$1" feature

	while IFS= read -r feature; do
		if registryEntriesJson "${feature}" | jq -e --arg scope "${scope}" '.[] | select(.scope == $scope)' >/dev/null; then
			return 0
		fi
	done < <(registryFeatures)

	return 1
}

validateRegistry() {
	local feature entryJson

	while IFS= read -r feature; do
		while IFS= read -r entryJson; do
			[[ -n "${entryJson}" ]] || continue

			local scope level runCommand
			scope="$(jq -r '.scope // ""' <<<"${entryJson}")"
			level="$(jq -r '.level // ""' <<<"${entryJson}")"
			runCommand="$(jq -r '.run // ""' <<<"${entryJson}")"

			if [[ -z "${scope}" || -z "${runCommand}" || "$(levelRank "${level}")" == "0" ]]; then
				echo "MalformedRegistryEntry: feature ${feature}: ${entryJson}" >&2
				exit 2
			fi
		done < <(registryEntriesJson "${feature}" | jq -c '.[]')
	done < <(registryFeatures)
}

validateSelection() {
	if [[ "$(levelRank "${requestedLevel}")" == "0" ]]; then
		echo "UnknownLevel: ${requestedLevel} (valid: ${validLevels[*]})" >&2
		exit 2
	fi

	if ! [[ "${jobs}" =~ ^[1-9][0-9]*$ ]]; then
		echo "InvalidJobs: ${jobs} (expected a positive integer)" >&2
		exit 2
	fi

	local feature
	for feature in "${requestedFeatures[@]}"; do
		if ! registryFeatureExists "${feature}"; then
			echo "UnknownFeature: ${feature} (run 'tests/tests.sh list' to see registered features)" >&2
			exit 2
		fi
	done

	local scope
	for scope in "${requestedScopes[@]}"; do
		if ! registryScopeExists "${scope}"; then
			echo "UnknownScope: ${scope} (run 'tests/tests.sh list' to see registered scopes)" >&2
			exit 2
		fi
	done
}

featureIsRequested() {
	if ((${#requestedFeatures[@]} == 0)); then
		return 0
	fi

	local feature
	for feature in "${requestedFeatures[@]}"; do
		if [[ "${feature}" == "$1" ]]; then
			return 0
		fi
	done

	return 1
}

scopeIsRequested() {
	if ((${#requestedScopes[@]} == 0)); then
		return 0
	fi

	local scope
	for scope in "${requestedScopes[@]}"; do
		if [[ "${scope}" == "$1" ]]; then
			return 0
		fi
	done

	return 1
}

readSelectedEntriesJson() {
	local feature="$1" level="$2"
	local requestedRank
	requestedRank="$(levelRank "${level}")"

	local entryJson
	while IFS= read -r entryJson; do
		[[ -n "${entryJson}" ]] || continue

		local entryLevel entryScope
		entryLevel="$(jq -r '.level' <<<"${entryJson}")"
		entryScope="$(jq -r '.scope' <<<"${entryJson}")"

		if (($(levelRank "${entryLevel}") <= requestedRank)) && scopeIsRequested "${entryScope}"; then
			printf '%s\n' "${entryJson}"
		fi
	done < <(registryEntriesJson "${feature}" | jq -c '.[]')
}

printList() {
	local lines=()
	local provisionalLegend=false
	local feature

	while IFS= read -r feature; do
		featureIsRequested "${feature}" || continue

		local entryJson
		while IFS= read -r entryJson; do
			[[ -n "${entryJson}" ]] || continue

			local scope level runCommand provisional marker=""
			scope="$(jq -r '.scope' <<<"${entryJson}")"
			level="$(jq -r '.level' <<<"${entryJson}")"
			runCommand="$(jq -r '.run' <<<"${entryJson}")"
			scopeIsRequested "${scope}" || continue

			provisional="$(jq -r '.provisional // false' <<<"${entryJson}")"
			if [[ "${provisional}" == "true" ]]; then
				marker="*"
				provisionalLegend=true
			fi

			lines+=("$(printf '%-32s %s' "${feature}:${scope}/${level}${marker}" "${runCommand}")")
		done < <(registryEntriesJson "${feature}" | jq -c '.[]')
	done < <(registryFeatures)

	if ((${#lines[@]} == 0)); then
		printf 'No registered tests match the filters.\n' >&2
		return 0
	fi

	printf '%s\n' "${lines[@]}" | sort

	if [[ "${provisionalLegend}" == "true" ]]; then
		printf '(* provisional: registered but implementation-coupled, awaiting refactoring)\n' >&2
	fi
}

selectionHasScope() {
	local wantedScope="$1" feature entryJson

	while IFS= read -r feature; do
		featureIsRequested "${feature}" || continue

		while IFS= read -r entryJson; do
			[[ -n "${entryJson}" ]] || continue

			local scope level
			scope="$(jq -r '.scope' <<<"${entryJson}")"
			level="$(jq -r '.level' <<<"${entryJson}")"

			if [[ "${scope}" == "${wantedScope}" ]] && (($(levelRank "${level}") <= $(levelRank "${requestedLevel}"))) && scopeIsRequested "${scope}"; then
				return 0
			fi
		done < <(registryEntriesJson "${feature}" | jq -c '.[]')
	done < <(registryFeatures)

	return 1
}

selectionHasIntegrationEntry() {
	local feature entryJson

	while IFS= read -r feature; do
		featureIsRequested "${feature}" || continue

		while IFS= read -r entryJson; do
			[[ -n "${entryJson}" ]] || continue

			local scope level
			scope="$(jq -r '.scope' <<<"${entryJson}")"
			level="$(jq -r '.level' <<<"${entryJson}")"

			if [[ "${scope}" != "unit" ]] && (($(levelRank "${level}") <= $(levelRank "${requestedLevel}"))) && scopeIsRequested "${scope}"; then
				return 0
			fi
		done < <(registryEntriesJson "${feature}" | jq -c '.[]')
	done < <(registryFeatures)

	return 1
}

runEntry() {
	local feature="$1" scope="$2" level="$3" runCommand="$4"
	local entryKey="${feature}:${scope}/${level}"
	local startedAt=$SECONDS
	local exitCode=0

	printf '\n[%s] %s\n' "${entryKey}" "${runCommand}"

	if [[ "${scope}" == "unit" ]]; then
		timeout --signal=TERM --kill-after=30 "${osTestCommandTimeoutSeconds}" \
			bash -c 'source "${OS_TEST_TESTS_DIR}/lib/container.sh"; osTestRunUnitEntry "$1"' \
			_ "${runCommand}"
		exitCode=$?
	fi

	if [[ "${scope}" != "unit" ]]; then
		timeout --signal=TERM --kill-after=30 "${osTestCommandTimeoutSeconds}" \
			bash "${repoRoot}/${runCommand}"
		exitCode=$?
	fi

	local elapsed=$((SECONDS - startedAt))

	if ((exitCode == 124 || exitCode == 137)); then
		printf 'TIMEOUT %s exceeded the %ss ceiling\n' "${entryKey}" "${osTestCommandTimeoutSeconds}"
		printf '%s\t%s\t%s\n' "${entryKey}" "error" "${elapsed}" >>"${resultsFile}"
		return 0
	fi

	if ((exitCode != 0)); then
		printf 'FAIL  %s (%ss)\n' "${entryKey}" "${elapsed}"
		printf '%s\t%s\t%s\n' "${entryKey}" "fail" "${elapsed}" >>"${resultsFile}"
		return 0
	fi

	printf 'PASS  %s (%ss)\n' "${entryKey}" "${elapsed}"
	printf '%s\t%s\t%s\n' "${entryKey}" "pass" "${elapsed}" >>"${resultsFile}"
}

runFeatureGroup() {
	local feature="$1" level="$2"
	local entriesJson
	entriesJson="$(readSelectedEntriesJson "${feature}" "${level}")"

	if [[ -z "${entriesJson}" ]]; then
		return 0
	fi

	local hasIntegrationEntry=false
	local entryJson
	while IFS= read -r entryJson; do
		if [[ "$(jq -r '.scope' <<<"${entryJson}")" != "unit" ]]; then
			hasIntegrationEntry=true
			break
		fi
	done <<<"${entriesJson}"

	local containerName=""
	if [[ "${hasIntegrationEntry}" == "true" ]]; then
		containerName="$(osTestStartContainer "${feature}" "${level}")"
		if [[ -z "${containerName}" ]]; then
			printf 'FAIL  %s: container did not start\n' "${feature}"
			printf '%s\t%s\t%s\n' "${feature}:container/start" "error" "0" >>"${resultsFile}"
			return 0
		fi
	fi

	export OS_TEST_FEATURE="${feature}"
	export OS_TEST_RUNTIME
	export OS_TEST_RUN_ID="${runId}"
	export OS_TEST_REPO_ROOT="${repoRoot}"
	export OS_TEST_ARTIFACTS_DIR="${artifactsDir}"
	export OS_TEST_TIMINGS_FILE="${timingsFile}"
	export OS_TEST_CONTAINER="${containerName}"

	if [[ -n "${containerName}" ]]; then
		local apiHostPort httpHostPort httpsHostPort
		apiHostPort="$(resolvePublishedPort "${containerName}" "${osTestApiContainerPort}")"
		httpHostPort="$(resolvePublishedPort "${containerName}" "${osTestHttpContainerPort}")"
		httpsHostPort="$(resolvePublishedPort "${containerName}" "${osTestHttpsContainerPort}")"
		export OS_TEST_API_URL="https://127.0.0.1:${apiHostPort}"
		export OS_TEST_HTTP_PORT="${httpHostPort}"
		export OS_TEST_HTTPS_PORT="${httpsHostPort}"
	fi

	while IFS= read -r entryJson; do
		[[ -n "${entryJson}" ]] || continue

		local scope entryLevel runCommand
		scope="$(jq -r '.scope' <<<"${entryJson}")"
		entryLevel="$(jq -r '.level' <<<"${entryJson}")"
		runCommand="$(jq -r '.run' <<<"${entryJson}")"

		export OS_TEST_SCOPE="${scope}"
		runEntry "${feature}" "${scope}" "${entryLevel}" "${runCommand}"
	done <<<"${entriesJson}"

	if [[ -n "${containerName}" && "${level}" != "fast" ]]; then
		exportContainerLogs "${containerName}" "${artifactsDir}/${feature}.log"
		osTestStopContainer "${containerName}"
	fi
}

readSelectedFeatureNames() {
	local feature
	while IFS= read -r feature; do
		featureIsRequested "${feature}" || continue

		if [[ -n "$(readSelectedEntriesJson "${feature}" "${requestedLevel}")" ]]; then
			printf '%s\n' "${feature}"
		fi
	done < <(registryFeatures)
}

runAllFeatures() {
	local features=("$@")
	local feature

	if ((jobs <= 1)); then
		for feature in "${features[@]}"; do
			runFeatureGroup "${feature}" "${requestedLevel}"
		done
		return 0
	fi

	for feature in "${features[@]}"; do
		runFeatureGroup "${feature}" "${requestedLevel}" >"${artifactsDir}/${feature}.out" 2>&1 &
		while (($(jobs -rp | wc -l) >= jobs)); do
			wait -n
		done
	done
	wait

	for feature in "${features[@]}"; do
		[[ -f "${artifactsDir}/${feature}.out" ]] || continue
		printf '\n===== %s =====\n' "${feature}"
		cat "${artifactsDir}/${feature}.out"
	done
}

printSummary() {
	printf '\n========================================\n'
	printf 'Summary\n'

	local passCount=0 failCount=0 errorCount=0
	local entryKey status elapsed
	while IFS=$'\t' read -r entryKey status elapsed; do
		[[ -n "${entryKey}" ]] || continue

		case "${status}" in
		pass)
			passCount=$((passCount + 1))
			printf '  %-36s PASS (%ss)\n' "${entryKey}" "${elapsed}"
			;;
		fail)
			failCount=$((failCount + 1))
			printf '  %-36s FAIL (%ss)\n' "${entryKey}" "${elapsed}"
			;;
		error)
			errorCount=$((errorCount + 1))
			printf '  %-36s ERROR (%ss)\n' "${entryKey}" "${elapsed}"
			;;
		esac
	done < <(sort -u "${resultsFile}")

	printf '  passed: %s, failed: %s, errors: %s\n' "${passCount}" "${failCount}" "${errorCount}"

	local fatalTiming=false
	if ! reportTimingBudgets "${goldenPath}" "${timingsFile}"; then
		fatalTiming=true
	fi

	if ((errorCount > 0)); then
		printf 'Overall: ERROR (exit 2)\n'
		return 2
	fi

	if ((failCount > 0)) || [[ "${fatalTiming}" == "true" ]]; then
		printf 'Overall: FAIL (exit 1)\n'
		return 1
	fi

	printf 'Overall: PASS (exit 0)\n'
	return 0
}

handleInterrupt() {
	trap - INT TERM

	local childPid
	for childPid in $(jobs -p); do
		kill "${childPid}" 2>/dev/null || true
	done

	osTestCleanupContainers
	exit 130
}

main() {
	parseArguments "$@"
	requireHostTools

	if [[ ! -f "${registryPath}" ]]; then
		echo "RegistryNotFound: ${registryPath}" >&2
		exit 2
	fi

	validateRegistry
	validateSelection

	if [[ "${listOnly}" == "true" ]]; then
		printList
		exit 0
	fi

	OS_TEST_RUNTIME="$(osTestRuntimeResolver)" || exit 2
	export OS_TEST_RUNTIME

	if [[ "${cleanOnly}" == "true" ]]; then
		osTestRemoveAllTestContainers
		exit 0
	fi

	local featureNames
	featureNames="$(readSelectedFeatureNames)"
	if [[ -z "${featureNames}" ]]; then
		echo "NoTestsSelected: no registered entry matches level '${requestedLevel}' with the requested filters" >&2
		exit 2
	fi

	mkdir -p "${artifactsDir}"
	touch "${resultsFile}" "${timingsFile}"

	export OS_TEST_REPO_ROOT="${repoRoot}"
	export OS_TEST_ARTIFACTS_DIR="${artifactsDir}"
	export OS_TEST_RUN_ID="${runId}"
	export OS_TEST_TIMINGS_FILE="${timingsFile}"
	export OS_TEST_TESTS_DIR="${testsDir}"

	if selectionHasScope "unit"; then
		if [[ ! -f "${repoRoot}/.env" ]]; then
			echo "MissingEnvFile: copy .env.example to .env before running unit tests" >&2
			exit 2
		fi
		osTestEnsureUnitImage || exit 2
	fi

	if selectionHasIntegrationEntry; then
		osTestBuildBinary || exit 2
		osTestEnsureRuntimeImage || exit 2
	fi

	local features=()
	while IFS= read -r feature; do
		features+=("${feature}")
	done <<<"${featureNames}"

	trap handleInterrupt INT TERM
	runAllFeatures "${features[@]}"
	osTestCleanupContainers
	trap - INT TERM

	printSummary
	exit $?
}

main "$@"
