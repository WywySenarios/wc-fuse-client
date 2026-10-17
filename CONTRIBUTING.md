# Contributing to wc-fuse-client

## Layout

- `protocol/`, the control-socket JSON messages, shared with the host daemon.
- `client/`
- `scripts/tests/`, test suites

## Testing

`scripts/tests/` suites must pass before commit. pre-commit runs them.

```sh
pre-commit run
```

It is highly advised to run tests before making modifications to set a baseline.

## Consuming contract

The host `wc-fuse` module and the harness consume this module by module path (GOPROXY at a published tag). `wc-log` is consumed the same way.
