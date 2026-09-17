#!/usr/bin/env bash

osTestRuntimeImage="os-test:latest"
osTestUnitImage="os-unit-test:latest"
osTestApiContainerPort=1618
osTestHttpContainerPort=80
osTestHttpsContainerPort=443
osTestContainerCpuLimit=2
osTestContainerMemoryLimit=2g
osTestApiWaitSeconds=180
osTestCommandTimeoutSeconds=1800

osTestRuntimeResolver() {
	if command -v podman >/dev/null 2>&1; then
		printf 'podman'
		return 0
	fi

	if command -v docker >/dev/null 2>&1; then
		printf 'docker'
		return 0
	fi

	echo "NoContainerRuntimeFound: install podman or docker" >&2
	return 2
}

osTestImageExists() {
	"${OS_TEST_RUNTIME}" image inspect "$1" >/dev/null 2>&1
}

resolveBinaryArch() {
	local containerArch=""
	if osTestImageExists "${osTestRuntimeImage}"; then
		containerArch="$("${OS_TEST_RUNTIME}" image inspect --format '{{.Architecture}}' "${osTestRuntimeImage}")"
	fi

	case "${containerArch}" in
	amd64 | arm64)
		printf '%s' "${containerArch}"
		return 0
		;;
	esac

	case "$(uname -m)" in
	aarch64 | arm64) printf 'arm64' ;;
	*) printf 'amd64' ;;
	esac
}

osTestBuildBinary() {
	if ! command -v templ >/dev/null 2>&1 || ! command -v go >/dev/null 2>&1; then
		echo "MissingRequiredTool: templ and go (run 'mise install')" >&2
		return 2
	fi

	local goArch
	goArch="$(resolveBinaryArch)"

	echo "Building the application binary for linux/${goArch}..."
	(
		cd "${OS_TEST_REPO_ROOT}" &&
			templ generate -path src/presentation/ui >/dev/null &&
			CGO_ENABLED=0 GOOS=linux GOARCH="${goArch}" go build -o ./bin/os ./os.go
	)
}

osTestEnsureRuntimeImage() {
	if osTestImageExists "${osTestRuntimeImage}" && [[ "${OS_TEST_REBUILD}" != "true" ]]; then
		return 0
	fi

	echo "Building the runtime image ${osTestRuntimeImage}..."
	timeout "${osTestCommandTimeoutSeconds}" \
		"${OS_TEST_RUNTIME}" build -t "${osTestRuntimeImage}" "${OS_TEST_REPO_ROOT}"
}

osTestEnsureUnitImage() {
	if osTestImageExists "${osTestUnitImage}" && [[ "${OS_TEST_REBUILD}" != "true" ]]; then
		return 0
	fi

	echo "Building the unit test image ${osTestUnitImage}..."
	timeout "${osTestCommandTimeoutSeconds}" \
		"${OS_TEST_RUNTIME}" build --target test -t "${osTestUnitImage}" "${OS_TEST_REPO_ROOT}"
}

exportContainerLogs() {
	local containerName="$1" outputFile="$2"
	"${OS_TEST_RUNTIME}" logs "${containerName}" >"${outputFile}" 2>&1 || true
}

osTestWaitForApi() {
	local containerName="$1"
	local deadline=$((SECONDS + osTestApiWaitSeconds))

	while ((SECONDS < deadline)); do
		if "${OS_TEST_RUNTIME}" exec "${containerName}" \
			curl -ks -o /dev/null --max-time 2 "https://127.0.0.1:${osTestApiContainerPort}/"; then
			return 0
		fi

		sleep 1
	done

	echo "Container ${containerName} did not serve the API within ${osTestApiWaitSeconds}s" >&2
	exportContainerLogs "${containerName}" "${OS_TEST_ARTIFACTS_DIR}/${containerName}.log"
	return 1
}

osTestStartContainer() {
	local feature="$1"
	local containerName="os-test-${feature}-${OS_TEST_RUN_ID}"

	if ! "${OS_TEST_RUNTIME}" run -d --name "${containerName}" \
		--hostname "${feature}.test" \
		--env "PRIMARY_VHOST=${feature}.test" \
		--cpus "${osTestContainerCpuLimit}" \
		--memory "${osTestContainerMemoryLimit}" \
		--volume "${OS_TEST_REPO_ROOT}/bin:/infinite/bin:ro" \
		--security-opt label=disable \
		--publish "127.0.0.1::${osTestApiContainerPort}" \
		--publish "127.0.0.1::${osTestHttpContainerPort}" \
		--publish "127.0.0.1::${osTestHttpsContainerPort}" \
		--entrypoint /bin/bash \
		"${osTestRuntimeImage}" \
		-c 'rm -f /infinite/os && ln -s /infinite/bin/os /infinite/os && exec /usr/bin/supervisord -c /infinite/supervisord.conf' \
		>/dev/null; then
		echo "ContainerStartFailed: ${containerName}" >&2
		return 1
	fi

	osTestWaitForApi "${containerName}"
	printf '%s' "${containerName}"
}

resolvePublishedPort() {
	local containerName="$1" containerPort="$2"
	"${OS_TEST_RUNTIME}" port "${containerName}" "${containerPort}/tcp" | sed 's/.*://'
}

osTestStopContainer() {
	local containerName="$1"
	"${OS_TEST_RUNTIME}" rm -f "${containerName}" >/dev/null 2>&1 || true
}

osTestRecordContainer() {
	printf '%s\n' "$1" >>"${OS_TEST_ARTIFACTS_DIR}/containers.txt"
}

osTestCleanupContainers() {
	local containersFile="${OS_TEST_ARTIFACTS_DIR}/containers.txt"

	if [[ ! -f "${containersFile}" ]]; then
		return 0
	fi

	local containerName
	while IFS= read -r containerName; do
		[[ -n "${containerName}" ]] || continue
		osTestStopContainer "${containerName}"
	done <"${containersFile}"
}

osTestRemoveAllTestContainers() {
	local containerNames
	containerNames="$("${OS_TEST_RUNTIME}" ps -a --format '{{.Names}}' | grep -E '^os-test-' || true)"

	if [[ -z "${containerNames}" ]]; then
		echo "No test containers found."
		return 0
	fi

	local containerName
	while IFS= read -r containerName; do
		[[ -n "${containerName}" ]] || continue
		echo "Removing ${containerName}..."
		osTestStopContainer "${containerName}"
	done <<<"${containerNames}"
}

osTestRunUnitEntry() {
	local runCommand="$1"
	local testArguments="${runCommand#go test }"

	if [[ "${testArguments}" == "${runCommand}" ]]; then
		echo "Unit registry entry must start with 'go test': ${runCommand}" >&2
		return 2
	fi

	local testPackages=()
	read -ra testPackages <<<"${testArguments}"

	"${OS_TEST_RUNTIME}" run --rm --entrypoint go \
		--volume "${OS_TEST_REPO_ROOT}:/infinite:rw" \
		--security-opt label=disable \
		--workdir /infinite \
		"${osTestUnitImage}" test -v "${testPackages[@]}"
}
