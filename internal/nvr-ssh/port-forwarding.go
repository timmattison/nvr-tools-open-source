package nvr_ssh

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"math/rand"
	"net"
)

func ForwardPort(ctx context.Context, cancelFunc context.CancelCauseFunc, sshClient *ssh.Client, destinationHost string, destinationPort int, localListeningIp string) (string, int, error) {
	var listener net.Listener
	var err error

	forwardedPort := rand.Intn(32768) + 10000

	if localListeningIp == "" {
		localListeningIp = "127.0.0.1"
	}

	if listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", localListeningIp, forwardedPort)); err != nil {
		return "", 0, err
	}

	go waitForConnectionsAndForward(ctx, cancelFunc, sshClient, destinationHost, destinationPort, listener)

	return localListeningIp, forwardedPort, nil
}

func waitForConnectionsAndForward(ctx context.Context, cancelFunc context.CancelCauseFunc, sshClient *ssh.Client, destinationHost string, destinationPort int, listener net.Listener) {
	defer listener.Close()
	var err error

	for {
		var local net.Conn

		if local, err = listener.Accept(); err != nil {
			return
		}

		// Issue a dial to the remote server on our SSH sshClient; here "127.0.0.1" refers to the remote server
		var remote net.Conn

		if remote, err = sshClient.Dial("tcp", fmt.Sprintf("%s:%d", destinationHost, destinationPort)); err != nil {
			return
		}

		go runTunnel(ctx, cancelFunc, local, remote)
	}
}

func runTunnel(ctx context.Context, cancelFunc context.CancelCauseFunc, local net.Conn, remote net.Conn) {
	defer local.Close()
	defer remote.Close()

	go socketCopy(ctx, cancelFunc, local, remote)
	go socketCopy(ctx, cancelFunc, remote, local)

	<-ctx.Done()
}

func socketCopy(ctx context.Context, cancelFunc context.CancelCauseFunc, connection1 net.Conn, connection2 net.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		var err error

		if _, err = io.Copy(connection1, connection2); err != nil {
			cancelFunc(err)
			return
		}
	}
}
