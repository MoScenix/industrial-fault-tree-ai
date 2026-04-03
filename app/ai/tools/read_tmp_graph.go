package tools

import (
	"context"

	lutils "github.com/MoScenix/industrial-fault-tree-ai/app/ai/utils"
	"github.com/cloudwego/eino/components/tool"
	einotool "github.com/cloudwego/eino/components/tool/utils"
)

// ReadTmpGraphRequest loads the current tmp graph or a specified version snapshot.
type ReadTmpGraphRequest struct {
	Version string `json:"version,omitempty" jsonschema:"description=可选，指定要读取的版本；为空时读取当前tmp图"`
}

// ReadTmpGraphResponse returns the graph the AI is allowed to operate on.
type ReadTmpGraphResponse struct {
	Graph          *lutils.GraphFile
	BasedOnVersion string
}

func ReadTmpGraphFunc(ctx context.Context, req *ReadTmpGraphRequest) (*ReadTmpGraphResponse, error) {
	_ = ctx
	_ = req
	return &ReadTmpGraphResponse{}, nil
}

func NewReadTmpGraphTool() (tool.InvokableTool, error) {
	return einotool.InferTool(
		"read_tmp_graph",
		"读取当前tmp图。修改图前先读取，确认节点和邻接关系的现状。",
		ReadTmpGraphFunc,
	)
}
