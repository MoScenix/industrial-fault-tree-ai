package service

import (
	"context"
	"fmt"

	"github.com/MoScenix/industrial-fault-tree-ai/app/document/biz/model"
	"github.com/MoScenix/industrial-fault-tree-ai/app/document/utils"
	document "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/document"
)

type ParseProjectPDFService struct {
	ctx context.Context
} // NewParseProjectPDFService new ParseProjectPDFService
func NewParseProjectPDFService(ctx context.Context) *ParseProjectPDFService {
	return &ParseProjectPDFService{ctx: ctx}
}

// Run create note info
func (s *ParseProjectPDFService) Run(req *document.ParseProjectPDFReq) (resp *document.ParsePDFResp, err error) {
	if req.GetProjectId() == "" || req.GetPdfId() == "" {
		return &document.ParsePDFResp{
			Success:      false,
			ErrorMessage: "project_id and pdf_id are required",
		}, nil
	}

	content, err := utils.ParsePDFFile(req.GetPdfId(), req.GetFileName())
	if err != nil {
		return &document.ParsePDFResp{
			Success:      false,
			DocumentId:   req.GetPdfId(),
			ErrorMessage: fmt.Sprintf("parse pdf failed: %v", err),
		}, nil
	}

	query := model.NewDocumentQuery(s.ctx)
	doc := model.Document{
		DocumentID:  req.GetPdfId(),
		OwnerType:   "PROJECT",
		OwnerID:     req.GetProjectId(),
		PdfID:       req.GetPdfId(),
		FileName:    req.GetFileName(),
		DisplayName: req.GetDisplayName(),
		ParseStatus: "SUCCESS",
		Summary:     content,
		Chunks:      utils.BuildChunks(req.GetPdfId(), content),
	}
	if doc.DisplayName == "" {
		doc.DisplayName = req.GetFileName()
	}

	if err = query.CreateDocument(doc); err != nil {
		return &document.ParsePDFResp{
			Success:      false,
			DocumentId:   req.GetPdfId(),
			ErrorMessage: err.Error(),
		}, nil
	}

	return &document.ParsePDFResp{
		Success:    true,
		DocumentId: req.GetPdfId(),
	}, nil
}
