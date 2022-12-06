package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/ghodss/yaml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
)

var config Config

const (
	collector = "databricks_exporter"
	RunsTotal = "runs_total"
)

func main() {
	var err error
	var configFile, bind string
	// =====================
	// Get OS parameter
	// =====================
	flag.StringVar(&configFile, "config", "metrics.yaml", "configuration file")
	flag.StringVar(&bind, "bind", "0.0.0.0:9971", "bind")
	flag.Parse()

	// =====================
	// Load config & yaml
	// =====================
	var b []byte
	if b, err = ioutil.ReadFile(configFile); err != nil {
		log.Errorf("Failed to read config file: %s", err)
		os.Exit(1)
	}

	// Load yaml
	if err := yaml.Unmarshal(b, &config); err != nil {
		log.Errorf("Failed to load config: %s", err)
		os.Exit(1)
	}

	// ========================
	// Regist handler
	// ========================
	log.Infof("Regist version collector - %s", collector)
	prometheus.Register(version.NewCollector(collector))
	prometheus.Register(&QueryCollector{})

	// Regist http handler
	log.Infof("HTTP handler path - %s", "/metrics")
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h := promhttp.HandlerFor(prometheus.Gatherers{
			prometheus.DefaultGatherer,
		}, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	})

	// start server
	log.Infof("Starting http server - %s", bind)
	if err := http.ListenAndServe(bind, nil); err != nil {
		log.Errorf("Failed to start http server: %s", err)
	}
}

// Describe prometheus describe
func (e *QueryCollector) Describe(ch chan<- *prometheus.Desc) {
	for metricName, metric := range config.Metrics {
		metric.metricDesc = prometheus.NewDesc(
			prometheus.BuildFQName(collector, "", metricName),
			metric.Description,
			metric.Labels, nil,
		)
		config.Metrics[metricName] = metric
		log.Infof("metric description for \"%s\" registerd", metricName)
	}
}

// Collect prometheus collect
func (e *QueryCollector) Collect(ch chan<- prometheus.Metric) {

	// Connect to database
	log.Infof("Collecting runs...")
	runs, err := GetRuns()
	if err != nil {
		log.Infof("Failed to get runs: %s", err)
	}
	e.collectRunsTotal(runs, ch)
}

func (e *QueryCollector) collectRunsTotal(runs *[]Run, ch chan<- prometheus.Metric) {
	// Execute each queries in metrics
	metric := config.Metrics[RunsTotal]
	for _, run := range *runs {
		// Metric labels
		labelVals, err := run.ToLabelValues(metric.Labels)
		if err != nil {
			log.Errorf("Failed to get label values: %s", err)
		}
		// Metric value
		countRuns := "1"
		val, _ := strconv.ParseFloat(countRuns, 64)
		ch <- prometheus.MustNewConstMetric(metric.metricDesc, prometheus.CounterValue, val, labelVals...)
	}
}
