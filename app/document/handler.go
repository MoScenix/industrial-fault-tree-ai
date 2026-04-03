package main

import (
	"context"
	document "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/document"
	"github.com/MoScenix/industrial-fault-tree-ai/app/document/biz/service"
)

// DocumentServiceImpl implements the last service interface defined in the IDL.
type DocumentServiceImpl struct{}

// ParsePersonalPDF implements the DocumentServiceImpl interface.
func (s *DocumentServiceImpl) ParsePersonalPDF(ctx context.Context, req *document.ParsePersonalPDFReq) (resp *document.ParsePDFResp, err error) {
	resp, err = service.NewParsePersonalPDFService(ctx).Run(req)

	return resp, err
}

// ParseProjectPDF implements the DocumentServiceImpl interface.
func (s *DocumentServiceImpl) ParseProjectPDF(ctx context.Context, req *document.ParseProjectPDFReq) (resp *document.ParsePDFResp, err error) {
	resp, err = service.NewParseProjectPDFService(ctx).Run(req)

	return resp, err
}

// GetDocument implements the DocumentServiceImpl interface.
func (s *DocumentServiceImpl) GetDocument(ctx context.Context, req *document.GetDocumentReq) (resp *document.GetDocumentResp, err error) {
	resp, err = service.NewGetDocumentService(ctx).Run(req)

	return resp, err
}

// ListDocuments implements the DocumentServiceImpl interface.
func (s *DocumentServiceImpl) ListDocuments(ctx context.Context, req *document.ListDocumentsReq) (resp *document.ListDocumentsResp, err error) {
	resp, err = service.NewListDocumentsService(ctx).Run(req)

	return resp, err
}

// SearchDocuments implements the DocumentServiceImpl interface.
func (s *DocumentServiceImpl) SearchDocuments(ctx context.Context, req *document.SearchDocumentsReq) (resp *document.SearchDocumentsResp, err error) {
	resp, err = service.NewSearchDocumentsService(ctx).Run(req)

	return resp, err
}
