package service

import (
	"context"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetCurrentSuggestionService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetCurrentSuggestionService(Context context.Context, RequestContext *app.RequestContext) *GetCurrentSuggestionService {
	return &GetCurrentSuggestionService{RequestContext: RequestContext, Context: Context}
}

func (h *GetCurrentSuggestionService) Run(req *graph.GetCurrentSuggestionRequest) (resp *graph.BaseResponseGraphSuggestion, err error) {
	item, err := loadAuthorizedGraphRecord(h.Context, req.GraphId)
	if err != nil {
		return &graph.BaseResponseGraphSuggestion{Code: 1, Message: graphAccessError(err).Error()}, nil
	}
	version := req.Version
	if version == "" {
		version = currentVersion(item.ProjectDir)
	}
	content, updateTime, err := readOptionalFile(suggestionPath(item.ProjectDir, version))
	if err != nil {
		return &graph.BaseResponseGraphSuggestion{Code: 1, Message: err.Error()}, err
	}
	return &graph.BaseResponseGraphSuggestion{
		Code:    0,
		Message: "success",
		Data: &graph.GraphSuggestionVO{
			GraphId: req.GraphId,
			Version: version,
			Content: content,
			UpdateTime: func() string {
				if updateTime.IsZero() {
					return ""
				}
				return updateTime.Format("2006-01-02 15:04:05")
			}(),
		},
	}, nil
}
