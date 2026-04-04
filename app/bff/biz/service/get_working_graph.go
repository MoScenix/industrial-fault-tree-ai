package service

import (
	"context"
	"os"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetWorkingGraphService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetWorkingGraphService(Context context.Context, RequestContext *app.RequestContext) *GetWorkingGraphService {
	return &GetWorkingGraphService{RequestContext: RequestContext, Context: Context}
}

func (h *GetWorkingGraphService) Run(req *graph.GetWorkingGraphRequest) (resp *graph.BaseResponseWorkingGraph, err error) {
	item, err := loadAuthorizedGraphRecord(h.Context, req.GraphId)
	if err != nil {
		return &graph.BaseResponseWorkingGraph{Code: 1, Message: graphAccessError(err).Error()}, nil
	}
	version := req.Version
	if version == "" {
		version = currentVersion(item.ProjectDir)
	}
	isTmp := false
	path := treePath(item.ProjectDir, version, false)
	if _, statErr := os.Stat(treePath(item.ProjectDir, version, true)); statErr == nil {
		isTmp = true
		path = treePath(item.ProjectDir, version, true)
	}
	content, _, err := readOptionalFile(path)
	if err != nil {
		return &graph.BaseResponseWorkingGraph{Code: 1, Message: err.Error()}, err
	}
	return &graph.BaseResponseWorkingGraph{
		Code:    0,
		Message: "success",
		Data: &graph.WorkingGraphVO{
			GraphId: req.GraphId, Version: version, IsTmp: isTmp, Content: normalizeJSONString(content),
		},
	}, nil
}
