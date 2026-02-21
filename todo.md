# TimeSlot Enterprise TODO

This document tracks enterprise-readiness work for `github.com/Melpic13/timeslot`.

## Status Snapshot (2026-02-21)

- Overall test coverage: `91.7%`
- Core packages implemented: `slot`, `availability`, `provider`, `query`, `conflict`, `recurrence`, `ical`, `timezone`
- CI automation: enabled (`test`, `lint`, `security`, `fuzz`, `release`, `codeql`)
- Release automation: enabled via GoReleaser and tag-based release workflow

## Milestone Timeline

### 2023-09-15
- Foundation complete (`internal/*`, `slot`, `timezone`, root errors/version/doc)

### 2024-06-21
- Core scheduling engine complete (`availability`, `provider`, `query`, `conflict`)

### 2025-08-10
- Integration and delivery complete (`recurrence`, `ical`, examples, docs, project automation)

### 2025-11-20
- Test maturity raised to `90%+` overall coverage

## Enterprise Readiness Checklist

### Architecture & API
- [x] Zero external runtime dependencies (stdlib only)
- [x] Immutable-style collection and availability operations
- [x] Timezone-aware scheduling primitives
- [x] Typed error surfaces and wrapping support
- [x] Package boundaries and dependency graph remain acyclic

### Quality Engineering
- [x] Unit tests for all major public APIs
- [x] Fuzz tests for parser surfaces (`availability`, `recurrence`, `ical`)
- [x] Benchmarks for key hotspots (`slot`, `provider`, `conflict`, `ical`)
- [x] Coverage at or above 90%
- [x] Race-enabled test execution in CI

### Security & Compliance
- [x] CodeQL workflow
- [x] gosec scan workflow
- [x] Input validation for external/parsing inputs
- [x] Security policy documented

### Developer Experience
- [x] Make targets for test/race/coverage/fuzz
- [x] Coverage threshold enforcement target (`make coverage-check`)
- [x] Contributing guide with local quality workflow
- [x] Working examples for common use cases

### Delivery Automation
- [x] CI quality gates on push/PR
- [x] Tag-driven release workflow
- [x] GoReleaser configuration aligned to a buildable entrypoint
- [x] Keep a Changelog format maintained

## Next Enterprise Backlog

### v1.1 hardening
- [ ] Introduce context-aware APIs for potentially expensive search paths
- [ ] Add API-level SLA benchmarks and regression thresholds in CI
- [ ] Add deterministic clock injection hooks for time-sensitive APIs
- [ ] Add stricter RFC 5545 recurrence compatibility tests
- [ ] Add consumer-facing migration notes and API stability matrix

### v1.2 operations
- [ ] Publish signed release artifacts and provenance attestations
- [ ] Add SBOM generation during release
- [ ] Add dependency and license reporting artifacts in CI

## Standard Verification Commands

```bash
make test
make test-race
make coverage-check
make vet
go test -bench=. -benchmem ./...
```
