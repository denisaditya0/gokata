# Gokata

Keyword-driven API test automation framework built on Go + Ginkgo.

## Modules

| Module | Description |
|--------|-------------|
| `core/` | Test framework (fluent builder, context, tags, retry) |
| `api/` | Backend API — run orchestration, GitLab, Jira (future) |
| `web/` | SvelteKit frontend — scenario editor, run console (future) |

## Quick Start

```bash
cd core
run-tests.bat fast service products sit
run-tests.bat list "products"
```

See [core/README.md](core/README.md) for framework docs.
