package service

import (
	"context"
	"fmt"

	lutils "github.com/MoScenix/industrial-fault-tree-ai/app/bff/biz/utils"
	aibff "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/ai"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcai "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/ai"
	"github.com/cloudwego/hertz/pkg/app"
)

type UpdatePromptService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdatePromptService(Context context.Context, RequestContext *app.RequestContext) *UpdatePromptService {
	return &UpdatePromptService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdatePromptService) Run(req *aibff.UpdatePromptRequest) (resp *aibff.BaseResponseBoolean, err error) {
	operator := "system"
	if v := h.Context.Value(lutils.UserIdKey); v != nil {
		operator = fmt.Sprintf("%v", v)
	}
	_, err = rpc.AiClient.UpdatePrompt(h.Context, &rpcai.UpdatePromptReq{
		Mode:     rpcai.PromptMode(req.Mode),
		Content:  req.Content,
		Operator: operator,
	})
	if err != nil {
		return &aibff.BaseResponseBoolean{Code: 1, Message: err.Error()}, nil
	}
	return &aibff.BaseResponseBoolean{Code: 0, Message: "success", Data: true}, nil
}
