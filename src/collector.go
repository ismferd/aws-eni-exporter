package main

import (
	"log"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

//Struct to collect eni
//to prometheus descriptors for each metric you wish to expose.
type eniCollector struct {
	gauge *prometheus.Desc
	eni   eni
}

//newEniCollector must create a constructor for the ENI collector
func newEniCollector(eni eni) *eniCollector {
	return &eniCollector{
		gauge: prometheus.NewDesc(strings.ReplaceAll(eni.name, "-", "_"),
			"Shows the currect value of ip availables for each subnet",
			nil, nil,
		),
		eni: eni,
	}
}

//Describe fucntion to ENI collector
func (collector *eniCollector) Describe(ch chan<- *prometheus.Desc) {

	ch <- collector.gauge
}

//Implements the ENI collector
func (collector *eniCollector) Collect(ch chan<- prometheus.Metric) {
	awsConfig := AWSconfiguration(region, vpc)
	cli := newAWSClient(awsConfig)
	eni, err := cli.GetSubnets(vpc)
	if err != nil {
		log.Fatal(err)
	}

	value := 0
	for _, v := range eni {
		if v.name == collector.eni.name {
			value = v.ipsAvailable
		}
	}
	ch <- prometheus.MustNewConstMetric(collector.gauge, prometheus.GaugeValue, float64(value))
}
