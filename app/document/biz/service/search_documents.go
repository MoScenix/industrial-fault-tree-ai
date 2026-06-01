package service

import (
	"context"
	"encoding/json"
	"math"
	"os"
	"strings"
	"time"

	dalmilvus "github.com/MoScenix/industrial-fault-tree-ai/app/document/biz/dal/milvus"
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
	logRAGEval(s.ctx, req, items)

	return &document.SearchDocumentsResp{
		Results: items,
	}, nil
}

type ragEvalChunk struct {
	DocumentID string  `json:"document_id"`
	ChunkID    string  `json:"chunk_id"`
	Rank       int     `json:"rank"`
	Cosine     float64 `json:"cosine"`
	Weight     float64 `json:"weight"`
	TextLen    int     `json:"text_len"`
}

type ragEvalRecord struct {
	TimestampISO       string         `json:"timestamp_iso"`
	Query              string         `json:"query"`
	UserID             string         `json:"user_id,omitempty"`
	ProjectID          string         `json:"project_id,omitempty"`
	TopK               int64          `json:"top_k"`
	ReturnedCount      int            `json:"returned_count"`
	WeightedSimilarity float64        `json:"weighted_similarity"`
	EmbeddingError     string         `json:"embedding_error,omitempty"`
	Chunks             []ragEvalChunk `json:"chunks"`
}

func logRAGEval(ctx context.Context, req *document.SearchDocumentsReq, items []*document.SearchResult) {
	logPath := strings.TrimSpace(os.Getenv("RAG_EVAL_LOG_FILE"))
	if logPath == "" {
		return
	}

	record := ragEvalRecord{
		TimestampISO:  time.Now().UTC().Format(time.RFC3339),
		Query:         req.GetQuery(),
		UserID:        req.GetUserId(),
		ProjectID:     req.GetProjectId(),
		TopK:          req.GetTopK(),
		ReturnedCount: len(items),
		Chunks:        make([]ragEvalChunk, 0, len(items)),
	}

	if dalmilvus.Embedder != nil && strings.TrimSpace(req.GetQuery()) != "" && len(items) > 0 {
		texts := make([]string, 0, len(items)+1)
		texts = append(texts, req.GetQuery())
		for _, item := range items {
			texts = append(texts, item.GetText())
		}

		vecs, err := dalmilvus.Embedder.EmbedStrings(ctx, texts)
		if err != nil {
			record.EmbeddingError = err.Error()
		} else if len(vecs) == len(texts) {
			queryVec := vecs[0]
			var weightedSum float64
			var weightTotal float64
			for i, item := range items {
				cosine := cosineSimilarity(queryVec, vecs[i+1])
				weight := 1.0 / math.Log2(float64(i+2))
				weightedSum += cosine * weight
				weightTotal += weight
				record.Chunks = append(record.Chunks, ragEvalChunk{
					DocumentID: item.GetDocumentId(),
					ChunkID:    item.GetChunkId(),
					Rank:       i + 1,
					Cosine:     cosine,
					Weight:     weight,
					TextLen:    len([]rune(item.GetText())),
				})
			}
			if weightTotal > 0 {
				record.WeightedSimilarity = weightedSum / weightTotal
			}
		}
	}

	appendJSONL(logPath, record)
}

func cosineSimilarity(a, b []float64) float64 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	if n == 0 {
		return 0
	}

	var dot, na, nb float64
	for i := 0; i < n; i++ {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}

func appendJSONL(path string, record ragEvalRecord) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	_ = enc.Encode(record)
}
