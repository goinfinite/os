#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

serviceName="mariadb"
dbName="os_test_${OS_TEST_RUN_ID}"
dbUser="os_test_${OS_TEST_RUN_ID}"
dbPassword='abc123!abc'

echo "Feature ${OS_TEST_FEATURE}: database lifecycle on mariadb"

osCliCapture db create -t sqlite -n "${dbName}"
assertCliJq '.status != "success"' "creating a database with an unsupported type fails"

osCliCapture services create-installable -n mariadb -v 10.11
assertCliStatus created "install mariadb service"
assertCliExitCode 0 "install mariadb exit code"

waitForCliJq 900 "mariadb service reaches the running state" \
	'.body.installedServices[] | select(.name == "mariadb" and .status == "running")' \
	services get -n "${serviceName}" -m 50

osCliCapture db get -t mariadb -n "${dbName}"
assertCliStatus success "read databases before creation"
assertCliJq ".body.databases | map(select(.name == \"${dbName}\")) | length == 0" "database is absent before creation"

osCliCapture db create -t mariadb -n "${dbName}"
assertCliStatus created "create database"
assertCliExitCode 0 "create database exit code"

osCliCapture db get -t mariadb -n "${dbName}"
assertCliStatus success "read created database"
assertCliJq ".body.databases | map(select(.name == \"${dbName}\" and .type == \"mariadb\")) | length == 1" "created database is persisted"

osCliCapture db create-user -t mariadb -n "${dbName}" -u "${dbUser}" -p "${dbPassword}" -r 'ALL'
assertCliStatus created "create database user"
assertCliExitCode 0 "create database user exit code"

osCliCapture db get -t mariadb -n "${dbName}"
assertCliJq ".body.databases[0].users | map(select(.username == \"${dbUser}\")) | length == 1" "created database user is persisted"

assertCommandSucceeds "database user can connect over TCP and query" \
	osBash "mariadb --protocol=tcp -h 127.0.0.1 -u '${dbUser}' -p'${dbPassword}' '${dbName}' -e 'SELECT 1' >/dev/null 2>&1"

osCliCapture db delete-user -t mariadb -n "${dbName}" -u "${dbUser}"
assertCliStatus success "delete database user"
assertCliExitCode 0 "delete database user exit code"

osCliCapture db get -t mariadb -n "${dbName}"
assertCliJq '.body.databases[0].users | length == 0' "deleted database user is absent"

assertCommandFails "deleted database user can no longer connect" \
	osBash "mariadb --protocol=tcp -h 127.0.0.1 -u '${dbUser}' -p'${dbPassword}' '${dbName}' -e 'SELECT 1' >/dev/null 2>&1"

osCliCapture db delete -t mariadb -n "${dbName}"
assertCliStatus success "delete database"
assertCliExitCode 0 "delete database exit code"

osCliCapture db get -t mariadb -n "${dbName}"
assertCliJq ".body.databases | map(select(.name == \"${dbName}\")) | length == 0" "deleted database is absent"

osCliCapture services delete -n "${serviceName}"
assertCliStatus success "delete mariadb service"
assertCliExitCode 0 "delete mariadb exit code"

waitForCliJq 300 "deleted mariadb service is absent" \
	'.body.installedServices | map(select(.name == "mariadb")) | length == 0' \
	services get -n "${serviceName}" -m 50

finishTests
