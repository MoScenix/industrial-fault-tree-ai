package mtl

import (
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
		panic(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", MetricPort)
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
		panic(err)
	}
	server.RegisterShutdownHook(func() {
		r.Deregister(registerInfo)
	})
	http.Handle("/metrics", promhttp.HandlerFor(Registry, promhttp.HandlerOpts{}))
	go http.ListenAndServe(MetricPort, nil)
}
