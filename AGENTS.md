# Repository Guidelines

## Project Structure & Module Organization

- Root module: `go.ipao.vip/gen` (Go 1.18+).
- Packages: core `.go` files at repo root; support packages in `internal/`, `field/`, `types/`, `helper/`, `tools/`.
- Tests: colocated `*_test.go` next to implementation (e.g., `do_test.go`, `generator_test.go`).

## Build, Test, and Development Commands

- `go build ./...`: build all packages.
- `go test ./...`: run all tests.
- `go test -race -cover ./...`: race detector + coverage.
- `go vet ./...`: static checks.
- `golangci-lint run`: lint and format per `.golangci.yml`.

## Coding Style & Naming Conventions

- Formatting: Go defaults. Run `golangci-lint run` or `go fmt ./...` before pushing.
- Imports: organized with `goimports` (configured in `.golangci.yml`).
- Indentation: tabs (standard Go). Line length: keep readable; avoid overly long expressions.
- Naming: exported identifiers use `PascalCase`; unexported `camelCase`; package names lower-case, no underscores. Test files named `*_test.go`.

## Testing Guidelines

- Framework: standard `testing` package.
- Test names: `TestXxx` and table-driven where useful.
- Coverage: prefer meaningful tests; run `go test -cover ./...` locally.
- Add tests alongside code changes; keep tests deterministic (no network/DB by default).

## Commit & Pull Request Guidelines

- Commits: concise, imperative subject (e.g., `fix: handle nil columns` or `feat(generator): add where clause builder`).
- Scope small, logical changes; include rationale in body when non-obvious.
- PRs: clear description, link related issues, list behavior changes, and include test updates. Add examples if APIs change.

## Security & Configuration Tips

- Do not commit secrets or credentials. Use environment variables for local experiments.
- Run `go vet` and `golangci-lint` to catch common issues.
- Prefer safe APIs and validate inputs in code generation paths.

## Architecture Overview (Brief)

- Codegen and parsing live in `internal/`.
- SQL expression helpers in `field/`; nullable helpers in `types/`.
- Utilities in `helper/`; auxiliary tooling under `tools/`.
