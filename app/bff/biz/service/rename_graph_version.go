package service

import (
	"context"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcgraph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
	"github.com/cloudwego/hertz/pkg/app"
)

type RenameGraphVersionService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewRenameGraphVersionService(Context context.Context, RequestContext *app.RequestContext) *RenameGraphVersionService {
	return &RenameGraphVersionService{RequestContext: RequestContext, Context: Context}
}

func (h *RenameGraphVersionService) Run(req *graph.RenameGraphVersionRequest) (resp *graph.BaseResponseBoolean, err error) {
	if _, err := loadAuthorizedGraphRecord(h.Context, req.GraphId); err != nil {
		return &graph.BaseResponseBoolean{Code: 1, Message: graphAccessError(err).Error()}, nil
	}
	res, err := rpc.GraphClient.RenameVersion(h.Context, &rpcgraph.RenameVersionReq{
		GraphId: req.GraphId, Version: req.Version, VersionName: req.VersionName,
	})
	if err != nil {
		return &graph.BaseResponseBoolean{Code: 1, Message: err.Error()}, err
	}
	return &graph.BaseResponseBoolean{Code: 0, Message: "success", Data: res.Success}, nil
}
