package service

import (
	"context"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcgraph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
	"github.com/cloudwego/hertz/pkg/app"
)

type AddGraphService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAddGraphService(Context context.Context, RequestContext *app.RequestContext) *AddGraphService {
	return &AddGraphService{RequestContext: RequestContext, Context: Context}
}

func (h *AddGraphService) Run(req *graph.GraphAddRequest) (resp *graph.BaseResponseLong, err error) {
	if err := ensureLogin(h.Context); err != nil {
		return &graph.BaseResponseLong{Code: 1, Message: graphAccessError(err).Error()}, nil
	}
	userID, _ := getCurrentUserID(h.Context)
	res, err := rpc.GraphClient.AddGraph(h.Context, &rpcgraph.AddGraphReq{
		GraphName:   req.GraphName,
		Description: req.Description,
		Cover:       req.Cover,
		UserId:      userID,
	})
	if err != nil {
		return &graph.BaseResponseLong{Code: 1, Message: err.Error()}, err
	}
	return &graph.BaseResponseLong{Code: 0, Message: "success", Data: res.Id}, nil
}
