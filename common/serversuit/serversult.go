package serversuit

import (
	"github.com/MoScenix/industrial-fault-tree-ai/common/mtl"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	prometheus "github.com/kitex-contrib/monitor-prometheus"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
)

type CommonServerSuilt struct {
	ServerName   string
	RegistryAddr string
}

func (s *CommonServerSuilt) Options() []server.Option {
	opts := []server.Option{
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: s.ServerName,
		}),
		server.WithMetaHandler(transmeta.ServerHTTP2Handler),
		server.WithMetaHandler(transmeta.MetainfoServerHandler),
		server.WithTracer(
			prometheus.NewServerTracer("", "", prometheus.WithDisableServer(true), prometheus.WithRegistry(mtl.Registry)),
		),
		server.WithSuite(tracing.NewServerSuite()),
	}
	reg, err := consul.NewConsulRegister(s.RegistryAddr)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithRegistry(reg))
	return opts
}
