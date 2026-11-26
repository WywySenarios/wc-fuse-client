// Package client is the in-pod client for the wc-fuse control socket.
package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/WywySenarios/wc-fuse/protocol"
	wclog "github.com/WywySenarios/wc-log/log"
)

// Client talks to the host wc-fuse daemon over the control socket.
type Client struct {
	socket string
	logger *wclog.Logger
}

// New returns a client for socket. logger may be nil.
func New(socket string, logger *wclog.Logger) *Client {
	return &Client{socket: socket, logger: logger}
}

// Request a workspace mount and return its mount point.
func (c *Client) Mount(sessionID, projectRoot string) (string, error) {
	resp, err := c.do(protocol.Request{
		Op:          protocol.OpMount,
		SessionID:   sessionID,
		ProjectRoot: projectRoot,
	})
	if err != nil {
		return "", err
	}
	c.log("wc-fuse: mounted %s for session %s at %s", projectRoot, sessionID, resp.MountPoint)
	return resp.MountPoint, nil
}

// Unmount releases sessionID's hold on projectRoot's workspace mount.
func (c *Client) Unmount(sessionID, projectRoot string) error {
	_, err := c.do(protocol.Request{
		Op:          protocol.OpUnmount,
		SessionID:   sessionID,
		ProjectRoot: projectRoot,
	})
	if err != nil {
		return err
	}
	c.log("wc-fuse: unmounted %s for session %s", projectRoot, sessionID)
	return nil
}

// Returns the daemon's open-mount count.
func (c *Client) Status() (int, error) {
	resp, err := c.do(protocol.Request{Op: protocol.OpStatus})
	if err != nil {
		return 0, err
	}
	return resp.OpenMounts, nil
}

func (c *Client) do(req protocol.Request) (*protocol.Response, error) {
	if c.socket == "" {
		return nil, errors.New("wc-fuse client: no socket path configured")
	}
	conn, err := net.Dial("unix", c.socket)
	if err != nil {
		return nil, fmt.Errorf("wc-fuse client: dial %s: %w", c.socket, err)
	}
	defer conn.Close()

	if err := json.NewEncoder(conn).Encode(req); err != nil {
		return nil, fmt.Errorf("wc-fuse client: send %s: %w", req.Op, err)
	}
	var resp protocol.Response
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		return nil, fmt.Errorf("wc-fuse client: read %s response: %w", req.Op, err)
	}
	if !resp.OK {
		return nil, &protocol.Error{Code: resp.Code, Message: resp.Error}
	}
	return &resp, nil
}

func (c *Client) log(format string, args ...any) {
	if c.logger != nil {
		c.logger.Info(fmt.Sprintf(format, args...))
	}
}
