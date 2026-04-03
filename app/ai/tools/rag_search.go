package tools

import (
	"context"

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
	_ = ctx
	_ = req
	return &RAGSearchResponse{}, nil
}

func NewRAGSearchTool() (tool.InvokableTool, error) {
	return einotool.InferTool(
		"rag_search",
		"在当前项目文档中检索证据片段。需要依据资料回答、校验或修改图时使用。",
		RAGSearchFunc,
	)
}
