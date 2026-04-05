package graph

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/biz/service"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/biz/utils"
	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// AddGraph .
// @router /graph/add [POST]
func AddGraph(ctx context.Context, c *app.RequestContext) {
	var err error
	var req graph.GraphAddRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &graph.BaseResponseLong{}
	resp, err = service.NewAddGraphService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// DeleteGraph .
// @router /graph/delete [POST]
func DeleteGraph(ctx context.Context, c *app.RequestContext) {
	var err error
	var req graph.DeleteRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &graph.BaseResponseBoolean{}
	resp, err = service.NewDeleteGraphService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// UpdateGraph .
// @router /graph/update [POST]
func UpdateGraph(ctx context.Context, c *app.RequestContext) {
	var err error
	var req graph.GraphUpdateRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &graph.BaseResponseBoolean{}
	resp, err = service.NewUpdateGraphService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// GetGraphVOById .
// @router /graph/get/vo [GET]
func GetGraphVOById(ctx context.Context, c *app.RequestContext) {
	var err error
	var req graph.GetGraphVOByIdRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &graph.BaseResponseGraphVO{}
	resp, err = service.NewGetGraphVOByIdService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// ListGraphVOByPage .
// @router /graph/list/page/vo [POST]
func ListGraphVOByPage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req graph.GraphQueryRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &graph.BaseResponsePageGraphVO{}
	resp, err = service.NewListGraphVOByPageService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// StartEdit .
// @router /graph/edit/start [POST]
func StartEdit(ctx context.Context, c *app.RequestContext) {
	var err error
	var req graph.StartEditRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &graph.BaseResponseGraphEditState{}
	resp, err = service.NewStartEditService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// GetWorkingGraph .
// @router /graph/edit/working [GET]
func GetWorkingGraph(ctx context.Context, c *app.RequestContext) {
	var err error
	var req graph.GetWorkingGraphRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &graph.BaseResponseWorkingGraph{}
	resp, err = service.NewGetWorkingGraphService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// DiscardWorkingGraph .
// @router /graph/edit/discard [POST]
func DiscardWorkingGraph(ctx context.Context, c *app.RequestContext) {
	var err error
	var req graph.DiscardWorkingGraphRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &graph.BaseResponseBoolean{}
	resp, err = service.NewDiscardWorkingGraphService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// SaveGraph .
// @router /graph/save [POST]
func SaveGraph(ctx context.Context, c *app.RequestContext) {
	var err error
	var req graph.SaveGraphRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &graph.BaseResponseSaveResult{}
	resp, err = service.NewSaveGraphService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// ListGraphVersion .
// @router /graph/version/list [GET]
func ListGraphVersion(ctx context.Context, c *app.RequestContext) {
	var err error
	var req graph.ListGraphVersionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &graph.BaseResponsePageGraphVersionVO{}
	resp, err = service.NewListGraphVersionService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// CreateGraphVersion .
// @router /graph/version/create [POST]
func CreateGraphVersion(ctx context.Context, c *app.RequestContext) {
	var err error
	var req graph.CreateGraphVersionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &graph.BaseResponseString{}
	resp, err = service.NewCreateGraphVersionService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// DeleteGraphVersion .
// @router /graph/version/delete [POST]
func DeleteGraphVersion(ctx context.Context, c *app.RequestContext) {
	var err error
	var req graph.DeleteGraphVersionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &graph.BaseResponseBoolean{}
	resp, err = service.NewDeleteGraphVersionService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// RenameGraphVersion .
// @router /graph/version/rename [POST]
func RenameGraphVersion(ctx context.Context, c *app.RequestContext) {
	var err error
	var req graph.RenameGraphVersionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &graph.BaseResponseBoolean{}
	resp, err = service.NewRenameGraphVersionService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// ChatToModifyGraph .
// @router /graph/chat [GET]
func ChatToModifyGraph(ctx context.Context, c *app.RequestContext) {
	var err error
	var req graph.ChatToModifyGraphRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &graph.ServerSentEventString{}
	resp, err = service.NewChatToModifyGraphService(ctx, c).Run(&req)
	if err != nil {
		return
	}
	_ = resp
}

// ListGraphMessage .
// @router /graph/message/list [GET]
func ListGraphMessage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req graph.ListGraphMessageRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &graph.BaseResponsePageGraphMessageVO{}
	resp, err = service.NewListGraphMessageService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// DownloadGraph .
// @router /graph/download/:graphId [GET]
func DownloadGraph(ctx context.Context, c *app.RequestContext) {
	var err error
	var req graph.DownloadGraphRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &graph.BaseResponseBytes{}
	resp, err = service.NewDownloadGraphService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// GetCurrentSuggestion .
// @router /graph/suggestion/current [GET]
func GetCurrentSuggestion(ctx context.Context, c *app.RequestContext) {
	var err error
	var req graph.GetCurrentSuggestionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewGetCurrentSuggestionService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// ValidateGraph .
// @router /graph/validate [POST]
func ValidateGraph(ctx context.Context, c *app.RequestContext) {
	var err error
	var req graph.ValidateGraphRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewValidateGraphService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
