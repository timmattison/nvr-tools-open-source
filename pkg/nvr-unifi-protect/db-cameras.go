package nvr_unifi_protect

import "github.com/jmoiron/sqlx"

type CameraRecord struct {
	Mac  string `db:"mac"`
	Type string `db:"type"`
	Name string `db:"name"`
	Id   string `db:"id"`
}

func SelectCameras(db *sqlx.DB) ([]CameraRecord, error) {
	var cameras []CameraRecord

	query := "SELECT mac, type, name, id FROM cameras"

	if err := db.Select(&cameras, query); err != nil {
		return nil, err
	}

	return cameras, nil
}
