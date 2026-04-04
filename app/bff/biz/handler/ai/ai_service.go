package ai

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/biz/service"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/biz/utils"
	ai "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/ai"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// GetPrompt .
// @router /ai/prompt/get [GET]
func GetPrompt(ctx context.Context, c *app.RequestContext) {
	var err error
	var req ai.GetPromptRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &ai.BaseResponsePromptVO{}
	resp, err = service.NewGetPromptService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// UpdatePrompt .
// @router /ai/prompt/update [POST]
func UpdatePrompt(ctx context.Context, c *app.RequestContext) {
	var err error
	var req ai.UpdatePromptRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &ai.BaseResponseBoolean{}
	resp, err = service.NewUpdatePromptService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
