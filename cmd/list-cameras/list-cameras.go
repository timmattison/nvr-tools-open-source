package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
	"github.com/timmattison/nvr-tools-open-source/pkg/nvr-unifi-protect"
)

func main() {
	var unifiProtectHost string

	flag.StringVar(&unifiProtectHost, "h", "", "The UniFi Protect host to connect to (IP address or hostname)")

	flag.Parse()

	if unifiProtectHost == "" {
		flag.Usage()
		return
	}

	ctx, cancelFunc := context.WithCancelCause(context.Background())

	var tunneledDbSqlx *sqlx.DB
	var err error

	if tunneledDbSqlx, err = nvr_unifi_protect.GetTunneledUnifiProtectDbSqlx(ctx, cancelFunc, unifiProtectHost); err != nil {
		log.Fatal("Failed to get tunneled DB connection", "error", err)
	}

	var cameraRecords []nvr_unifi_protect.CameraRecord

	if cameraRecords, err = nvr_unifi_protect.SelectCameras(tunneledDbSqlx); err != nil {
		log.Fatal("Failed to select cameras", "error", err)
	}

	var jsonCameras []byte

	if jsonCameras, err = json.MarshalIndent(cameraRecords, "", "  "); err != nil {
		log.Fatal("Failed to marshal cameras", "error", err)
	}

	fmt.Println(string(jsonCameras))

	return
}
