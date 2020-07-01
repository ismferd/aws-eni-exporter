package main

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)



//Define a struct for you collector that contains pointers
//to prometheus descriptors for each metric you wish to expose.
//Note you can also include fields of other types if they provide utility
//but we just won't be exposing them as metrics.
type eniCollector struct {
	gauge *prometheus.Desc
	eni eni

}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func newEniCollector(eni eni) *eniCollector {
	return &eniCollector{
		gauge: prometheus.NewDesc(strings.ReplaceAll(eni.name,"-","_"),
			"Shows whether a foo has occurred in our cluster",
			nil, nil,
		),
		eni: eni,
	}
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *eniCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.gauge
}

//Collect implements required collect function for all promehteus collectors
func (collector *eniCollector) Collect(ch chan<- prometheus.Metric) {
	eni := AWSconfiguration()
	value := 0
	for _,v := range(eni){
		if v.name == collector.eni.name{
			value = v.ipsAvailable
		}
	}
	ch <- prometheus.MustNewConstMetric(collector.gauge, prometheus.GaugeValue, float64(value))
}