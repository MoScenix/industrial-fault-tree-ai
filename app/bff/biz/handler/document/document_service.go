package document

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/biz/service"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/biz/utils"
	document "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/document"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// UploadUserDocument .
// @router /document/user/upload [POST]
func UploadUserDocument(ctx context.Context, c *app.RequestContext) {
	var err error
	var req document.UploadUserDocumentRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &document.BaseResponseBoolean{}
	resp, err = service.NewUploadUserDocumentService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// UploadProjectDocument .
// @router /document/project/upload [POST]
func UploadProjectDocument(ctx context.Context, c *app.RequestContext) {
	var err error
	var req document.UploadProjectDocumentRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &document.BaseResponseBoolean{}
	resp, err = service.NewUploadProjectDocumentService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
