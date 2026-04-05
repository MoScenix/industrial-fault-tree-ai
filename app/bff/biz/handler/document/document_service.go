package document

import (
	"context"
	"errors"
	"strconv"

	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/biz/service"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/biz/utils"
	document "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/document"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// UploadUserDocument .
// @router /document/user/upload [POST]
func UploadUserDocument(ctx context.Context, c *app.RequestContext) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &document.BaseResponseBoolean{}
	resp, err = service.NewUploadUserDocumentService(ctx, c).Run(fileHeader)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// UploadProjectDocument .
// @router /document/project/upload [POST]
func UploadProjectDocument(ctx context.Context, c *app.RequestContext) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	graphID, err := strconv.ParseInt(c.PostForm("graphId"), 10, 64)
	if err != nil || graphID <= 0 {
		if err == nil {
			err = errors.New("graphId is required")
		}
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &document.BaseResponseBoolean{}
	resp, err = service.NewUploadProjectDocumentService(ctx, c).Run(graphID, fileHeader)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
