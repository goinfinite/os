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

## Unit Tests

- Run Go tests inside the container described in `docs/DEVELOPMENT.md`, never on the
  host: application tests change system state, and host results are not trustworthy.
- Re-test the working tree against the already-built image with the bind-mounted
  procedure in the same document, so a one-shot edit does not cost a fresh image.
- When a package fails, check the same package at the base commit before treating it
  as your regression. Report pre-existing failures instead of fixing them inside an
  unrelated change.

## UI Verification

- Run `bash dev-build.sh http-unpriv --pid-only` yourself to bring up a disposable
  instance for dashboard checks. Do not stop to ask for permission.
- Start it detached with `setsid`, so a foreground timeout cannot reach its process
  group. Wait for port 1618. The script builds the image, creates the `dev` account
  with the password `abc123!`, then keeps rebuilding on every code change.
- Stop it by signalling the pid recorded in `logs/dev-build.pid`. Its trap stops the
  container, the rebuild watcher and the attach session, then removes the file. The
  container runs with `--rm`, so it leaves nothing behind.
- Do not use the `http` argument: it changes host sysctl settings through `sudo`.
