#!/usr/bin/env bash

recordTiming() {
	local operation="$1" elapsedMs="$2"

	if [[ -z "${OS_TEST_TIMINGS_FILE:-}" || -z "${operation}" ]]; then
		return 0
	fi

	printf '%s\t%s\n' "${operation}" "${elapsedMs}" >>"${OS_TEST_TIMINGS_FILE}"
}

calculateTimingPercentile() {
	local samples="$1" percent="$2"

	printf '%s\n' "${samples}" | awk -v p="${percent}" '
		{ values[NR] = $1 }
		END {
			if (NR == 0) { exit }
			idx = int((p / 100) * NR)
			if (idx < 1) { idx = 1 }
			if (idx > NR) { idx = NR }
			print values[idx]
		}
	'
}

reportOperationTiming() {
	local operation="$1" p50="$2" p90="$3" baselineP90="$4"
	local infoMs="$5" warnMs="$6" fatalMs="$7"
	local relativeWarnPercent="$8" relativeFatalPercent="$9"

	local summary="${operation} p50 ${p50}ms p90 ${p90}ms"
	local relativePercent=""
	if [[ -n "${baselineP90}" && "${baselineP90}" -gt 0 ]]; then
		relativePercent="$(awk -v p90="${p90}" -v base="${baselineP90}" 'BEGIN { printf "%.1f", ((p90 - base) / base) * 100 }')"
		summary="${summary} (baseline p90 ${baselineP90}ms, ${relativePercent}%)"
	fi

	if [[ -n "${fatalMs}" && "${p90}" -gt "${fatalMs}" ]]; then
		printf 'FATAL %s exceeds fatal %sms\n' "${summary}" "${fatalMs}"
		return 1
	fi

	if [[ -n "${relativeFatalPercent}" && -n "${relativePercent}" ]]; then
		if awk -v r="${relativePercent}" -v limit="${relativeFatalPercent}" 'BEGIN { exit !(r > limit) }'; then
			printf 'FATAL %s exceeds relative fatal %s%%\n' "${summary}" "${relativeFatalPercent}"
			return 1
		fi
	fi

	if [[ -n "${warnMs}" && "${p90}" -gt "${warnMs}" ]]; then
		printf 'WARN  %s exceeds warn %sms\n' "${summary}" "${warnMs}"
		return 0
	fi

	if [[ -n "${relativeWarnPercent}" && -n "${relativePercent}" ]]; then
		if awk -v r="${relativePercent}" -v limit="${relativeWarnPercent}" 'BEGIN { exit !(r > limit) }'; then
			printf 'WARN  %s exceeds relative warn %s%%\n' "${summary}" "${relativeWarnPercent}"
			return 0
		fi
	fi

	if [[ -n "${infoMs}" && "${p90}" -gt "${infoMs}" ]]; then
		printf 'INFO  %s exceeds info %sms\n' "${summary}" "${infoMs}"
		return 0
	fi

	printf 'OK    %s\n' "${summary}"
	return 0
}

reportTimingBudgets() {
	local goldenPath="$1" timingsPath="$2"

	if [[ ! -f "${goldenPath}" || ! -f "${timingsPath}" ]]; then
		return 0
	fi

	local fatalCount=0
	local operation
	while IFS= read -r operation; do
		[[ -n "${operation}" ]] || continue

		local samples p50 p90 baselineP90 warnMs fatalMs infoMs relativeWarnPercent relativeFatalPercent
		samples="$(awk -F'\t' -v op="${operation}" '$1 == op { print $2 }' "${timingsPath}" | sort -n)"
		p50="$(calculateTimingPercentile "${samples}" 50)"
		p90="$(calculateTimingPercentile "${samples}" 90)"
		baselineP90="$(yq -r ".[\"${operation}\"][\"baseline_p90_ms\"] // \"\"" "${goldenPath}")"
		infoMs="$(yq -r ".[\"${operation}\"][\"info_ms\"] // \"\"" "${goldenPath}")"
		warnMs="$(yq -r ".[\"${operation}\"][\"warn_ms\"] // \"\"" "${goldenPath}")"
		fatalMs="$(yq -r ".[\"${operation}\"][\"fatal_ms\"] // \"\"" "${goldenPath}")"
		relativeWarnPercent="$(yq -r ".[\"${operation}\"][\"relative_warn_percent\"] // \"\"" "${goldenPath}")"
		relativeFatalPercent="$(yq -r ".[\"${operation}\"][\"relative_fatal_percent\"] // \"\"" "${goldenPath}")"

		reportOperationTiming \
			"${operation}" "${p50}" "${p90}" "${baselineP90}" "${infoMs}" \
			"${warnMs}" "${fatalMs}" "${relativeWarnPercent}" "${relativeFatalPercent}"
		if [[ "$?" == "1" ]]; then
			fatalCount=$((fatalCount + 1))
		fi
	done < <(cut -f1 "${timingsPath}" | sort -u)

	if ((fatalCount > 0)); then
		return 1
	fi

	return 0
}
