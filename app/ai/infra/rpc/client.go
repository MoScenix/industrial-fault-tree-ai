package rpc

import (
	"github.com/MoScenix/industrial-fault-tree-ai/app/ai/conf"
	"github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/document/documentservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
)

var DocumentClient documentservice.Client

func Init() {
	initDocumentClient()
}

func initDocumentClient() {
	r := mustResolver()
	var err error
	DocumentClient, err = documentservice.NewClient(
		"document",
		client.WithResolver(r),
		client.WithMetaHandler(transmeta.MetainfoClientHandler),
		client.WithSuite(tracing.NewClientSuite()),
	)
	if err != nil {
		klog.Fatal(err)
	}
}

func mustResolver() discovery.Resolver {
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		klog.Fatal(err)
	}
	return r
}
