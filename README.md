# wc-fuse-client

The in-pod half of wc-fuse: the control-socket wire protocol and the client
the harness embeds to request and release per-project workspace mounts.

It is split out of the host-side `wc-fuse` repo so the pod's dependency graph
stays small — this module depends only on `wc-log`, never on FUSE or host
provisioning code.

## Contents

- `protocol/` — the JSON `Request`/`Response` messages and the op/error codes.
- `client/` — `client.Client`, with `Mount`, `Unmount`, and `Status`.

## Usage

```go
import "github.com/WywySenarios/wc-fuse-client/client"

c := client.New("/dev/wc/control/wc-fuse.sock", logger)
mountPoint, err := c.Mount(sessionID, projectRoot)
// ... use the workspace at mountPoint ...
err = c.Unmount(projectRoot)
```

Requests carry no policy: they name only a session and a project root. The host
daemon authorizes them against the host-side project-roots allowlist.

## Consuming

Consumed by the harness (in-pod) and by the `wc-fuse` daemon (host), both by
module path. Development:

```sh
go test ./...
```
