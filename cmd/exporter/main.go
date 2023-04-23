package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	exporter "github.com/BenjaminGlusa/prometheus-solarman-exporter"
	"github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	var (
		solarHost = flag.String("solar.host", "Solaranlage", "address of the solarman microinverter website")
		solarUser = flag.String("solar.user", "admin", "username for the solarman microinverter website")
		solarPass = flag.String("solar.pass", "admin", "password for the solarman microinverter website")
		promPort = flag.Int("prom.port", 9150, "port to expose prometheus metrics")
	)

	flag.Parse()

	uri := fmt.Sprintf("http://%s/status.html", *solarHost)

	solarStats := func()([]exporter.SolarStats, error) {
		netClient := &http.Client{
			Timeout: time.Second * 10,
		}
		req, _ := http.NewRequest("GET", uri, nil)
		req.SetBasicAuth(*solarUser, *solarPass)

		resp, err := netClient.Do(req)
		if err != nil {
			log.Fatalf("could not scrap solarman on %s: %s", uri, err)
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("could not ready body %s", err)
		}
		r := bytes.NewReader(body)
		return exporter.ScanSolarStats(r)

	}

	sc := exporter.NewSolarCollector(solarStats)

	reg := prometheus.NewRegistry()
	reg.MustRegister(sc)
	// reg.MustRegister(collectors.NewGoCollector())

	mux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", promHandler)

	port := fmt.Sprintf(":%d", *promPort)
	log.Printf("starting solarman exporter on port %q/metrics", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("cannot start solarman exporter: %s", err.Error())
	}
}