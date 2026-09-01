# AGENTS.override.md

## Before Changes

- Read `README.md`.
- Read `docs/DEVELOPMENT.md` before build, test, UI, or API work.
- Read `docs/SECURITY.md` before changing authentication, secrets, filesystem access,
  command execution, or network behavior.
- Read the nearest `.context.md` before opening code in a directory.
- Use `docs/FEATURE-MAP.md` when tracing or changing a user-facing flow.

## Rules

- Edit source files, not generated files. Edit `.templ` files and companion state
  files instead of `*_templ.go`; run `templ generate -path src/presentation/ui`.
- Edit Swagger annotations in `src/presentation/api/api.go` and
  `src/presentation/api/controller/`. Do not edit `src/presentation/api/docs/`; run
  `swag init --pdl 3 -g src/presentation/api/api.go -o src/presentation/api/docs`.
- Use the tool versions in `.mise.toml`. Run `mise trust` once and `mise install`
  when they are missing.
- Tests must verify behavior. Prefer the containerized test procedure in
  `docs/DEVELOPMENT.md` because application tests can change system state.
