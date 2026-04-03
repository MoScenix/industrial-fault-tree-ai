package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/document/biz/model"
	document "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/document"
)

type GetDocumentService struct {
	ctx context.Context
} // NewGetDocumentService new GetDocumentService
func NewGetDocumentService(ctx context.Context) *GetDocumentService {
	return &GetDocumentService{ctx: ctx}
}

// Run create note info
func (s *GetDocumentService) Run(req *document.GetDocumentReq) (resp *document.GetDocumentResp, err error) {
	query := model.NewDocumentQuery(s.ctx)
	doc, err := query.GetDocumentByDocumentID(req.GetDocumentId())
	if err != nil {
		return nil, err
	}

	return &document.GetDocumentResp{
		Document: toProtoDocument(doc),
	}, nil
}
