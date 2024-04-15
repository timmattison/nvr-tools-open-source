package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/timmattison/nvr-tools-open-source/internal/nvr-environment"
	"github.com/timmattison/nvr-tools-open-source/pkg/nvr-unifi-protect"
	"os"
)

func main() {
	var unifiProtectHost string

	flag.StringVar(&unifiProtectHost, "h", "", fmt.Sprintf("The UniFi Protect host to connect to (IP address or hostname) (can also be set via %s environment variable or .env file value)", nvr_environment.UnifiProtectHostKey))

	flag.Parse()

	if unifiProtectHost == "" {
		if err := godotenv.Load(".env"); err != nil {
			log.Warn("No .env file found and no host specified")
		} else {
			unifiProtectHost = os.Getenv(nvr_environment.UnifiProtectHostKey)
		}

		if unifiProtectHost == "" {
			flag.Usage()
			return
		}
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
