package nvr_unifi_protect

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type LicensePlateRecord struct {
	Start        int64  `db:"start"`
	End          int64  `db:"end"`
	LicensePlate string `db:"license_plate"`
}

type LicensePlateWithLocalTime struct {
	Start        time.Time
	End          time.Time
	LicensePlate string
}

func SelectLicensePlates(db *sqlx.DB) ([]LicensePlateWithLocalTime, error) {
	var licensePlates []LicensePlateRecord

	query := "SELECT start, \"end\", metadata->'licensePlate'->>'name' as license_plate FROM events WHERE metadata::jsonb @? '$.licensePlate.name' AND \"end\" IS NOT NULL;"

	if err := db.Select(&licensePlates, query); err != nil {
		return nil, err
	}

	var licensePlatesWithLocalTime []LicensePlateWithLocalTime

	for i := range licensePlates {
		licensePlatesWithLocalTime = append(licensePlatesWithLocalTime, LicensePlateWithLocalTime{
			Start:        time.Unix(licensePlates[i].Start/1000, 0),
			End:          time.Unix(licensePlates[i].End/1000, 0),
			LicensePlate: licensePlates[i].LicensePlate,
		})
	}

	return licensePlatesWithLocalTime, nil
}
