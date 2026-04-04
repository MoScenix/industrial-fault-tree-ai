package service

import (
	"context"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcgraph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
	"github.com/cloudwego/hertz/pkg/app"
)

type ListGraphMessageService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListGraphMessageService(Context context.Context, RequestContext *app.RequestContext) *ListGraphMessageService {
	return &ListGraphMessageService{RequestContext: RequestContext, Context: Context}
}

func (h *ListGraphMessageService) Run(req *graph.ListGraphMessageRequest) (resp *graph.BaseResponsePageGraphMessageVO, err error) {
	res, err := rpc.GraphClient.ListGraphMessage(h.Context, &rpcgraph.ListGraphMessageReq{
		GraphId: req.GraphId, PageSize: req.PageSize, LastCreateTime: req.LastCreateTime,
	})
	if err != nil {
		return &graph.BaseResponsePageGraphMessageVO{Code: 1, Message: err.Error()}, err
	}
	records := make([]*graph.GraphMessageVO, 0, len(res.MessageList))
	for _, item := range res.MessageList {
		records = append(records, &graph.GraphMessageVO{
			Id: item.Id, GraphId: item.GraphId, UserId: item.UserId, Role: item.Role,
			Content: item.Content, CreateTime: item.CreateTime, UpdateTime: item.UpdateTime,
		})
	}
	return &graph.BaseResponsePageGraphMessageVO{
		Code:    0,
		Message: "success",
		Data:    makeGraphMessagePage(records, req.PageSize, res.Total),
	}, nil
}
