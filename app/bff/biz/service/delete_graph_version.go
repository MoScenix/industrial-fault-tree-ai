package service

import (
	"context"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcgraph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
	"github.com/cloudwego/hertz/pkg/app"
)

type DeleteGraphVersionService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDeleteGraphVersionService(Context context.Context, RequestContext *app.RequestContext) *DeleteGraphVersionService {
	return &DeleteGraphVersionService{RequestContext: RequestContext, Context: Context}
}

func (h *DeleteGraphVersionService) Run(req *graph.DeleteGraphVersionRequest) (resp *graph.BaseResponseBoolean, err error) {
	res, err := rpc.GraphClient.DeleteVersion(h.Context, &rpcgraph.DeleteVersionReq{
		GraphId: req.GraphId, Version: req.Version,
	})
	if err != nil {
		return &graph.BaseResponseBoolean{Code: 1, Message: err.Error()}, err
	}
	return &graph.BaseResponseBoolean{Code: 0, Message: "success", Data: res.Success}, nil
}
