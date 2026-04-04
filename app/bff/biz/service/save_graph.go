package service

import (
	"context"
	"os"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcgraph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
	"github.com/cloudwego/hertz/pkg/app"
)

type SaveGraphService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewSaveGraphService(Context context.Context, RequestContext *app.RequestContext) *SaveGraphService {
	return &SaveGraphService{RequestContext: RequestContext, Context: Context}
}

func (h *SaveGraphService) Run(req *graph.SaveGraphRequest) (resp *graph.BaseResponseSaveResult, err error) {
	item, err := loadAuthorizedGraphRecord(h.Context, req.GraphId)
	if err != nil {
		return &graph.BaseResponseSaveResult{Code: 1, Message: graphAccessError(err).Error()}, nil
	}
	if !req.UseTmp && req.Content != "" {
		if err := os.MkdirAll(tmpVersionDir(item.ProjectDir, req.FromVersion), 0o755); err != nil {
			return &graph.BaseResponseSaveResult{Code: 1, Message: err.Error()}, err
		}
		if err := os.WriteFile(treePath(item.ProjectDir, req.FromVersion, true), []byte(normalizeJSONString(req.Content)), 0o644); err != nil {
			return &graph.BaseResponseSaveResult{Code: 1, Message: err.Error()}, err
		}
	}
	res, err := rpc.GraphClient.Save(h.Context, &rpcgraph.SaveReq{
		GraphId: req.GraphId, FromVersion: req.FromVersion, ToVersion: req.ToVersion, Remark: req.Remark,
	})
	if err != nil {
		return &graph.BaseResponseSaveResult{Code: 1, Message: err.Error()}, err
	}
	return &graph.BaseResponseSaveResult{
		Code:    0,
		Message: "success",
		Data: &graph.SaveResult{
			FromVersion: res.FromVersion,
			ToVersion:   res.ToVersion,
			Message:     res.Message,
		},
	}, nil
}
