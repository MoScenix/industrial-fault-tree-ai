package service

import (
	"context"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/cloudwego/hertz/pkg/app"
)

type DownloadGraphService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDownloadGraphService(Context context.Context, RequestContext *app.RequestContext) *DownloadGraphService {
	return &DownloadGraphService{RequestContext: RequestContext, Context: Context}
}

func (h *DownloadGraphService) Run(req *graph.DownloadGraphRequest) (resp *graph.BaseResponseBytes, err error) {
	item, err := loadAuthorizedGraphRecord(h.Context, req.GraphId)
	if err != nil {
		return &graph.BaseResponseBytes{Code: 1, Message: graphAccessError(err).Error()}, nil
	}
	version := req.Version
	if version == "" {
		version = currentVersion(item.ProjectDir)
	}
	content, _, err := readOptionalFile(treePath(item.ProjectDir, version, req.IsTmp))
	if err != nil {
		return &graph.BaseResponseBytes{Code: 1, Message: err.Error()}, err
	}
	return &graph.BaseResponseBytes{Code: 0, Message: "success", Data: []byte(content)}, nil
}
