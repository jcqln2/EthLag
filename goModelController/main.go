package main

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	publicNodeURL = "https://ethereum-rpc.publicnode.com"
	alchemyURL    = "https://eth-mainnet.g.alchemy.com/v2/88eZBls2st3aenXrIVk4p" // Ensure this has a valid API key
)

var rpcBody = []byte(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`)

type Model struct {
	PublicNodeLatency float64
	AlchemyLatency    float64
	PublicNodeGauge   prometheus.Gauge
	AlchemyGauge      prometheus.Gauge
}

type Controller struct {
	Model *Model
}

func NewModel() *Model {
	return &Model{
		PublicNodeGauge: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "eth_rpc_publicnode_latency_seconds",
			Help: "Latency of eth_blockNumber to Public Node",
		}),
		AlchemyGauge: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "eth_rpc_alchemy_latency_seconds",
			Help: "Latency of eth_blockNumber to Alchemy",
		}),
	}
}

func (c *Controller) MeasureLatency(url string, gauge prometheus.Gauge) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(rpcBody))
	if err != nil {
		log.Printf("Error creating request to %s: %v", url, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error calling %s: %v", url, err)
		return
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response from %s: %v", url, err)
		return
	}
	latency := time.Since(start).Seconds()
	gauge.Set(latency) // Ensure gauge is updated
	log.Printf("Latency for %s: %v seconds", url, latency)
}

func (c *Controller) UpdateModel() {
	c.MeasureLatency(publicNodeURL, c.Model.PublicNodeGauge)
	c.MeasureLatency(alchemyURL, c.Model.AlchemyGauge)
}

func main() {
	model := NewModel()
	controller := &Controller{Model: model}

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8000", nil) // Run server in a goroutine

	for {
		controller.UpdateModel()
		time.Sleep(30 * time.Second)
	}
}
