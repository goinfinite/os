# Changelog

```log
0.3.0 - 2026/07/03
fix(internalSetup): create /app/html at runtime to prevent fuse-overlayfs whiteout accumulation
refactor(infra): rename webServer package to internalSetup
fix(db): enforce UTC timestamps in gorm NowFunc for all database services

0.2.9 - 2026/04/09
feat(auth): migrate to tk TrustedCidrsReader and add TRUSTED_CIDRS env var
docs(readme): document verbose API panic responses for trusted networks
fix(ssl): clean up stale /validateOwnership mappings before creating new ones
feat(ssl): add SKIP_DNS_OWNERSHIP_CHECK env var to bypass DNS ownership check
refactor(ssl): move mapping repos and ownership path to SslCmdRepo struct fields
docs(readme): add environment variables reference table
chore: update go deps
fix(swagger): add missing dummy imports for swagger annotation compliance

0.2.8 - 2026/03/24
chore: remove temporary echo-swagger replace directive and upgrade to v1.5.2
fix(ssl): wrap SslPrivateKey around tk EnvelopedPrivateKey
fix(ssl): repair altnames filter and add pagination to ssl listing
feat(api): add swagger dto import for account endpoint
feat(ui): add database user and alias shortcut buttons
fix(files): handle root directory in file tree builder
fix(database): use WeakPassword for database user creation
fix(api): read operatorAccountId from echo context in file endpoints
refactor(liaison): migrate manual pagination to tk PaginationParser
refactor(cli): migrate to tk SimpleCliResponseRenderer
refactor(api): migrate to tk LiaisonApiResponseEmitter
refactor(liaison): migrate LiaisonOutput to tk LiaisonResponse
fix: get key path via replace
fix: show error on ssl list error
fix: add primaryKeyColumn to PaginationQueryBuilder calls
fix: sort services by name in ReadFirstInstalledItem
refactor(valueObject): eliminate pure alias files with direct tk usage
feat(files): add file privileges normalizer with ownership resolution
fix(services): use name sort for installable services pagination
refactor(infra): replace PaginationQueryBuilder alias with tk direct usage
test(auth): add security tests with repo constructors
feat(auth): add jwt v5 algorithm validation and api key hash security
refactor: remove obsolete migrateOperatorAccountIdToSri
refactor: use ActivityRecordLevelSecurity const and inline NewSriAccount
refactor: delete presentation wrappers and use tk directly
refactor: migrate infra consumers from wrappers to tk directly
fix: replace panic with ResponseWrapper in files API controller
refactor: add fileClerk field to repo structs for locality of behavior
refactor: eliminate tk type aliases from all layers
chore: add context to every src/ dir
fix: use tk input reader for api
chore: update go and deps

0.2.7.1 - 2025/10/31
fix: remove / from database and runtime hx-get/post to respect base href
fix: add / to footer fragment
chore: update go version

0.2.7 - 2025/10/02
feat: allow account update via username
fix: autologin layout adjustments
fix: custom db connection params
fix: allow for extract .tar.gz on file manager #277
fix: add "delete index.html message" to index.html #278
fix: .trash being shown as file on file manager #279

0.2.6 - 2025/06/06
feat: runtime php run
feat: implement clearable fields for cron and mappings
feat: add dash link to default index
feat: prevent only account deletion
fix: add debug logs to ssl watchdog
fix(api): rename accountId to operatorAccountId
feat(ui): safe prefill user and pass on login via query params
fix(ui): missing capital letters on modals
fix(ui): use UiToolset instead of local utilities
fix(ui): add loading until login finishes redirect
fix(ui): add file.name to file manager update file content
fix(ui): resize code editor when modal is resized
fix(ui): keep file line in active state when selected
fix(ui): code editor height when full screen
fix(ui): php update settings missing vhost

0.2.5 - 2025/05/27
refactor(ui): sidebar
refactor(ui): merge page and presenters
refactor(ui): move state.js to individual embeds
feat(ui): add cloak and loading overlay (from goinfinite/ui)
fix(improvement): lazy load marketplace and services avatar images
fix(bug): database service is creating a default mapping [#261]

0.2.4 - 2025/05/08
feat!: add pagination to database read ops
feat: mappings security rules
feat: add should upgrade insecure requests to mappings
feat: add default index html page
fix: allow multi instances of multi nature services
fix: gorm sqlite memory file spec
fix: typo on vhost type select input
fix: hash gen on ssl ownership check
chore: install (and use) ui and tk projects

0.2.3 - 2025/04/08
refactor: replace BadgerDB with SQLite in-memory
fix(critical): x-bind misplacement on InputClientSide component
fix: remove / from marketplace temp dir
fix: check type of vhost on aliasesHostname ToEntity()
tests: add tests for transientDbSvc

0.2.2 - 2025/04/04
refactor: ssl infra implementation
refactor: ssl watchdog after ssl infra refactor
feat(front): add visual clues on ssl list pages
feat(front): add ca bundle field on import ssl modal
feat: issue valid SSL endpoint and UI button
fix(front): swap self-signed ssl bug
fix: skip aliases on ssl watchdog
fix: vhosts key on ssl create pair
fix: add error enums to auth query
fix: add head routes to public api routes
fix: disable default super admin for first account
fix: remove altNames from create ssl pair
fix: stop uninstall of databases on mktplace uninstall
fix: only remove services without mappings on mktplace uninstall
fix: improve first setup with presentation helper
docs: improve project readme

0.2.1 - 2025/04/01
refactor: vhost and mapping infra implementation
refactor: ssl watchdog
refactor: delete service mappings
feat: add mapping hostname and path to service mapping auto create
feat: add support for wildcard vhosts
feat: add custom response code to url mappings
feat: add marketplace item reference to mappings
fix: aliases addition replacing parent ssl
fix: remove mappings and ssl files when vhost is deleted
fix: move vhost aliases to parent row
fix: remove vhost removal from ssl pair

0.2.0 - 2025/03/25
refactor(front): file manager with HTMX+Alpine.js
fix: missing service name on mapping targetValue input

0.1.9 - 2025/03/10
fix(front): file manager permissions and account sudoers
feat: add WorkingDir property to CreateInstallableService DTO

0.1.8 - 2025/01/24
refactor(front): overview page with HTMX+Alpine.js
feat: add time to the overview chart tooltip
feat: allow for o11y overview without resource usage
fix: file manager download files bug
fix: update service table on events
fix: fix start/stop button colors
fix: cron routing from previous ui
fix: replace the storage usage area chart with a line chart
fix: remove unixFile SRI and use record details instead
fix: reorder the OS system info tags

0.1.7 - 2025/01/06
fix(front): install params section on marketplace deploy
refactor(front): login page with HTMX+Alpine.js
feat: /setup/ (api and front)

0.1.6 - 2024/12/23
refactor(front): crons page with HTMX+Alpine.js

0.1.5 - 2024/12/06
feat: add all missing security records to write ops
fix: download files with more than 5MB in size
fix: correctly mapping dataFieldName to ensure the marketplace item installation api route works as expected

0.1.4 - 2024/12/06
feat: move marketplace registry to its own git repository
feat: move services registry to its own git repository

0.1.3 - 2024/12/06
feat: manage ssh keys for accounts
feat: add all missing security records to write ops

0.1.2 - 2024/10/26
refactor(front): marketplace page with HTMX+Alpine.js

0.1.1 - 2024/10/22
refactor(front): accounts page with HTMX+Alpine.js

0.1.0 - 2024/10/08
refactor!: migrate/rename speedia to infinite
refactor: os-api banner
refactor: scheduled tasks
feat: add footer bar with resource usage and scheduled tasks
fix: uptime with proc/1 on overview

0.0.9 - 2024/10/04
refactor(front): runtime page with HTMX+Alpine.js
refactor(front): ssls page with HTMX+Alpine.js
feat: chown default dirs after service install/add
feat: add jsonAjax helper
feat: add dev-build.sh script and make file
fix: adjust mapping layout for lower resolutions
fix: bug on match pattern on mappings
fix: bug on error level and error log php update

0.0.8 - 2024/09/23
refactor(front): databases page with HTMX+Alpine.js

0.0.7 - 2024/09/23
refactor(front): mappings page with HTMX+Alpine.js
feat: opensearch and java support
feat: add OpenMage & Adobe Commerce to marketplace
fix: combine install url with mappings path properly

0.0.6 - 2024/08/22
feat: activity records
feat: limit login attempts by ip address

0.0.5 - 2024/08/20
refactor: api and cli controllers to use services layer
feat: add log handler middleware
fix: supervisorctl auth error when using cron
refactor: unify runtime controllers with service layer
refactor: unify services controllers with service layer
refactor: unify authentication controllers with service layer
refactor: unify account controllers with service layer
refactor: unify cron controllers with service layer
refactor: unify database controllers with service layer
refactor: unify o11y controllers with service layer
refactor: unify ssl controllers with service layer
refactor: unify vhost controllers with service layer
refactor: vos to new format and remove all panics

0.0.4 - 2024/07/17
feat: add async tasks
refactor: marketplace presentation layer

0.0.3 - 2024/07/12
refactor: everything services related
fix: move /_/api to /api

0.0.2 - 2024/07/04
refactor: vhost infra with models
refactor: replace with valid ssl
feat: add restart status to services
feat: accept pkcs1 format for ssl private key
feat: uninstallSteps and files for mktplace
feat: .htaccess auto reload
fix: cron not working
fix: add read only mode middleware
fix: services deactivate loop on read only mode
fix: rename assets of php to php-webserver
fix: .htaccess not being loaded
fix: add missing wp permalink structure
fix: use sni on validation hash curl
fix: is publicly trusted renewal logic
fix: cpu usage calculation us bug

0.0.1 - 2024/05/23
feat: initial release
```
