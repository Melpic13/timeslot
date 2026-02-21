# Contributing

## Development setup

1. Install Go 1.22+
2. Clone the repository
3. Run `go test ./...`

## Standard local workflow

1. `make fmt`
2. `make test`
3. `make test-race`
4. `make coverage-check`
5. `make vet`

## Code style

- Follow standard Go conventions
- Keep API docs on exported symbols
- Prefer table-driven tests for behavior coverage
- Keep the package dependency graph acyclic and dependency-light

## Testing requirements

- Add tests for any feature or bug fix
- Keep overall coverage at or above 90%
- Add/update fuzz tests for external input parsers where relevant

## PR process

1. Open an issue if scope or design is non-trivial
2. Submit a focused PR with tests
3. Ensure all CI jobs pass before requesting review

## Issue guidelines

Use bug/feature templates under `.github/ISSUE_TEMPLATE`.
