package service

import (
	"context"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcgraph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
	"github.com/cloudwego/hertz/pkg/app"
)

type DeleteGraphService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDeleteGraphService(Context context.Context, RequestContext *app.RequestContext) *DeleteGraphService {
	return &DeleteGraphService{RequestContext: RequestContext, Context: Context}
}

func (h *DeleteGraphService) Run(req *graph.DeleteRequest) (resp *graph.BaseResponseBoolean, err error) {
	if _, err := loadAuthorizedGraphRecord(h.Context, req.Id); err != nil {
		return &graph.BaseResponseBoolean{Code: 1, Message: graphAccessError(err).Error()}, nil
	}
	res, err := rpc.GraphClient.DeleteGraph(h.Context, &rpcgraph.DeleteGraphReq{Id: req.Id})
	if err != nil {
		return &graph.BaseResponseBoolean{Code: 1, Message: err.Error()}, err
	}
	return &graph.BaseResponseBoolean{Code: 0, Message: "success", Data: res.Success}, nil
}
