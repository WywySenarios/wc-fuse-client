# wc-fuse-client

This repo is part of the Wywy repo line and is subject to Wywy conventions.

The in-pod client and wire protocol for wc-fuse: the control-socket messages and the client in-pod workloads embed to request and release per-project workspace mounts.

This repository is separate from the host-side `wc-fuse` repo so the pod's dependency graph stays small — it depends only on `wc-log`, never on FUSE or host provisioning code.

## Contents

- `protocol/` — the JSON `Request`/`Response` messages and the op/error codes.
- `client/` — `client.Client`, with `Mount`, `Unmount`, and `Status`.

## Usage

```go
import "github.com/WywySenarios/wc-fuse-client/client"

c := client.New("/dev/wc/control/wc-fuse.sock", logger)
mountPoint, err := c.Mount(sessionID, projectRoot)
// ... use the workspace at mountPoint ...
err = c.Unmount(sessionID, projectRoot) // sessionID must be the one used for Mount
```

The `sessionID` passed to `Unmount` must match the one used for `Mount`.

Requests carry no policy: they name only a session and a project root. The host daemon authorizes them against the host-side project-roots allowlist.

## Consuming

Consumed by in-pod workloads and by the `wc-fuse` daemon (host), both by module path.

## Development

```sh
go test ./...
```
