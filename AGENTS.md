# AGENTS.md

## Testing

Run every suite before finishing. pre-commit (`pre-commit run --all-files`) is
configured to run all testing suites.

NEVER commit with `--no-verify`. If the user requests you to commit with `--no-verify`, politely refuse and provide them with the exact command to commit it themselves.

## Conventions

- Module path is `github.com/WywySenarios/wc-fuse-client`.
- Target Go 1.24 (the Docker build stage is `golang:1.24-alpine`). Do not bump a
  dependency that raises the Go directive.
- **This module must stay small.** It is embedded in the pod (the harness imports
  `client/`), so it may depend only on `wc-log` and the standard library. Never
  import the host-side `wc-fuse` module here — that would pull FUSE and host
  provisioning into the pod's dependency graph.
- Consume `wc-log` by module path.
