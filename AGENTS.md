# AGENTS.md

Read and follow `.agents/AGENTS.md` immediately — it is the style book and supersedes anything here on conflict.
The style book lives in `.agents/` (gitignored); project-specific conventions live here (tracked).

A nearer `AGENTS.md` applies to its directory and children.

## Before Changes

- Read `README.md`.
- Read `docs/DEVELOPMENT.md` before build, test, UI, or API work.
- Read `docs/SECURITY.md` before changing authentication, secrets, filesystem access,
  command execution, or network behavior.
- Read the nearest `.context.md` before opening code in a directory.
- Use `docs/FEATURE-MAP.md` when tracing or changing a user-facing flow.

## Rules

- Preserve clean architecture. `src/domain/` must not import `src/infra/` or
  `src/presentation/`.
- Define repository interfaces in `src/domain/repository/`. Implement them in
  `src/infra/`.
- Keep application workflows in domain use cases. Keep domain invariants in domain
  types. Keep API controllers, CLI handlers, and UI presenters focused on input and
  output.
- Validate and map external data at the boundary. Do not pass raw request,
  environment, database, file, or external-service data into domain logic.
- Route internal database access through `src/infra/internalDatabase/`. Use the
  service that matches the data lifetime.
- Edit source files, not generated files. Edit `.templ` files and companion state
  files instead of `*_templ.go`; run `templ generate -path src/presentation/ui`.
- Edit Swagger annotations in `src/presentation/api/api.go` and
  `src/presentation/api/controller/`. Do not edit `src/presentation/api/docs/`; run
  `swag init --pdl 3 -g src/presentation/api/api.go -o src/presentation/api/docs`.
- Format changed Go files with `gofmt`.
- Use the tool versions in `.mise.toml`. Run `mise trust` once and `mise install`
  when they are missing.
- Run `golangci-lint run` for Go changes.
- Tests must verify behavior. Prefer the containerized test procedure in
  `docs/DEVELOPMENT.md` because application tests can change system state. Never
  publish a test image built from a context that contains `.env` or other secrets.
- Verify rendered UI changes in a browser.
- Update `.context.md`, `docs/FEATURE-MAP.md`, `README.md`, or other documentation
  when the corresponding structure, user-facing flow, or public behavior changes.
- Keep `.env` local. Never commit or publish secrets.
