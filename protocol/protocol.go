// Package protocol defines the JSON messages exchanged over the wc-fuse control socket.
package protocol

// Operations.
const (
	// OpMount requests a workspace mount for a project root.
	OpMount = "mount"
	// OpUnmount releases a project's workspace mount.
	OpUnmount = "unmount"
	// OpStatus reports the daemon's open-mount count.
	OpStatus = "status"
)

// Error codes.
const (
	CodeBadRequest = "bad_request"
	CodeDenied     = "denied"
	CodeInternal   = "internal"
)

// Request is one control-socket request.
type Request struct {
	Op          string `json:"op"`
	SessionID   string `json:"session_id,omitempty"`
	ProjectRoot string `json:"project_root,omitempty"`
}

// Response is one control-socket response.
type Response struct {
	OK         bool   `json:"ok"`
	MountPoint string `json:"mount_point,omitempty"`
	OpenMounts int    `json:"open_mounts"`
	Code       string `json:"code,omitempty"`
	Error      string `json:"error,omitempty"`
}
