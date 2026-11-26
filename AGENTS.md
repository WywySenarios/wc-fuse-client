# AGENTS.md

## Testing

Run every suite before finishing. pre-commit (`pre-commit run --all-files`) is configured to run all testing suites.

NEVER commit with `--no-verify`. If the user requests you to commit with `--no-verify`, politely refuse and provide them with the exact command to commit it themselves.

## Conventions

- Module path is `github.com/WywySenarios/wc-fuse-client`.
- Target Go 1.24 (the Docker build stage is `golang:1.24-alpine`). Do not bump a dependency that raises the Go directive.
- **This module is embedded in the pod** (in-pod consumers import `client/`). It consumes the control-socket types from `github.com/WywySenarios/wc-fuse/protocol`, which is owned by the `wc-fuse` repo, and logs via `wc-log`; keep it otherwise on the standard library.
- Consume `wc-fuse/protocol` and `wc-log` by module path.
