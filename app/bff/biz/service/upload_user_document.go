package service

import (
	"context"
	"os"
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

func (h *UploadUserDocumentService) Run(req *document.UploadUserDocumentRequest) (resp *document.BaseResponseBoolean, err error) {
	if err := ensureLogin(h.Context); err != nil {
		return &document.BaseResponseBoolean{Code: 1, Message: graphAccessError(err).Error()}, nil
	}
	pdfID := newObjectID()
	pdfPath := filepath.Join("/document", pdfID+".pdf")
	if err := os.MkdirAll(filepath.Dir(pdfPath), 0o755); err != nil {
		return &document.BaseResponseBoolean{Code: 1, Message: err.Error()}, err
	}
	if err := os.WriteFile(pdfPath, req.File, 0o644); err != nil {
		return &document.BaseResponseBoolean{Code: 1, Message: err.Error()}, err
	}
	userIDInt, _ := getCurrentUserID(h.Context)
	userID := strconv.FormatInt(userIDInt, 10)
	_, err = rpc.DocumentClient.ParsePersonalPDF(h.Context, &rpcdocument.ParsePersonalPDFReq{
		UserId: userID, PdfId: pdfID, FileName: req.FileName, DisplayName: req.FileName,
	})
	if err != nil {
		return &document.BaseResponseBoolean{Code: 1, Message: err.Error()}, err
	}
	return &document.BaseResponseBoolean{Code: 0, Message: "success", Data: true}, nil
}
