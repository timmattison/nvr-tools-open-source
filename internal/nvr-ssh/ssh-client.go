package nvr_ssh

import (
	"context"
	"fmt"
	"github.com/timmattison/nvr-tools-open-source/pkg/nvr-errors"
	"golang.org/x/crypto/ssh"
	"os"
)

var privateKeyLocations = []string{
	".ssh/id_rsa", // Linux
	".ssh/id_dsa",
	".ssh/id_ecdsa",
	".ssh/id_ed25519",
}

func GetSshClient(ctx context.Context, cancelFunc context.CancelCauseFunc, remoteHost string, remotePort int, remoteUser string, allowUnverifiedHosts bool) (*ssh.Client, error) {
	var err error

	defer func() {
		if err != nil {
			cancelFunc(err)
		}
	}()

	if remoteHost == "" {
		err = nvr_errors.ErrNoRemoteHost
		return nil, err
	}

	if remotePort == 0 {
		err = nvr_errors.ErrNoRemotePort
		return nil, err
	}

	var sshClientConfig *ssh.ClientConfig

	if sshClientConfig, err = getSshClientConfig(remoteUser, allowUnverifiedHosts); err != nil {
		return nil, err
	}

	var sshClient *ssh.Client

	if sshClient, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", remoteHost, remotePort), sshClientConfig); err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		sshClient.Close()
	}()

	return sshClient, nil
}

func getSshClientConfig(remoteUser string, allowUnverifiedHosts bool) (*ssh.ClientConfig, error) {
	var err error
	var home string

	if remoteUser == "" {
		return nil, nvr_errors.ErrNoRemoteUser
	}

	if home, err = os.UserHomeDir(); err != nil {
		return nil, nil
	}

	var signers []ssh.Signer

	for _, privateKeyLocation := range privateKeyLocations {
		keyPath := home + string(os.PathSeparator) + privateKeyLocation

		var key []byte

		if key, err = os.ReadFile(keyPath); err != nil {
			// Error reading the key file, ignoring it
			continue
		}

		var signer ssh.Signer

		if signer, err = ssh.ParsePrivateKey(key); err != nil {
			// Error parsing the private key, ignoring it
			continue
		}

		signers = append(signers, signer)
	}

	var hostKeyCallback ssh.HostKeyCallback

	if allowUnverifiedHosts {
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	config := &ssh.ClientConfig{
		User: remoteUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signers...),
		},
		HostKeyCallback: hostKeyCallback,
	}

	return config, nil
}
