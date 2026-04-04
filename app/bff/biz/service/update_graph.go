package service

import (
	"context"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcgraph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateGraphService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateGraphService(Context context.Context, RequestContext *app.RequestContext) *UpdateGraphService {
	return &UpdateGraphService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateGraphService) Run(req *graph.GraphUpdateRequest) (resp *graph.BaseResponseBoolean, err error) {
	res, err := rpc.GraphClient.UpdateGraph(h.Context, &rpcgraph.UpdateGraphReq{
		Id:          req.Id,
		GraphName:   req.GraphName,
		Description: req.Description,
		Cover:       req.Cover,
	})
	if err != nil {
		return &graph.BaseResponseBoolean{Code: 1, Message: err.Error()}, err
	}
	return &graph.BaseResponseBoolean{Code: 0, Message: "success", Data: res.Success}, nil
}
