package service

import (
	"context"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcgraph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
	"github.com/cloudwego/hertz/pkg/app"
)

type StartEditService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewStartEditService(Context context.Context, RequestContext *app.RequestContext) *StartEditService {
	return &StartEditService{RequestContext: RequestContext, Context: Context}
}

func (h *StartEditService) Run(req *graph.StartEditRequest) (resp *graph.BaseResponseGraphEditState, err error) {
	if _, err := loadAuthorizedGraphRecord(h.Context, req.GraphId); err != nil {
		return &graph.BaseResponseGraphEditState{Code: 1, Message: graphAccessError(err).Error()}, nil
	}
	res, err := rpc.GraphClient.StartEdit(h.Context, &rpcgraph.StartEditReq{
		GraphId: req.GraphId,
		Version: req.Version,
	})
	if err != nil {
		return &graph.BaseResponseGraphEditState{Code: 1, Message: err.Error()}, err
	}
	return &graph.BaseResponseGraphEditState{
		Code:    0,
		Message: "success",
		Data: &graph.GraphEditState{
			GraphId:        req.GraphId,
			TmpReady:       res.TmpReady,
			BasedOnVersion: res.BasedOnVersion,
			Message:        res.Message,
		},
	}, nil
}
