package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func readData(u *url.URL, httpClient *http.Client) (map[string]stats, error) {
	resp, err := httpClient.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve HTTP data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error code received: %s (%d)", http.StatusText(resp.StatusCode), resp.StatusCode)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read body text: %v", err)
	}

	var data map[string]stats
	err = json.Unmarshal(buf, &data)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal json: %v", err)
	}

	return data, nil
}

func writeData(data map[string]stats, iClient influxdb2.Client, database string) {
	writeAPI := iClient.WriteAPIBlocking("", database)

	for key, row := range data {
		p := influxdb2.NewPoint(key, map[string]string{"unit": "random"},
			map[string]interface{}{
				"alt_suppressed":          row.AltSuppressed,
				"messages":                row.Messages,
				"local_samples_processed": row.Local.SamplesProcessed,
				"local_samples_dropped":   row.Local.SamplesDropped,
				"local_modead":            row.Local.ModeAC,
				"local_modes":             row.Local.ModeS,
				"local_bad":               row.Local.Bad,
				"local_unknown_icao":      row.Local.UnknownICAO,
				"local_accepted_n":        row.Local.Accepted[0],
				"local_accepted_bits":     row.Local.Accepted[1],
				"local_signal":            row.Local.Signal,
				"local_noise":             row.Local.PeakSignal,
				"local_strong_signals":    row.Local.StrongSignals,
				"cpr_surface":             row.CPR.Surface,
				"cpr_airborne":            row.CPR.Airborne,
				"cpr_global_ok":           row.CPR.GlobalOK,
				"cpr_global_bad":          row.CPR.GlobalBad,
				"cpr_global_badrange":     row.CPR.GlobalBadRange,
				"cpr_global_badspeed":     row.CPR.GlobalBadSpeed,
				"cpr_global_skipped":      row.CPR.GlobalSkipped,
				"cpr_local_ok":            row.CPR.LocalOK,
				"cpr_local_acft_rel":      row.CPR.LocalAcftRel,
				"cpr_local_recv_rel":      row.CPR.LocalRecvRel,
				"cpr_local_skip":          row.CPR.LocalSkip,
				"cpr_local_skip_range":    row.CPR.LocalSkipRange,
				"cpr_local_skip_speed":    row.CPR.LocalSkipSpeed,
				"cpr_filtered":            row.CPR.Filtered,
				"tracks_all":              row.Tracks.All,
				"tracks_single_msg":       row.Tracks.SingleMessage,
			},
			time.Now())

		writeAPI.WritePoint(context.Background(), p)
	}
}

func main() {
	// Parse parameters
	dump1090URL, _ := url.Parse("http://fr24.in.sffreak.com/dump1090/data/stats.json")
	influxURL, _ := url.Parse("http://dockerhost.in.sffreak.com:8086")
	influxToken := ""
	influxDB := "dump1090"
	pollTime, _ := time.ParseDuration("10s")

	if os.Getenv("HOST") != "" {
		if u, err := url.Parse(os.Getenv("HOST")); err != nil {
			log.Fatalf("%s is not a valid URL, quitting...", os.Getenv("HOST"))
		} else {
			dump1090URL = u
		}
	}

	if os.Getenv("INFLUX_URL") != "" {
		if u, err := url.Parse(os.Getenv("INFLUX_URL")); err != nil {
			log.Fatalf("%s is not a valid InfluxDB URL, quitting...", os.Getenv("INFLUX_URL"))
		} else {
			influxURL = u
		}
	}

	if os.Getenv("INFLUX_TOKEN") != "" {
		influxToken = os.Getenv("INFLUX_TOKEN")
	}

	if os.Getenv("INFLUX_DB") != "" {
		influxDB = os.Getenv("INFLUX_DB")
	}

	if os.Getenv("POLL_TIME") != "" {
		var err error
		pollTime, err = time.ParseDuration(os.Getenv("POLL_TIME"))
		if err != nil {
			log.Fatalf("Unable to parse duration: %v", err)
		}
	}

	// Set up InfluxDB
	iClient := influxdb2.NewClient(influxURL.String(), influxToken)
	defer iClient.Close()

	log.Println("Program started")

	for {
		go func() {
			data, err := readData(dump1090URL, http.DefaultClient)
			if err != nil {
				log.Printf("Cannot read data: %v\n", err)
			}

			writeData(data, iClient, influxDB)
		}()

		time.Sleep(pollTime)
	}
}
