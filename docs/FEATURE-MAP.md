# Feature Map

> Auto-maintained index of every user-facing feature and the code path that implements it. Updated alongside the code — not after the fact.

## Login & Authentication

User authentication via username/password, creating session tokens for subsequent API/UI requests.

**Flow:**

1. `src/presentation/api/controller/authentication.go` — POST /api/v1/auth/login receives credentials
2. `src/presentation/ui/presenter/login/` — UI login form submission
3. `src/domain/useCase/createSessionToken.go` — Validates credentials and generates JWT token
4. `src/infra/auth/` — Token creation and signing
5. `src/infra/internalDatabase/` — Persists and retrieves account credentials

---

## Account Management

Create, read, update, and delete user accounts. Create and manage secure access public keys for API authentication.

**Flow:**

1. `src/presentation/api/controller/account.go` — Account REST endpoints (GET, POST, PUT, DELETE)
2. `src/presentation/cli/controller/account.go` — CLI account commands
3. `src/presentation/ui/presenter/accounts/` — Web dashboard account management pages
4. `src/domain/useCase/createAccount.go` — Account creation logic with validation
5. `src/domain/useCase/updateAccount.go` — Account update operations
6. `src/domain/useCase/deleteAccount.go` — Account deletion with cleanup
7. `src/domain/useCase/createSecureAccessPublicKey.go` — Generate API keys
8. `src/infra/account/` — Account repository implementation
9. `src/infra/internalDatabase/` — Persistent storage of accounts and keys

---

## Database Management

Create, configure, and delete relational databases (PostgreSQL, MySQL, MariaDB). Manage database users and their credentials.

**Flow:**

1. `src/presentation/api/controller/database.go` — Database REST endpoints (GET, POST, DELETE)
2. `src/presentation/cli/controller/database.go` — CLI database commands
3. `src/presentation/ui/presenter/databases/` — Web dashboard database management UI
4. `src/domain/useCase/createDatabase.go` — Database creation orchestration
5. `src/domain/useCase/deleteDatabase.go` — Database deletion and cleanup
6. `src/domain/useCase/createDatabaseUser.go` — Add user to database
7. `src/domain/useCase/deleteDatabaseUser.go` — Remove user from database
8. `src/infra/database/` — Multi-engine database service implementations
9. `src/infra/internalDatabase/` — Metadata persistence

---

## Service Deployment (Marketplace)

Browse and install applications and services from the marketplace catalog. One-click deployments with configurable parameters.

**Flow:**

1. `src/presentation/api/controller/marketplace.go` — Marketplace REST endpoints
2. `src/presentation/cli/controller/marketplace.go` — CLI marketplace commands
3. `src/presentation/ui/presenter/marketplace/` — Web dashboard marketplace UI
4. `src/domain/useCase/createInstallableService.go` — Marketplace app installation
5. `src/domain/useCase/createCustomService.go` — Custom service creation
6. `src/infra/marketplace/` — Marketplace catalog and installation orchestration
7. `src/infra/services/` — Service lifecycle management
8. `src/infra/internalDatabase/` — Service metadata and installation history

---

## Application Mappings

Configure domain/path mappings to route traffic to deployed applications. Secure mappings with rate limits, bandwidth caps, and connection limits.

**Flow:**

1. `src/presentation/api/controller/virtualHost.go` — Mapping REST endpoints
2. `src/presentation/cli/controller/virtualHost.go` — CLI mapping commands
3. `src/presentation/ui/presenter/mappings/` — Web dashboard mapping configuration UI
4. `src/domain/useCase/createMapping.go` — Mapping creation and routing setup
5. `src/domain/useCase/createMappingSecurityRule.go` — Add rate/bandwidth/connection limits
6. `src/infra/vhost/` — NGINX virtual host configuration generation
7. `src/infra/internalSetup/` — Web server reload after config changes
8. `src/infra/internalDatabase/` — Mapping metadata persistence

---

## Scheduled Tasks (Cron)

Create, configure, and execute scheduled background tasks on a recurring schedule.

**Flow:**

1. `src/presentation/api/controller/cron.go` — Cron REST endpoints (GET, POST, PUT, DELETE)
2. `src/presentation/cli/controller/cron.go` — CLI cron commands
3. `src/presentation/ui/presenter/crons/` — Web dashboard cron management UI
4. `src/domain/useCase/createCron.go` — Cron job creation with schedule validation
5. `src/domain/useCase/updateCron.go` — Cron modification
6. `src/domain/useCase/deleteCron.go` — Cron deletion
7. `src/infra/cron/` — Cron repository implementation
8. `src/infra/scheduledTask/` — Task scheduler and execution engine
9. `src/infra/internalDatabase/` — Cron schedule and execution log persistence

---

## SSL/TLS Certificates

Manage SSL certificates including Let's Encrypt integration for automated HTTPS. Support both self-signed and publicly trusted certificates.

**Flow:**

1. `src/presentation/api/controller/ssl.go` — SSL REST endpoints
2. `src/presentation/cli/controller/ssl.go` — CLI SSL commands
3. `src/presentation/ui/presenter/ssls/` — Web dashboard SSL management UI
4. `src/domain/useCase/createPubliclyTrustedSslPair.go` — Let's Encrypt certificate provisioning
5. `src/infra/ssl/` — SSL certificate management and renewal
6. `src/infra/helper/` — Certificate generation utilities (self-signed)
7. `src/presentation/http.go` — HTTPS server setup with certificates
8. `src/infra/internalDatabase/` — Certificate metadata and tracking

---

## File Manager

Browse, upload, download, delete, and compress files within the container filesystem.

**Flow:**

1. `src/presentation/api/controller/files.go` — File REST endpoints (GET, POST, DELETE)
2. `src/presentation/cli/controller/files.go` — CLI file commands
3. `src/presentation/ui/presenter/fileManager/` — Web dashboard file browser UI
4. `src/domain/useCase/copyUnixFile.go` — File copy operation
5. `src/domain/useCase/compressUnixFiles.go` — Archive creation
6. `src/domain/useCase/deleteUnixFile.go` — File deletion
7. `src/infra/files/` — File repository implementation
8. `src/infra/helper/` — Filesystem utilities and path validation

---

## Runtime Management

Configure programming language runtimes (PHP, Node.js, Python) including version selection and environment variables.

**Flow:**

1. `src/presentation/api/controller/runtime.go` — Runtime REST endpoints
2. `src/presentation/cli/controller/runtime.go` — CLI runtime commands
3. `src/presentation/ui/presenter/runtimes/` — Web dashboard runtime configuration UI
4. `src/domain/useCase/createInstallableService.go` — Runtime installation as a service
5. `src/infra/runtime/` — Runtime environment setup and version management
6. `src/infra/services/` — Service lifecycle (start/stop runtimes)
7. `src/infra/internalDatabase/` — Runtime configuration persistence

---

## System Overview & Monitoring

Display system health, hardware specs, uptime, IP address, and operational metrics. Quick access to common actions.

**Flow:**

1. `src/presentation/api/controller/o11y.go` — System metrics REST endpoint
2. `src/presentation/ui/presenter/overview/` — Web dashboard system overview page
3. `src/infra/o11y/` — System observability and metrics collection
4. `src/infra/internalDatabase/` — Historical metric storage

---

## Initial Setup Wizard

One-time setup flow to create the first admin account and configure system defaults.

**Flow:**

1. `src/presentation/api/controller/setup.go` — Setup REST endpoints (gated on account count)
2. `src/presentation/ui/presenter/setup/` — Web dashboard multi-step setup form
3. `src/domain/useCase/createFirstAccount.go` — Initial admin account creation
4. `src/infra/internalSetup/` — Initial web server configuration
5. `src/infra/internalDatabase/` — First account persistence

---

## Virtual Host Configuration

Automatically generate and manage NGINX configuration for hosting multiple applications and mappings.

**Flow:**

1. `src/infra/vhost/` — Virtual host configuration generator and manager
2. `src/infra/internalSetup/` — Web server reload and NGINX integration
3. `src/presentation/api/controller/virtualHost.go` — API endpoints for vhost management
4. `src/infra/internalDatabase/` — Vhost configuration metadata

---

## Activity Auditing

Log all user actions, system events, and administrative operations for compliance and debugging.

**Flow:**

1. `src/infra/activityRecord/` — Activity record repository implementation
2. `src/presentation/api/middleware/` — Request logging middleware captures API calls
3. `src/presentation/ui/middleware/` — Request logging middleware captures UI actions
4. `src/presentation/cli/middleware/` — Command execution logging
5. `src/infra/internalDatabase/` — Trail database stores complete audit history

---

## Service Lifecycle Management

Deploy, configure, start, stop, and remove user applications and system services.

**Flow:**

1. `src/presentation/api/controller/marketplace.go` — Service installation endpoints
2. `src/presentation/api/controller/services.go` — Service lifecycle endpoints
3. `src/presentation/cli/controller/services.go` — CLI service commands
4. `src/domain/useCase/createInstallableService.go` — Service installation orchestration
5. `src/domain/useCase/deleteService.go` — Service removal
6. `src/infra/services/` — Service deployment, configuration, and lifecycle
7. `src/infra/database/` — Database service management (PostgreSQL, MySQL, etc.)
8. `src/infra/runtime/` — Runtime service management
9. `src/infra/internalDatabase/` — Service metadata and configuration

---

## Request Authentication & Authorization

Verify user identity and permissions for all API and UI requests via JWT tokens or session cookies.

**Flow:**

1. `src/presentation/api/middleware/` — API Bearer token and API key verification
2. `src/presentation/ui/middleware/` — UI session cookie verification
3. `src/infra/auth/` — JWT token validation
4. `src/infra/internalDatabase/` — Account lookup and permission verification
5. `src/domain/useCase/createSessionToken.go` — Token generation on login

---

## Scheduled Task Execution

Execute background tasks at configured schedules, logging execution history and results.

**Flow:**

1. `src/infra/scheduledTask/` — Task scheduler and execution engine
2. `src/infra/cron/` — Cron job repository
3. `src/presentation/api/controller/scheduledTask.go` — Task execution tracking REST endpoint
4. `src/infra/internalDatabase/` — Task schedule and execution log persistence

---

## Secure API Key Management

Create, revoke, and manage API keys for secure programmatic access to the system.

**Flow:**

1. `src/presentation/api/controller/account.go` — API key management endpoints
2. `src/domain/useCase/createSecureAccessPublicKey.go` — Key generation
3. `src/infra/auth/` — Key validation middleware
4. `src/infra/internalDatabase/` — API key storage and retrieval

---

## Responsive Web Dashboard

Serve embedded static assets (CSS, JavaScript, images) for the web UI with responsive design.

**Flow:**

1. `src/presentation/ui/` — UI initialization and routing
2. `src/presentation/ui/router.go` — Route registration for all pages
3. `src/presentation/ui/assets/` — Embedded static files (Tailwind CSS, JS)
4. `src/presentation/ui/layout/` — Page layout templates
5. `src/presentation/ui/component/` — Reusable UI components
6. `src/presentation/ui/presenter/` — Page-specific presenters
7. `src/presentation/ui/middleware/` — Asset serving and authentication
8. `src/infra/internalDatabase/` — Page data queries

</context>
