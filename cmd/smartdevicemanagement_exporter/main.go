package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	smartdevicemanagementexporter "github.com/ryancurrah/smartdevicemanagement_exporter"
	"github.com/ryancurrah/smartdevicemanagement_exporter/partnerconnmanager"
)

var addr = flag.String("listen-address", ":8080", "Address to listen on for HTTP requests.")
var pid = flag.String("project-id", "", "ID of the Smart Device Management project.")
var credentials = flag.String("credentials", "client_secret.json", "Location on disk to the Oauth2 credentials JSON file.")
var refreshToken = flag.String("refresh-token", "refresh_token.json", "Location on disk to store the Oauth2 refresh token JSON file.")
var recordMetricsDelay = flag.Duration("record-metrics-delay", time.Second*60, "Delay between queries to the Smart Device Management API for recording metrics.")

func main() {
	flag.Parse()

	addrEnv := os.Getenv("LISTEN_ADDRESS")
	if addrEnv != "" {
		*addr = addrEnv
	}

	pidEnv := os.Getenv("PROJECT_ID")
	if pidEnv != "" {
		*pid = pidEnv
	}

	credentialsEnv := os.Getenv("CREDENTIALS")
	if credentialsEnv != "" {
		*credentials = credentialsEnv
	}

	refreshTokenEnv := os.Getenv("REFRESH_TOKEN")
	if refreshTokenEnv != "" {
		*refreshToken = refreshTokenEnv
	}

	recordMetricsDelayEnv := os.Getenv("RECORD_METRICS_DELAY")

	if recordMetricsDelayEnv != "" {
		recordMetricsDelayDuration, err := time.ParseDuration(recordMetricsDelayEnv)
		if err != nil {
			log.Fatal(err)
		}

		*recordMetricsDelay = recordMetricsDelayDuration
	}

	ctx := context.Background()

	config, err := smartdevicemanagementexporter.LoadOauth2Config(*credentials)
	if err != nil {
		log.Fatal(err)
	}

	authorizationCodeChan := make(chan partnerconnmanager.AuthorizationCode)

	pcm := partnerconnmanager.PartnerConnManager{
		AuthorizationCodeChan: authorizationCodeChan,
		ClientID:              config.ClientID,
		ProjectID:             *pid,
	}

	sdme := smartdevicemanagementexporter.SmartDeviceManagementExporter{
		AuthorizationCodeChan: authorizationCodeChan,
		Config:                config,
		ProjectID:             *pid,
		Ctx:                   ctx,
		RefreshTokenFile:      *refreshToken,
		RecordMetricsDelay:    *recordMetricsDelay,
	}

	go func() {
		err := sdme.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/authstatus", func(w http.ResponseWriter, r *http.Request) {
		authStatus := "not authorized"
		if sdme.IsClientRunning() {
			authStatus = "authorized"
		}

		fmt.Fprint(w, authStatus)
	})
	http.HandleFunc("/authorize", pcm.AuthorizeHandler)
	http.HandleFunc("/authorized", pcm.AuthorizedHandler)

	log.Println("running")

	log.Fatal(http.ListenAndServe(*addr, nil))
}
