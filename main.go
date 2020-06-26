package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var berta_gauge = promauto.NewGauge(
	prometheus.GaugeOpts{
		Name: "berta_free_spots",
		Help: "Free spots available"})

func get_free_spots() int {
	pattern, err := regexp.Compile(`\((\d*)\sfreie`)

	url := os.Getenv("BERTA_URL")
	resp, err := http.Get(url)

	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	match := pattern.FindStringSubmatch(string(body))

	result, err := strconv.Atoi(match[1])

	berta_gauge.Set(float64(result))

	return result
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":7999", nil)

	for {
		go get_free_spots()
		time.Sleep(time.Second * 60)
	}
}
