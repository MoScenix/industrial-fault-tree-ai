package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/document/biz/model"
	document "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/document"
)

type ListDocumentsService struct {
	ctx context.Context
} // NewListDocumentsService new ListDocumentsService
func NewListDocumentsService(ctx context.Context) *ListDocumentsService {
	return &ListDocumentsService{ctx: ctx}
}

// Run create note info
func (s *ListDocumentsService) Run(req *document.ListDocumentsReq) (resp *document.ListDocumentsResp, err error) {
	query := model.NewDocumentQuery(s.ctx)
	ownerType := req.GetOwnerType().String()

	docs, err := query.ListDocuments(ownerType, req.GetOwnerId(), req.GetPage(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	total, err := query.CountDocuments(ownerType, req.GetOwnerId())
	if err != nil {
		return nil, err
	}

	items := make([]*document.Document, 0, len(docs))
	for _, docItem := range docs {
		items = append(items, toProtoDocument(docItem))
	}

	return &document.ListDocumentsResp{
		Documents: items,
		Total:     total,
	}, nil
}
