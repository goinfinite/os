# Changelog

```log
0.1.5 - 2024/X/X
feat: add all missing security records to write ops
fix: download files with more than 5MB in size

0.1.4 - 2024/X/X
feat: move marketplace registry to its own git repository
feat: move services registry to its own git repository

0.1.3 - 2024/X/X
Nothing yet

0.1.2 - 2024/X/X
refactor(front): marketplace page with HTMX+Alpine.js

0.1.1 - 2024/X/X
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
