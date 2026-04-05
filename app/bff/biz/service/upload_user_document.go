package service

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"strconv"

	document "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/document"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcdocument "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/document"
	"github.com/cloudwego/hertz/pkg/app"
)

type UploadUserDocumentService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUploadUserDocumentService(Context context.Context, RequestContext *app.RequestContext) *UploadUserDocumentService {
	return &UploadUserDocumentService{RequestContext: RequestContext, Context: Context}
}

func (h *UploadUserDocumentService) Run(fileHeader *multipart.FileHeader) (resp *document.BaseResponseBoolean, err error) {
	if err := ensureLogin(h.Context); err != nil {
		return &document.BaseResponseBoolean{Code: 1, Message: graphAccessError(err).Error()}, nil
	}
	if fileHeader == nil {
		return &document.BaseResponseBoolean{Code: 1, Message: "file is required"}, nil
	}
	pdfID := newObjectID()
	fileName := sanitizeUploadFileName(fileHeader, pdfID+".pdf")
	pdfPath, err := saveUploadedFileToDocDir(h.RequestContext, fileHeader, filepath.Join("/document", pdfID), fileName)
	if err != nil {
		return &document.BaseResponseBoolean{Code: 1, Message: err.Error()}, nil
	}
	_ = pdfPath
	userIDInt, _ := getCurrentUserID(h.Context)
	userID := strconv.FormatInt(userIDInt, 10)
	parseResp, err := rpc.DocumentClient.ParsePersonalPDF(h.Context, &rpcdocument.ParsePersonalPDFReq{
		UserId: userID, PdfId: pdfID, FileName: fileName, DisplayName: fileName,
	})
	if err != nil {
		return &document.BaseResponseBoolean{Code: 1, Message: err.Error()}, nil
	}
	if parseResp == nil {
		return &document.BaseResponseBoolean{Code: 1, Message: "document parse returned empty response"}, nil
	}
	if !parseResp.GetSuccess() {
		message := parseResp.GetErrorMessage()
		if message == "" {
			message = "document parse failed"
		}
		return &document.BaseResponseBoolean{Code: 1, Message: message}, nil
	}
	return &document.BaseResponseBoolean{Code: 0, Message: "success", Data: true}, nil
}
