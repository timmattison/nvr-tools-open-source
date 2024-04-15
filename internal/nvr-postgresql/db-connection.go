package nvr_postgresql

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/url"
)

const driverName = "postgres"

func GetPostgresqlDbUrl(username string, password string, host string, port int, database string, disableSsl bool) *url.URL {
	rawQuery := ""

	if disableSsl {
		rawQuery = "sslmode=disable"
	}

	dbURL := &url.URL{
		Scheme:   driverName,
		User:     url.UserPassword(username, password),
		Host:     fmt.Sprintf("%s:%d", host, port),
		Path:     fmt.Sprintf("/%s", database),
		RawQuery: rawQuery,
	}

	return dbURL
}

func OpenPostgresqlDbSql(dbURL *url.URL) (*sql.DB, error) {
	var err error
	var db *sql.DB

	if db, err = sql.Open(driverName, dbURL.String()); err != nil {
		return nil, err
	}

	return db, nil
}

func OpenPostgresqlDbSqlx(dbURL *url.URL) (*sqlx.DB, error) {
	var err error
	var db *sqlx.DB

	if db, err = sqlx.Connect(driverName, dbURL.String()); err != nil {
		return nil, err
	}

	return db, nil
}
