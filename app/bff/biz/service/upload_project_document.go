package service

import (
	"context"
	"mime/multipart"
	"path/filepath"

	document "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/document"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcdocument "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/document"
	"github.com/cloudwego/hertz/pkg/app"
)

type UploadProjectDocumentService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUploadProjectDocumentService(Context context.Context, RequestContext *app.RequestContext) *UploadProjectDocumentService {
	return &UploadProjectDocumentService{RequestContext: RequestContext, Context: Context}
}

func (h *UploadProjectDocumentService) Run(graphID int64, fileHeader *multipart.FileHeader) (resp *document.BaseResponseBoolean, err error) {
	item, err := loadAuthorizedGraphRecord(h.Context, graphID)
	if err != nil {
		return &document.BaseResponseBoolean{Code: 1, Message: graphAccessError(err).Error()}, nil
	}
	if fileHeader == nil {
		return &document.BaseResponseBoolean{Code: 1, Message: "file is required"}, nil
	}
	projectID := projectIDFromDir(item.ProjectDir)
	if projectID == "" {
		return &document.BaseResponseBoolean{Code: 1, Message: "invalid project scope"}, nil
	}
	pdfID := newObjectID()
	fileName := sanitizeUploadFileName(fileHeader, pdfID+".pdf")
	globalDocDir := filepath.Join("/document", pdfID)
	if _, err := saveUploadedFileToDocDir(h.RequestContext, fileHeader, globalDocDir, fileName); err != nil {
		return &document.BaseResponseBoolean{Code: 1, Message: err.Error()}, nil
	}

	parseResp, err := rpc.DocumentClient.ParseProjectPDF(h.Context, &rpcdocument.ParseProjectPDFReq{
		ProjectId: projectID,
		PdfId:     pdfID, FileName: fileName, DisplayName: fileName,
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
