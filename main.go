package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/lukasCoppens/openhab-item-exporter/openhab"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	stateMetric *prometheus.GaugeVec
	errorMetric prometheus.Counter
)

func main() {
	logrus.Info("Starting openhab-item-exporter...")
	cnf := loadConfig()
	initLogging(cnf.LogLevel)
	initMetrics()
	logrus.Info("Creating client")
	openhab := openhab.InitClient(cnf.OpenhabEndpoint, cnf.OpenhabClientTimoutSeconds)
	logrus.Info("Loading initial state")
	err := syncState(openhab)
	if err != nil {
		logrus.Fatalf("failed to load initial state: %s", err.Error())
	}
	logrus.Info("Starting periodic sync")
	go runPeriodicSync(openhab, cnf.OpenhabPollIntervalSeconds)
	logrus.Info("Starting server")
	http.Handle(cnf.MetricsPath, promhttp.Handler())
	http.ListenAndServe(cnf.MetricsPort, nil)
}

func initMetrics() {
	logrus.Infof("Init metrics")
	stateMetric = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "openhab_item_state",
			Help: "State of the item",
		},
		[]string{"name", "state", "type", "groups", "tags"},
	)
	errorMetric = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "openhab_item_exporter_errors",
			Help: "Metric increased when an error is encountered",
		},
	)
}

func syncState(client *openhab.Client) error {
	items, err := client.GetItems()
	if err != nil {
		return err
	}
	logrus.Infof("Found %d items", len(items))
	stateMetric.Reset()
	for _, item := range items {
		intState := item.GetIntState()
		stateMetric.With(prometheus.Labels{"name": item.Name, "state": item.State, "type": item.Type, "groups": strings.Join(item.GroupNames, ","), "tags": strings.Join(item.Tags, ",")}).Set(float64(intState))
	}
	return nil
}

func runPeriodicSync(client *openhab.Client, seconds int) {
	uptimeTicker := time.NewTicker(time.Duration(seconds) * time.Second)
	for {
		<-uptimeTicker.C
		logrus.Infof("Syncing state")
		err := syncState(client)
		if err != nil {
			logrus.Errorf("Failed syncing state: %s", err.Error())
			errorMetric.Inc()
		}
	}
}

func initLogging(ll string) {
	logrus.Infof("Setting logLevel to %s", ll)
	logLevel, err := logrus.ParseLevel(ll)
	if err != nil {
		logrus.Warnf("Failed to parse log level (%s), using default Info. Error: %s", ll, err.Error())
	}
	logrus.SetLevel(logLevel)
}
