package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MoScenix/industrial-fault-tree-ai/app/ai/infra/rpc"
	lutils "github.com/MoScenix/industrial-fault-tree-ai/app/ai/utils"
	document "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/document"
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
	ChunkID      string  `json:"chunk_id,omitempty" jsonschema:"description=命中的证据片段ID"`
	DocumentName string  `json:"document_name,omitempty" jsonschema:"description=证据来源文档名称"`
	Text         string  `json:"text,omitempty" jsonschema:"description=返回的证据片段文本内容"`
	Score        float64 `json:"score" jsonschema:"description=当前项目内简单匹配得分，分数越高代表命中越多"`
}

// RAGSearchResponse contains retrieved evidence chunks.
type RAGSearchResponse struct {
	Chunks        []*RAGChunk `json:"chunks" jsonschema:"description=检索返回的证据片段列表"`
	ReturnedCount int         `json:"returned_count" jsonschema:"description=本次实际返回的片段数量"`
	ErrorMessage  string      `json:"error_message,omitempty" jsonschema:"description=检索异常时的错误信息；为空表示没有工具层错误"`
}

func logRAGSearchResult(resp *RAGSearchResponse) {
	if resp == nil {
		fmt.Printf("[tool:rag_search] result <nil>\n")
		return
	}
	payload, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("[tool:rag_search] result_marshal_error=%q\n", err.Error())
		return
	}
	fmt.Printf("[tool:rag_search] result=%s\n", string(payload))
}

func RAGSearchFunc(ctx context.Context, req *RAGSearchRequest) (*RAGSearchResponse, error) {
	projectCtx, _ := ctx.Value(lutils.ProjectContextKey).(*lutils.ProjectContext)
	if projectCtx == nil || projectCtx.ProjectID == "" || req == nil || strings.TrimSpace(req.Query) == "" {
		fmt.Printf("[tool:rag_search] skip project_ctx_nil=%v project=%q query=%q\n",
			projectCtx == nil, func() string {
				if projectCtx == nil {
					return ""
				}
				return projectCtx.ProjectID
			}(), func() string {
				if req == nil {
					return ""
				}
				return req.Query
			}())
		resp := &RAGSearchResponse{Chunks: []*RAGChunk{}, ReturnedCount: 0}
		logRAGSearchResult(resp)
		return resp, nil
	}
	fmt.Printf("[tool:rag_search] project=%s version=%s query=%q topk=%d\n",
		projectCtx.ProjectID, projectCtx.CurrentVersion, req.Query, req.TopK)
	if rpc.DocumentClient == nil {
		resp := &RAGSearchResponse{
			Chunks:        []*RAGChunk{},
			ReturnedCount: 0,
			ErrorMessage:  "document rpc client is unavailable",
		}
		logRAGSearchResult(resp)
		return resp, nil
	}

	topK := req.TopK
	if topK <= 0 {
		topK = 3
	}

	searchResp, err := rpc.DocumentClient.SearchDocuments(ctx, &document.SearchDocumentsReq{
		UserId:    "",
		ProjectId: projectCtx.ProjectID,
		Query:     req.Query,
		TopK:      int64(topK),
	})
	if err != nil {
		resp := &RAGSearchResponse{
			Chunks:        []*RAGChunk{},
			ReturnedCount: 0,
			ErrorMessage:  err.Error(),
		}
		logRAGSearchResult(resp)
		return resp, nil
	}

	results := searchResp.GetResults()
	chunks := make([]*RAGChunk, 0, len(results))
	for _, item := range results {
		if item == nil {
			continue
		}
		chunks = append(chunks, &RAGChunk{
			ChunkID:      item.GetChunkId(),
			DocumentName: item.GetDocumentName(),
			Text:         item.GetText(),
			Score:        item.GetScore(),
		})
	}

	fmt.Printf("[tool:rag_search] project=%s matched=%d returned=%d\n",
		projectCtx.ProjectID, len(results), len(chunks))
	resp := &RAGSearchResponse{
		Chunks:        chunks,
		ReturnedCount: len(chunks),
	}
	logRAGSearchResult(resp)
	return resp, nil
}

func NewRAGSearchTool() (tool.InvokableTool, error) {
	return einotool.InferTool(
		"rag_search",
		"在当前项目文档中检索证据片段。需要依据资料回答、校验或修改图时使用。",
		RAGSearchFunc,
	)
}
