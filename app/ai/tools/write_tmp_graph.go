package tools

import (
	"context"

	lutils "github.com/MoScenix/industrial-fault-tree-ai/app/ai/utils"
	"github.com/cloudwego/eino/components/tool"
	einotool "github.com/cloudwego/eino/components/tool/utils"
)

// WriteTmpGraphRequest writes a fully prepared tmp graph snapshot.
type WriteTmpGraphRequest struct {
	Graph         *lutils.GraphFile `json:"graph" jsonschema:"description=你要写回的完整tmp图"`
	ChangeSummary string            `json:"change_summary" jsonschema:"description=本次修改摘要，简短说明改了什么"`
}

// WriteTmpGraphResponse reports the tmp write outcome.
type WriteTmpGraphResponse struct {
	Success        bool
	TmpPath        string
	BasedOnVersion string
}

func WriteTmpGraphFunc(ctx context.Context, req *WriteTmpGraphRequest) (*WriteTmpGraphResponse, error) {
	_ = ctx
	_ = req
	return &WriteTmpGraphResponse{}, nil
}

func NewWriteTmpGraphTool() (tool.InvokableTool, error) {
	return einotool.InferTool(
		"write_tmp_graph",
		"把完整tmp图写回工作版本。只有在你已经整理好整张图后才调用，不能用于正式版本。",
		WriteTmpGraphFunc,
	)
}
