package service

import (
	"context"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcgraph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
	"github.com/cloudwego/hertz/pkg/app"
)

type CreateGraphVersionService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCreateGraphVersionService(Context context.Context, RequestContext *app.RequestContext) *CreateGraphVersionService {
	return &CreateGraphVersionService{RequestContext: RequestContext, Context: Context}
}

func (h *CreateGraphVersionService) Run(req *graph.CreateGraphVersionRequest) (resp *graph.BaseResponseString, err error) {
	res, err := rpc.GraphClient.CreateVersion(h.Context, &rpcgraph.CreateVersionReq{
		GraphId: req.GraphId, VersionName: req.VersionName,
	})
	if err != nil {
		return &graph.BaseResponseString{Code: 1, Message: err.Error()}, err
	}
	return &graph.BaseResponseString{Code: 0, Message: "success", Data: res.Version}, nil
}
