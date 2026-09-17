#!/usr/bin/env bash
set -uo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/../lib/assert.sh"
source "$(dirname "${BASH_SOURCE[0]}")/../lib/cli.sh"

accountUsername="t${OS_TEST_RUN_ID}ui"
accountPassword='abc123!abc'
apiBaseUrl="https://127.0.0.1:1618"
proxyBaseHref="/containerId/"

echo "Feature ${OS_TEST_FEATURE}: base href and swagger documentation"

osCliCapture account create -u "${accountUsername}" -p "${accountPassword}" -s true
assertCliStatus created "account created for base href tests"

osCliCapture auth login -u "${accountUsername}" -p "${accountPassword}" -i 127.0.0.1
assertCliStatus success "login to read the dashboard"
sessionToken="$(jq -r '.body.tokenStr' <<<"${cliOutput}")"
assertNotEmpty "${sessionToken}" "login returns a session token"

loginPageHtml="$(osBash "curl -sk ${apiBaseUrl}/login/")"
assertContains "${loginPageHtml}" '<base href="/">' "login page uses the root base href"

proxiedLoginPageHtml="$(osBash "curl -sk -H 'X-Base-Href: ${proxyBaseHref}' ${apiBaseUrl}/login/")"
assertContains "${proxiedLoginPageHtml}" "<base href=\"${proxyBaseHref}\">" "login page uses the reverse proxy base href"

swaggerIndexRedirect="$(osBash "curl -sk -o /dev/null -w '%{redirect_url}' ${apiBaseUrl}/api/swagger/")"
assertEquals "${apiBaseUrl}/api/swagger/index.html" "${swaggerIndexRedirect}" "swagger root redirects to its index"

proxiedSwaggerIndexRedirect="$(osBash "curl -sk -o /dev/null -w '%{redirect_url}' -H 'X-Base-Href: ${proxyBaseHref}' ${apiBaseUrl}/api/swagger/")"
assertEquals "${apiBaseUrl}${proxyBaseHref}api/swagger/index.html" "${proxiedSwaggerIndexRedirect}" "swagger root redirect carries the reverse proxy base href"

swaggerIndexHtml="$(osBash "curl -sk ${apiBaseUrl}/api/swagger/index.html")"
assertContains "${swaggerIndexHtml}" "Swagger UI" "swagger index serves the documentation"

swaggerTrailingSlashRedirect="$(osBash "curl -sk -o /dev/null -w '%{redirect_url}' ${apiBaseUrl}/api/swagger")"
assertEquals "${apiBaseUrl}/api/swagger/" "${swaggerTrailingSlashRedirect}" "api path without trailing slash is normalized"

proxiedTrailingSlashRedirect="$(osBash "curl -sk -o /dev/null -w '%{redirect_url}' -H 'X-Base-Href: ${proxyBaseHref}' ${apiBaseUrl}/api/swagger")"
assertEquals "${apiBaseUrl}${proxyBaseHref}api/swagger/" "${proxiedTrailingSlashRedirect}" "trailing slash normalization carries the reverse proxy base href"

protocolRelativeBaseHrefRedirect="$(osBash "curl -sk -o /dev/null -w '%{redirect_url}' -H 'X-Base-Href: //evil.com/' ${apiBaseUrl}/api/swagger")"
assertEquals "${apiBaseUrl}/api/swagger/" "${protocolRelativeBaseHrefRedirect}" "protocol-relative base href is rejected"

accountRedirect="$(osBash "curl -sk -o /dev/null -w '%{redirect_url}' ${apiBaseUrl}/api/v1/account")"
assertEquals "${apiBaseUrl}/api/v1/account/" "${accountRedirect}" "api endpoint without trailing slash is normalized"

proxiedAccountRedirect="$(osBash "curl -sk -o /dev/null -w '%{redirect_url}' -H 'X-Base-Href: ${proxyBaseHref}' ${apiBaseUrl}/api/v1/account")"
assertEquals "${apiBaseUrl}${proxyBaseHref}api/v1/account/" "${proxiedAccountRedirect}" "api endpoint normalization carries the reverse proxy base href"

dashboardHtml="$(osBash "curl -sk -H 'Authorization: Bearer ${sessionToken}' ${apiBaseUrl}/overview/")"
assertContains "${dashboardHtml}" 'href="/api/swagger/"' "sidebar documentation link uses the root base href"

proxiedDashboardHtml="$(osBash "curl -sk -H 'Authorization: Bearer ${sessionToken}' -H 'X-Base-Href: ${proxyBaseHref}' ${apiBaseUrl}/overview/")"
assertContains "${proxiedDashboardHtml}" "href=\"${proxyBaseHref}api/swagger/\"" "sidebar documentation link carries the reverse proxy base href"

finishTests
