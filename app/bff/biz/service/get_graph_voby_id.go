package service

import (
	"context"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcgraph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetGraphVOByIdService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetGraphVOByIdService(Context context.Context, RequestContext *app.RequestContext) *GetGraphVOByIdService {
	return &GetGraphVOByIdService{RequestContext: RequestContext, Context: Context}
}

func (h *GetGraphVOByIdService) Run(req *graph.GetGraphVOByIdRequest) (resp *graph.BaseResponseGraphVO, err error) {
	res, err := rpc.GraphClient.GetGraph(h.Context, &rpcgraph.GetGraphReq{Id: req.Id})
	if err != nil {
		return &graph.BaseResponseGraphVO{Code: 1, Message: err.Error()}, err
	}
	return &graph.BaseResponseGraphVO{
		Code:    0,
		Message: "success",
		Data: &graph.GraphVO{
			Id:             res.Graph.Id,
			GraphName:      res.Graph.GraphName,
			Description:    res.Graph.Description,
			Cover:          res.Graph.Cover,
			UserId:         res.Graph.UserId,
			CurrentVersion: res.Graph.CurrentVersion,
			HasTmp:         res.Graph.HasTmp,
			CreateTime:     res.Graph.CreateTime,
			UpdateTime:     res.Graph.UpdateTime,
		},
	}, nil
}
