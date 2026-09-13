#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

echo "Feature ${OS_TEST_FEATURE}: marketplace catalog"

osCliCapture mktplace list-catalog -m 5
assertCliStatus success "read marketplace catalog"
assertCliExitCode 0 "read marketplace catalog exit code"
assertCliJq '.body.marketplaceCatalogItems | length > 0' "catalog returns items"
assertCliJq '.body.marketplaceCatalogItems[0].id > 0' "catalog items expose an id"
assertCliJq '.body.marketplaceCatalogItems[0].name | length > 0' "catalog items expose a name"
assertCliJq '.body.marketplaceCatalogItems[0].slugs | length > 0' "catalog items expose slugs"

osCliCapture mktplace list-catalog -s pocketbase
assertCliStatus success "filter marketplace catalog by slug"
assertCliJq '.body.marketplaceCatalogItems | length >= 1 and (map(.slugs | index("pocketbase") != null) | all)' "slug filter returns only items that carry the pocketbase slug"

osCliCapture mktplace list
assertCliStatus success "read installed marketplace items"
assertCliJq '.body.marketplaceInstalledItems | type == "array"' "installed items response is a list"

osCliCapture mktplace install
assertCliStatus infraError "install without id or slug fails"
assertCliExitCode 69 "install without id or slug exit code"

osCliCapture mktplace delete -i 65000
assertCliStatus infraError "delete nonexistent installed item fails"
assertCliExitCode 69 "delete nonexistent installed item exit code"

finishTests
