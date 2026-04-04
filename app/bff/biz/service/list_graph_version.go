package service

import (
	"context"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcgraph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
	"github.com/cloudwego/hertz/pkg/app"
)

type ListGraphVersionService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListGraphVersionService(Context context.Context, RequestContext *app.RequestContext) *ListGraphVersionService {
	return &ListGraphVersionService{RequestContext: RequestContext, Context: Context}
}

func (h *ListGraphVersionService) Run(req *graph.ListGraphVersionRequest) (resp *graph.BaseResponsePageGraphVersionVO, err error) {
	if _, err := loadAuthorizedGraphRecord(h.Context, req.GraphId); err != nil {
		return &graph.BaseResponsePageGraphVersionVO{Code: 1, Message: graphAccessError(err).Error()}, nil
	}
	res, err := rpc.GraphClient.ListVersion(h.Context, &rpcgraph.ListVersionReq{GraphId: req.GraphId})
	if err != nil {
		return &graph.BaseResponsePageGraphVersionVO{Code: 1, Message: err.Error()}, err
	}
	records := make([]*graph.GraphVersionVO, 0, len(res.VersionList))
	for _, item := range res.VersionList {
		records = append(records, &graph.GraphVersionVO{
			Version:     item.Version,
			VersionName: item.VersionName,
			IsCurrent:   item.IsCurrent,
			CreateTime:  item.CreateTime,
			UpdateTime:  item.UpdateTime,
		})
	}
	return &graph.BaseResponsePageGraphVersionVO{
		Code:    0,
		Message: "success",
		Data:    makeGraphVersionPage(records, int64(len(records))),
	}, nil
}
