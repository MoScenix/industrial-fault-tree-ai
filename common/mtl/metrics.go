package mtl

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Registry *prometheus.Registry

func InitMetric(serviceName, MetricPort, registerAddr string) {
	Registry = prometheus.NewRegistry()
	Registry.MustRegister(collectors.NewGoCollector())
	Registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	r, err := consul.NewConsulRegister(registerAddr)
	if err != nil {
		log.Printf("init prometheus consul register failed: %v", err)
		startMetricHTTP(MetricPort)
		return
	}
	addr, err := net.ResolveTCPAddr("tcp", MetricPort)
	if err != nil {
		fmt.Printf("resolve metric addr failed: %v", err)
		log.Printf("resolve metric addr failed: %v", err)
		startMetricHTTP(MetricPort)
		return
	}
	registerInfo := &registry.Info{
		ServiceName: "prometheus",
		Weight:      1,
		Addr:        addr,
		Tags: map[string]string{
			"service": serviceName,
		},
	}
	err = r.Register(registerInfo)
	if err != nil {
		log.Printf("register prometheus service failed: %v", err)
	}
	server.RegisterShutdownHook(func() {
		r.Deregister(registerInfo)
	})
	startMetricHTTP(MetricPort)
}

func startMetricHTTP(metricPort string) {
	http.Handle("/metrics", promhttp.HandlerFor(Registry, promhttp.HandlerOpts{}))
	go func() {
		if err := http.ListenAndServe(metricPort, nil); err != nil {
			log.Printf("start metrics http failed: %v", err)
		}
	}()
}
