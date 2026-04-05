package service

import (
	"context"
	"fmt"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcai "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/ai"
	"github.com/cloudwego/hertz/pkg/app"
)

type ValidateGraphService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewValidateGraphService(Context context.Context, RequestContext *app.RequestContext) *ValidateGraphService {
	return &ValidateGraphService{RequestContext: RequestContext, Context: Context}
}

func (h *ValidateGraphService) Run(req *graph.ValidateGraphRequest) (resp *graph.BaseResponseBoolean, err error) {
	item, err := loadAuthorizedGraphRecord(h.Context, req.GraphId)
	if err != nil {
		return &graph.BaseResponseBoolean{Code: 1, Message: graphAccessError(err).Error()}, nil
	}
	version := req.Version
	if version == "" {
		version = currentVersion(item.ProjectDir)
	}
	if version == "" {
		return &graph.BaseResponseBoolean{Code: 1, Message: "当前项目没有可用版本"}, nil
	}
	fmt.Printf("[bff:validate] graph_id=%d project=%s version=%s\n",
		req.GraphId, projectIDFromDir(item.ProjectDir), version)
	_, err = rpc.AiClient.Validate(h.Context, &rpcai.ValidateReq{
		ProjectId: projectIDFromDir(item.ProjectDir),
		Version:   version,
	})
	if err != nil {
		return &graph.BaseResponseBoolean{Code: 1, Message: err.Error()}, err
	}
	return &graph.BaseResponseBoolean{Code: 0, Data: true, Message: "success"}, nil
}
