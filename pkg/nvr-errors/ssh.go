package nvr_errors

import "errors"

var (
	ErrNoRemoteHost         = errors.New("remote host must be specified")
	ErrNoRemotePort         = errors.New("remote port must be specified")
	ErrNoRemoteUser         = errors.New("remote user must be specified")
	ErrNoKnownHostKeyFound  = errors.New("no known host key found")
	ErrKnownHostKeyMismatch = errors.New("known host key mismatch")
)
