package service

import (
	"context"

	aibff "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/ai"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcai "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/ai"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetPromptService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetPromptService(Context context.Context, RequestContext *app.RequestContext) *GetPromptService {
	return &GetPromptService{RequestContext: RequestContext, Context: Context}
}

func (h *GetPromptService) Run(req *aibff.GetPromptRequest) (resp *aibff.BaseResponsePromptVO, err error) {
	res, err := rpc.AiClient.GetPrompt(h.Context, &rpcai.GetPromptReq{
		Mode: rpcai.PromptMode(req.Mode),
	})
	if err != nil {
		return &aibff.BaseResponsePromptVO{Code: 1, Message: err.Error()}, nil
	}
	return &aibff.BaseResponsePromptVO{
		Code:    0,
		Message: "success",
		Data: &aibff.PromptVO{
			Mode:      int32(res.Mode),
			Content:   res.Content,
			UpdatedAt: res.UpdatedAt,
		},
	}, nil
}
