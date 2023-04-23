package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &solarCollector{}

type solarCollector struct {
	CurrentPower *prometheus.Desc
	YieldToday *prometheus.Desc
	YieldTotal *prometheus.Desc

	stats func() ([]SolarStats, error)
}

func NewSolarCollector(stats func() ([]SolarStats, error)) prometheus.Collector {
	return &solarCollector {
		CurrentPower: prometheus.NewDesc(
			"solarman_now_p",
			"Amount of wattage currently produced in Watt",
			[]string{},
			nil,
		),
		YieldToday: prometheus.NewDesc(
			"solarman_today_e",
			"Amount of wattage produced today in kWh",
			[]string{},
			nil,
		),
		YieldTotal: prometheus.NewDesc(
			"solarman_total_e",
			"Amount of wattage produced in total in kWh",
			[]string{},
			nil,
		),
		stats: stats,
	}
}

func (c *solarCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.CurrentPower,
	}

	for _, d := range ds {
		ch <- d
	}
}

func (c *solarCollector) Collect(ch chan<- prometheus.Metric) {
	stats, err := c.stats()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.CurrentPower, err)
		return
	}

	for _, s := range stats {
		ch <- prometheus.MustNewConstMetric(
			c.CurrentPower,
			prometheus.GaugeValue,
			s.CurrentPower,
		)

		ch <- prometheus.MustNewConstMetric(
			c.YieldToday,
			prometheus.GaugeValue,
			s.YieldToday,
		)

		ch <- prometheus.MustNewConstMetric(
			c.YieldTotal,
			prometheus.GaugeValue,
			s.YieldTotal,
		)
	}
}