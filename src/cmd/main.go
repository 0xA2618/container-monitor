package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"github.com/zhangyy8lab/tusimaServerMonitor/client"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// var activeServer []string
var monitorServerList []string

type metrics struct {
	serverTotal prometheus.Gauge
	hdFailures  *prometheus.CounterVec
}

// NewMetrics creates a new metrics instance.
func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		serverTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "alive_server_count",
			Help: "check alive server",
		}),
		hdFailures: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "custom_check_alive_server",
				Help: "custom_check_alive_server",
			},
			monitorServerList,
		),
	}

	reg.MustRegister(m.serverTotal)
	reg.MustRegister(m.hdFailures)
	return m
}

// periodicallyUpdateMetrics
func periodicallyUpdateMetrics(labels prometheus.Labels, m *metrics) {

	ticker := time.NewTicker(time.Second * 5)

	for range ticker.C {
		lines := client.DockerPS()

		// int to float
		activeServerStr := strconv.Itoa(len(lines) - 1)

		activeServerFloat, _ := strconv.ParseFloat(activeServerStr, 64)

		if len(lines) > 0 {
			for _, name := range monitorServerList {
				labels[name] = client.CheckServetActive(name, lines)

			}
		}
		m.serverTotal.Set(activeServerFloat)
		m.hdFailures.With(labels)
	}

	return
}

func main() {
	config := viper.New()
	//config.AddConfigPath("/home/ubuntu/zyy/tusimaMonitorServer/src/config/")
	config.AddConfigPath("/app/config/")
	config.SetConfigName("service")
	config.SetConfigType("yaml")

	if err := config.ReadInConfig(); err != nil {
		panic(err)

		log.Fatal(err)
	}

	for _, name := range config.GetStringSlice("monitorServer.server") {
		monitorServerList = append(monitorServerList, strings.ReplaceAll(name, "-", "_"))
		fmt.Println("monitorServerList:", monitorServerList)
	}
	fmt.Println("monitorServerList2:", monitorServerList)
	fmt.Println("----------------")
	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	// Create new metrics and register them using the custom registry.
	m := NewMetrics(reg)

	labels := make(prometheus.Labels)

	go periodicallyUpdateMetrics(labels, m)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.Get("server.port")), nil))
}
