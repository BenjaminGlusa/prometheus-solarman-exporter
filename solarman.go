package exporter

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

type SolarStats struct {
	CurrentPower float64
	YieldToday float64
	YieldTotal float64
}

func ScanSolarStats(r io.Reader)([]SolarStats, error) {
	s := bufio.NewScanner(r)

	var stats []SolarStats
	var solarStats SolarStats

	for s.Scan() {
		fields := strings.Fields(string(s.Bytes()))
		if len(fields) >= 3 && fields[1] == "webdata_now_p" {
			s := strings.ReplaceAll(fields[3], "\"", "")
			s = strings.ReplaceAll(s, ";", "")
			currentPower, err := strconv.Atoi(s)
			if err != nil {
				log.Fatalf("Could not parse current power %s : %s", s, err.Error())
			}
			solarStats.CurrentPower = float64(currentPower)
		}

		if len(fields) >= 3 && fields[1] == "webdata_today_e" {
			s := strings.ReplaceAll(fields[3], "\"", "")
			s = strings.ReplaceAll(s, ";", "")
			yieldToday, err := strconv.ParseFloat(s, 64)
			if err != nil {
				log.Fatalf("Could not parse yield today %s : %s", s, err.Error())
			}
			solarStats.YieldToday = float64(yieldToday)
		}

		if len(fields) >= 3 && fields[1] == "webdata_total_e" {
			s := strings.ReplaceAll(fields[3], "\"", "")
			s = strings.ReplaceAll(s, ";", "")
			yieldTotal, err := strconv.ParseFloat(s, 64)
			if err != nil {
				log.Fatalf("Could not parse yield total %s : %s", s, err.Error())
			}
			solarStats.YieldToday = float64(yieldTotal)
		}

	}

	stats = append(stats, solarStats)

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("%w: failed to read metrics", err)
	}

	return stats, nil
}