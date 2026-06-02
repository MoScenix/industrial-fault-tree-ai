package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/document/biz/model"
	document "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/document"
)

type SearchDocumentsService struct {
	ctx context.Context
} // NewSearchDocumentsService new SearchDocumentsService
func NewSearchDocumentsService(ctx context.Context) *SearchDocumentsService {
	return &SearchDocumentsService{ctx: ctx}
}

// Run create note info
func (s *SearchDocumentsService) Run(req *document.SearchDocumentsReq) (resp *document.SearchDocumentsResp, err error) {
	query := model.NewDocumentQuery(s.ctx)
	results, err := query.SearchDocuments(req.GetUserId(), req.GetProjectId(), req.GetQuery(), req.GetTopK())
	if err != nil {
		return nil, err
	}

	items := make([]*document.SearchResult, 0, len(results))
	for _, result := range results {
		items = append(items, &document.SearchResult{
			DocumentId:   result.DocumentID,
			DocumentName: result.DocumentName,
			ChunkId:      result.ChunkID,
			Text:         result.Text,
			Score:        result.Score,
		})
	}

	return &document.SearchDocumentsResp{
		Results: items,
	}, nil
}
