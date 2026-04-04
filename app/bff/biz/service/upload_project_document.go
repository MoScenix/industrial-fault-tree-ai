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

type UploadProjectDocumentService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUploadProjectDocumentService(Context context.Context, RequestContext *app.RequestContext) *UploadProjectDocumentService {
	return &UploadProjectDocumentService{RequestContext: RequestContext, Context: Context}
}

func (h *UploadProjectDocumentService) Run(req *document.UploadProjectDocumentRequest) (resp *document.BaseResponseBoolean, err error) {
	item, err := loadAuthorizedGraphRecord(h.Context, req.GraphId)
	if err != nil {
		return &document.BaseResponseBoolean{Code: 1, Message: graphAccessError(err).Error()}, nil
	}
	pdfID := newObjectID()
	globalPDFPath := filepath.Join("/document", pdfID+".pdf")
	projectPDFPath := filepath.Join(item.ProjectDir, "documents", pdfID+".pdf")
	if err := os.MkdirAll(filepath.Dir(globalPDFPath), 0o755); err != nil {
		return &document.BaseResponseBoolean{Code: 1, Message: err.Error()}, err
	}
	if err := os.MkdirAll(filepath.Dir(projectPDFPath), 0o755); err != nil {
		return &document.BaseResponseBoolean{Code: 1, Message: err.Error()}, err
	}
	if err := os.WriteFile(globalPDFPath, req.File, 0o644); err != nil {
		return &document.BaseResponseBoolean{Code: 1, Message: err.Error()}, err
	}
	if err := os.WriteFile(projectPDFPath, req.File, 0o644); err != nil {
		return &document.BaseResponseBoolean{Code: 1, Message: err.Error()}, err
	}
	_, err = rpc.DocumentClient.ParseProjectPDF(h.Context, &rpcdocument.ParseProjectPDFReq{
		ProjectId: strconv.FormatInt(req.GraphId, 10),
		PdfId:     pdfID, FileName: req.FileName, DisplayName: req.FileName,
	})
	if err != nil {
		return &document.BaseResponseBoolean{Code: 1, Message: err.Error()}, err
	}
	return &document.BaseResponseBoolean{Code: 0, Message: "success", Data: true}, nil
}
