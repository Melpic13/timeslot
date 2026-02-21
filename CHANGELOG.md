# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog,
and this project adheres to Semantic Versioning.

## [Unreleased]
### Added
- Enterprise-quality test suite expansion across all packages and examples
- Coverage gate tooling with `make coverage-check` (90% minimum)
- Release CLI entrypoint at `cmd/timeslot`

### Changed
- CI pipeline now enforces `go mod tidy` cleanliness, race tests, lint, and security scans
- Booking system example now computes next Monday dynamically to avoid date drift regressions
- GoReleaser configuration now builds from `cmd/timeslot`

## [1.0.0] - 2025-08-10
### Added
- Initial public release of core scheduling packages
- Weekly availability, provider, query, conflict, recurrence, and iCal support
- CI/CD workflows and example applications
