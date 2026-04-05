package tools

import (
	"context"
	"fmt"
	"time"

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
	projectCtx, _ := ctx.Value(lutils.ProjectContextKey).(*lutils.ProjectContext)
	if projectCtx == nil || projectCtx.ProjectID == "" {
		fmt.Printf("[tool:write_tmp_graph] empty project context\n")
		return &WriteTmpGraphResponse{Success: false}, nil
	}
	if req == nil || req.Graph == nil {
		fmt.Printf("[tool:write_tmp_graph] project=%s version=%s empty graph request\n",
			projectCtx.ProjectID, projectCtx.CurrentVersion)
		return &WriteTmpGraphResponse{Success: false}, nil
	}
	fmt.Printf("[tool:write_tmp_graph] project=%s version=%s change_summary=%q nodes=%d\n",
		projectCtx.ProjectID, projectCtx.CurrentVersion, req.ChangeSummary, len(req.Graph.Nodes))

	req.Graph.Meta.BasedOnVersion = projectCtx.CurrentVersion
	req.Graph.Meta.GeneratedAt = time.Now().Format("2006-01-02 15:04:05")
	if req.Graph.Meta.Version == "" {
		req.Graph.Meta.Version = "tmp"
	}
	tmpPath := lutils.TmpVersionTreePath(projectCtx.ProjectID, projectCtx.CurrentVersion)
	if err := lutils.SaveGraphFile(tmpPath, req.Graph); err != nil {
		return nil, err
	}
	fmt.Printf("[tool:write_tmp_graph] project=%s version=%s wrote=%s\n",
		projectCtx.ProjectID, projectCtx.CurrentVersion, tmpPath)
	return &WriteTmpGraphResponse{
		Success:        true,
		TmpPath:        tmpPath,
		BasedOnVersion: projectCtx.CurrentVersion,
	}, nil
}

func NewWriteTmpGraphTool() (tool.InvokableTool, error) {
	return einotool.InferTool(
		"write_tmp_graph",
		"把完整tmp图写回工作版本。只有在你已经整理好整张图后才调用，不用于正式版本。",
		WriteTmpGraphFunc,
	)
}
