#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

echo "Feature ${OS_TEST_FEATURE}: system overview"

osCliCapture o11y overview
assertCliStatus success "read system overview"
assertCliExitCode 0 "system overview exit code"
assertCliJq '.body.hostname | length > 0' "overview reports a hostname"
assertCliJq '.body.uptimeSecs > 0' "overview reports uptime"
assertCliJq '.body.specs.cpuCores > 0' "overview reports cpu cores"
assertCliJq '.body.specs.memoryTotal | length > 0' "overview reports total memory"
assertCliJq '.body.currentUsage.cpuUsagePercent >= 0' "overview reports cpu usage"
assertCliJq '.body.currentUsage.storageUsage >= 0' "overview reports storage usage"

finishTests
