package tools

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	lutils "github.com/MoScenix/industrial-fault-tree-ai/app/ai/utils"
	"github.com/cloudwego/eino/components/tool"
	einotool "github.com/cloudwego/eino/components/tool/utils"
)

// RAGSearchRequest describes a scoped retrieval request within the current project.
type RAGSearchRequest struct {
	Query   string            `json:"query" jsonschema:"description=你要检索的故障树相关问题或证据需求"`
	TopK    int32             `json:"top_k" jsonschema:"description=返回的证据片段数量"`
	Filters map[string]string
}

// RAGChunk is a retrieval result item.
type RAGChunk struct {
	ChunkID      string
	DocumentName string
	Text         string
	Score        float64
}

// RAGSearchResponse contains retrieved evidence chunks.
type RAGSearchResponse struct {
	Chunks []*RAGChunk
}

func RAGSearchFunc(ctx context.Context, req *RAGSearchRequest) (*RAGSearchResponse, error) {
	projectCtx, _ := ctx.Value(lutils.ProjectContextKey).(*lutils.ProjectContext)
	if projectCtx == nil || projectCtx.ProjectID == "" || req == nil || strings.TrimSpace(req.Query) == "" {
		return &RAGSearchResponse{Chunks: []*RAGChunk{}}, nil
	}

	docRoot := filepath.Join(lutils.ProjectDir(projectCtx.ProjectID), "documents")
	matches := make([]*RAGChunk, 0)
	queryTerms := tokenize(req.Query)

	_ = filepath.Walk(docRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}
		if !isSearchableDocument(path) {
			return nil
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		text := string(content)
		score := scoreText(text, queryTerms)
		if score <= 0 {
			return nil
		}
		matches = append(matches, &RAGChunk{
			ChunkID:      path,
			DocumentName: filepath.Base(path),
			Text:         truncateText(text, 600),
			Score:        score,
		})
		return nil
	})

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})

	limit := int(req.TopK)
	if limit <= 0 || limit > len(matches) {
		limit = len(matches)
	}
	return &RAGSearchResponse{Chunks: matches[:limit]}, nil
}

func NewRAGSearchTool() (tool.InvokableTool, error) {
	return einotool.InferTool(
		"rag_search",
		"在当前项目文档中检索证据片段。需要依据资料回答、校验或修改图时使用。",
		RAGSearchFunc,
	)
}

func isSearchableDocument(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".txt" || ext == ".md" || ext == ".json"
}

func tokenize(query string) []string {
	fields := strings.Fields(strings.ToLower(query))
	return fields
}

func scoreText(text string, terms []string) float64 {
	lower := strings.ToLower(text)
	var score float64
	for _, term := range terms {
		if term == "" {
			continue
		}
		score += float64(strings.Count(lower, term))
	}
	return score
}

func truncateText(text string, max int) string {
	runes := []rune(text)
	if len(runes) <= max {
		return text
	}
	return string(runes[:max])
}
