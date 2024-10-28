package nvr_unifi_protect

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/timmattison/nvr-tools-open-source/internal/nvr-postgresql"
	"github.com/timmattison/nvr-tools-open-source/internal/nvr-ssh"
	"golang.org/x/crypto/ssh"
	"net/url"
)

const (
	unifiProtectPostgresqlUser     = "postgres"
	unifiProtectPostgresqlPassword = ""
	unifiProtectPostgresqlHost     = "127.0.0.1"
	unifiProtectPostgresqlPort     = 5433
	unifiProtectPostgresqlDatabase = "unifi-protect"

	unifiProtectSshPort  = 22
	unifiProtectUsername = "root"
)

func GetTunneledUnifiProtectDbSql(ctx context.Context, cancelFunc context.CancelCauseFunc, unifiProtectHost string, unifiProtectSshPort int, unifiProtectSshUser string, allowUnverifiedHosts bool) (*sql.DB, error) {
	var dbUrl *url.URL
	var err error

	if dbUrl, err = getTunneledUnifiProtectDbUrl(ctx, cancelFunc, unifiProtectHost, unifiProtectSshPort, unifiProtectSshUser, allowUnverifiedHosts); err != nil {
		return nil, err
	}

	return nvr_postgresql.OpenPostgresqlDbSql(dbUrl)
}

func GetTunneledUnifiProtectDbSqlx(ctx context.Context, cancelFunc context.CancelCauseFunc, unifiProtectHost string, unifiProtectSshPort int, unifiProtectSshUser string, allowUnverifiedHosts bool) (*sqlx.DB, error) {
	var dbUrl *url.URL
	var err error

	if dbUrl, err = getTunneledUnifiProtectDbUrl(ctx, cancelFunc, unifiProtectHost, unifiProtectSshPort, unifiProtectSshUser, allowUnverifiedHosts); err != nil {
		return nil, err
	}

	return nvr_postgresql.OpenPostgresqlDbSqlx(dbUrl)
}

func getTunneledUnifiProtectDbUrl(ctx context.Context, cancelFunc context.CancelCauseFunc, unifiProtectHost string, unifiProtectSshPort int, unifiProtectSshUser string, allowUnverifiedHosts bool) (*url.URL, error) {
	var sshClient *ssh.Client
	var err error

	if sshClient, err = getUnifiProtectSshClient(ctx, cancelFunc, unifiProtectHost, allowUnverifiedHosts); err != nil {
		return nil, err
	}

	var localIp string
	var localPort int

	if localIp, localPort, err = tunnelToUnifiProtectDb(ctx, cancelFunc, sshClient); err != nil {
		return nil, err
	}

	return getUnifiProtectDbUrl(localIp, localPort), nil
}

func getUnifiProtectSshClient(ctx context.Context, cancelFunc context.CancelCauseFunc, unifiProtectHost string, allowUnverifiedHosts bool) (*ssh.Client, error) {
	return nvr_ssh.GetSshClient(ctx, cancelFunc, unifiProtectHost, unifiProtectSshPort, unifiProtectUsername, allowUnverifiedHosts)
}

func tunnelToUnifiProtectDb(ctx context.Context, cancelFunc context.CancelCauseFunc, sshClient *ssh.Client) (string, int, error) {
	return nvr_ssh.ForwardPort(ctx, cancelFunc, sshClient, unifiProtectPostgresqlHost, unifiProtectPostgresqlPort, "")
}

func getUnifiProtectDbUrl(localIp string, localPort int) *url.URL {
	return nvr_postgresql.GetPostgresqlDbUrl(
		unifiProtectPostgresqlUser,
		unifiProtectPostgresqlPassword,
		localIp,
		localPort,
		unifiProtectPostgresqlDatabase,
		true,
	)
}
