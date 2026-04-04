package service

import (
	"context"
	"os"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/cloudwego/hertz/pkg/app"
)

type DiscardWorkingGraphService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDiscardWorkingGraphService(Context context.Context, RequestContext *app.RequestContext) *DiscardWorkingGraphService {
	return &DiscardWorkingGraphService{RequestContext: RequestContext, Context: Context}
}

func (h *DiscardWorkingGraphService) Run(req *graph.DiscardWorkingGraphRequest) (resp *graph.BaseResponseBoolean, err error) {
	item, err := loadAuthorizedGraphRecord(h.Context, req.GraphId)
	if err != nil {
		return &graph.BaseResponseBoolean{Code: 1, Message: graphAccessError(err).Error()}, nil
	}
	version := currentVersion(item.ProjectDir)
	if err := os.RemoveAll(tmpVersionDir(item.ProjectDir, version)); err != nil {
		return &graph.BaseResponseBoolean{Code: 1, Message: err.Error()}, err
	}
	return &graph.BaseResponseBoolean{Code: 0, Message: "success", Data: true}, nil
}
