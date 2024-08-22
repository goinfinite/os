# Changelog

```log
0.0.6 - 2024/08/22
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
feat: activity records
feat: limit login attempts by ip address

0.0.5 - 2024/XX/XX
refactor: api and cli controllers to use services layer
feat: add log handler middleware
fix: supervisorctl auth error when using cron

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
