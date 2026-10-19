package client_test

import (
	"encoding/json"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WywySenarios/wc-fuse-client/client"
	"github.com/WywySenarios/wc-fuse/protocol"
)

// fakeDaemon speaks the control protocol over a unix socket. It stands
// in for the host daemon so this module's tests do not depend on the
// wc-fuse (host) module.
type fakeDaemon struct {
	allow []string
	open  int
}

func newFakeDaemon(t *testing.T, allow ...string) string {
	t.Helper()
	socket := filepath.Join(t.TempDir(), "ctl", "wc-fuse.sock")
	if err := os.MkdirAll(filepath.Dir(socket), 0o711); err != nil {
		t.Fatal(err)
	}
	ln, err := net.Listen("unix", socket)
	if err != nil {
		t.Fatal(err)
	}
	d := &fakeDaemon{allow: allow}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go d.serve(conn)
		}
	}()
	t.Cleanup(func() { _ = ln.Close() })
	return socket
}

func (d *fakeDaemon) serve(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)
	for {
		var req protocol.Request
		if err := dec.Decode(&req); err != nil {
			return
		}
		var resp protocol.Response
		switch req.Op {
		case protocol.OpMount:
			if !d.allowed(req.ProjectRoot) {
				resp = protocol.Response{Code: protocol.CodeDenied, Error: "not allowed"}
			} else {
				d.open++
				resp = protocol.Response{OK: true, MountPoint: "/dev/wc/workspace/x", OpenMounts: d.open}
			}
		case protocol.OpUnmount:
			if d.open > 0 {
				d.open--
			}
			resp = protocol.Response{OK: true, OpenMounts: d.open}
		case protocol.OpStatus:
			resp = protocol.Response{OK: true, OpenMounts: d.open}
		default:
			resp = protocol.Response{Code: protocol.CodeBadRequest, Error: "unknown op"}
		}
		if err := enc.Encode(resp); err != nil {
			return
		}
	}
}

func (d *fakeDaemon) allowed(root string) bool {
	for _, a := range d.allow {
		if root == a || strings.HasPrefix(root, a+"/") {
			return true
		}
	}
	return false
}

func TestClientRoundTrip(t *testing.T) {
	socket := newFakeDaemon(t, "/home/eric/src")
	c := client.New(socket, nil)

	mp, err := c.Mount("s1", "/home/eric/src/proj")
	if err != nil {
		t.Fatalf("Mount: %v", err)
	}
	if mp == "" {
		t.Error("Mount returned an empty mount point")
	}

	if n, err := c.Status(); err != nil || n != 1 {
		t.Errorf("Status = %d, %v; want 1, nil", n, err)
	}

	if err := c.Unmount("/home/eric/src/proj"); err != nil {
		t.Fatalf("Unmount: %v", err)
	}
	if n, err := c.Status(); err != nil || n != 0 {
		t.Errorf("Status after unmount = %d, %v; want 0, nil", n, err)
	}
}

func TestClientRefusesDisallowedRoot(t *testing.T) {
	socket := newFakeDaemon(t, "/home/eric/src")
	c := client.New(socket, nil)

	if _, err := c.Mount("s1", "/etc"); err == nil {
		t.Fatal("Mount of a root outside the allowlist must fail")
	}
}

func TestClientDaemonErrorIsTyped(t *testing.T) {
	socket := newFakeDaemon(t, "/home/eric/src")
	c := client.New(socket, nil)

	_, err := c.Mount("s1", "/etc")
	if err == nil {
		t.Fatal("Mount of a root outside the allowlist must fail")
	}
	var perr *protocol.Error
	if !errors.As(err, &perr) {
		t.Fatalf("error is %T, want *protocol.Error", err)
	}
	if perr.Code != protocol.CodeDenied {
		t.Errorf("code = %q, want %q", perr.Code, protocol.CodeDenied)
	}
}
