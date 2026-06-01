package model

import (
	"context"
	"errors"
	"fmt"
	"strings"

	dalmilvus "github.com/MoScenix/industrial-fault-tree-ai/app/document/biz/dal/milvus"
	milvus2 "github.com/cloudwego/eino-ext/components/retriever/milvus2"
	eretriever "github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
)

type Document struct {
	DocumentID  string
	OwnerType   string
	OwnerID     string
	PdfID       string
	FileName    string
	DisplayName string
	ParseStatus string
	Summary     string
	Chunks      []DocumentChunk
	CreatedAt   string
	UpdatedAt   string
}

type DocumentChunk struct {
	ChunkID    string
	DocumentID string
	Text       string
	Page       int64
	Order      int64
}

type SearchResult struct {
	DocumentID   string
	DocumentName string
	ChunkID      string
	Text         string
	Score        float64
}

var (
	ErrMilvusIndexerUnavailable    = errors.New("milvus indexer is unavailable")
	ErrMilvusRetrieverUnavailable  = errors.New("milvus retriever is unavailable")
	ErrDocumentLookupUnimplemented = errors.New("document lookup is not implemented")
	ErrSearchScopeRequired         = errors.New("search scope required: project_id or user_id")
)

const embeddingBatchSize = 10

type DocumentQuery struct {
	ctx context.Context
}

func NewDocumentQuery(ctx context.Context) *DocumentQuery {
	return &DocumentQuery{
		ctx: ctx,
	}
}

func (q *DocumentQuery) CreateDocument(doc Document) error {
	if dalmilvus.Indexer == nil {
		return ErrMilvusIndexerUnavailable
	}

	docs := make([]*schema.Document, 0, len(doc.Chunks))
	for _, chunk := range doc.Chunks {
		docs = append(docs, &schema.Document{
			ID:      chunk.ChunkID,
			Content: chunk.Text,
			MetaData: map[string]any{
				"document_id":   doc.DocumentID,
				"document_name": doc.DisplayName,
				"owner_type":    doc.OwnerType,
				"owner_id":      doc.OwnerID,
				"pdf_id":        doc.PdfID,
				"file_name":     doc.FileName,
				"page":          chunk.Page,
				"order":         chunk.Order,
			},
		})
	}

	for start := 0; start < len(docs); start += embeddingBatchSize {
		end := start + embeddingBatchSize
		if end > len(docs) {
			end = len(docs)
		}

		if _, err := dalmilvus.Indexer.Store(q.ctx, docs[start:end]); err != nil {
			return err
		}
	}

	if err := dalmilvus.FlushCollection(q.ctx); err != nil {
		return err
	}

	return nil
}

func (q *DocumentQuery) GetDocumentByDocumentID(documentID string) (Document, error) {
	return Document{}, ErrDocumentLookupUnimplemented
}

func (q *DocumentQuery) ListDocuments(ownerType, ownerID string, page, pageSize int64) ([]Document, error) {
	return nil, ErrDocumentLookupUnimplemented
}

func (q *DocumentQuery) CountDocuments(ownerType, ownerID string) (int64, error) {
	return 0, ErrDocumentLookupUnimplemented
}

func (q *DocumentQuery) SearchDocuments(userID, projectID, query string, topK int64) ([]SearchResult, error) {
	if dalmilvus.Retriever == nil {
		return nil, ErrMilvusRetrieverUnavailable
	}
	filterExpr, err := buildOwnerFilter(userID, projectID)
	if err != nil {
		return nil, err
	}

	opts := make([]eretriever.Option, 0, 1)
	if topK > 0 {
		opts = append(opts, eretriever.WithTopK(int(topK)))
	}
	opts = append(opts, milvus2.WithFilter(filterExpr))

	docs, err := dalmilvus.Retriever.Retrieve(q.ctx, query, opts...)
	if err != nil {
		return nil, err
	}

	results := make([]SearchResult, 0, len(docs))
	for _, doc := range docs {
		if doc == nil {
			continue
		}

		result := SearchResult{
			DocumentID:   valueAsString(doc.MetaData["document_id"]),
			DocumentName: valueAsString(doc.MetaData["document_name"]),
			ChunkID:      doc.ID,
			Text:         doc.Content,
		}
		results = append(results, result)
	}

	return results, nil
}

func valueAsString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func buildOwnerFilter(userID, projectID string) (string, error) {
	projectID = strings.TrimSpace(projectID)
	userID = strings.TrimSpace(userID)

	if projectID != "" {
		return fmt.Sprintf(
			"owner_type == 'PROJECT' and owner_id == '%s'",
			escapeMilvusString(projectID),
		), nil
	}

	if userID != "" {
		return fmt.Sprintf(
			"owner_type == 'PERSONAL' and owner_id == '%s'",
			escapeMilvusString(userID),
		), nil
	}

	return "", ErrSearchScopeRequired
}

func escapeMilvusString(s string) string {
	escaped := strings.ReplaceAll(s, "\\", "\\\\")
	return strings.ReplaceAll(escaped, "'", "\\'")
}
