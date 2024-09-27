package nvr_ssh

import (
	"bytes"
	"context"
	"fmt"
	"github.com/timmattison/nvr-tools-open-source/pkg/nvr-errors"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"strings"
)

var privateKeyLocations = []string{
	".ssh/id_rsa",
	".ssh/id_dsa",
	".ssh/id_ecdsa",
	".ssh/id_ed25519",
}

var knownHostsLocations = []string{
	".ssh/known_hosts",
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
	} else {
		var knownHosts []byte

		separator := ""

		for _, knownHostsLocation := range knownHostsLocations {
			knownHostsPath := home + string(os.PathSeparator) + knownHostsLocation

			var currentKnownHosts []byte

			if currentKnownHosts, err = os.ReadFile(knownHostsPath); err != nil {
				// Error reading the known hosts file, ignoring it
				continue
			}

			knownHosts = append(knownHosts, separator...)
			separator = "\n"
			knownHosts = append(knownHosts, currentKnownHosts...)
		}

		hostKeys := make(map[string]ssh.PublicKey)

		for len(knownHosts) > 0 {
			var hosts []string
			var pubKey ssh.PublicKey

			if _, hosts, pubKey, _, knownHosts, err = ssh.ParseKnownHosts(knownHosts); err != nil {
				return nil, err
			}

			for _, host := range hosts {
				hostKeys[host] = pubKey
			}
		}

		hostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			addressFromRemote := strings.Split(remote.String(), ":")[0]

			if hostKey, ok := hostKeys[addressFromRemote]; ok {
				if bytes.Equal(hostKey.Marshal(), key.Marshal()) {
					return nil
				}

				return nvr_errors.ErrKnownHostKeyMismatch
			}

			addressFromHostname := strings.Split(hostname, ":")[0]

			if hostKey, ok := hostKeys[addressFromHostname]; ok {
				if bytes.Equal(hostKey.Marshal(), key.Marshal()) {
					return nil
				}

				return nvr_errors.ErrKnownHostKeyMismatch
			}

			return nvr_errors.ErrNoKnownHostKeyFound
		}
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
