package tools

import (
	"context"
	"fmt"
	"os"

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
	projectCtx, _ := ctx.Value(lutils.ProjectContextKey).(*lutils.ProjectContext)
	if projectCtx == nil || projectCtx.ProjectID == "" {
		fmt.Printf("[tool:read_tmp_graph] empty project context\n")
		return &ReadTmpGraphResponse{Graph: lutils.DefaultGraphFile("")}, nil
	}
	requestVersion := ""
	if req != nil {
		requestVersion = req.Version
	}
	fmt.Printf("[tool:read_tmp_graph] project=%s ctx_version=%s req_version=%s\n",
		projectCtx.ProjectID, projectCtx.CurrentVersion, requestVersion)

	if req != nil && req.Version != "" {
		graph, basedOnVersion, _, err := lutils.LoadWorkingGraph(projectCtx.ProjectID, req.Version)
		if err != nil {
			if os.IsNotExist(err) {
				return &ReadTmpGraphResponse{
					Graph:          lutils.DefaultGraphFile(projectCtx.ProjectID),
					BasedOnVersion: req.Version,
				}, nil
			}
			return nil, err
		}
		fmt.Printf("[tool:read_tmp_graph] loaded project=%s version=%s based_on=%s\n",
			projectCtx.ProjectID, req.Version, basedOnVersion)
		return &ReadTmpGraphResponse{Graph: graph, BasedOnVersion: basedOnVersion}, nil
	}

	graph, basedOnVersion, _, err := lutils.LoadWorkingGraph(projectCtx.ProjectID, projectCtx.CurrentVersion)
	if err != nil {
		if os.IsNotExist(err) {
			return &ReadTmpGraphResponse{
				Graph:          lutils.DefaultGraphFile(projectCtx.ProjectID),
				BasedOnVersion: projectCtx.CurrentVersion,
			}, nil
		}
		return nil, err
	}
	fmt.Printf("[tool:read_tmp_graph] loaded project=%s version=%s based_on=%s\n",
		projectCtx.ProjectID, projectCtx.CurrentVersion, basedOnVersion)
	return &ReadTmpGraphResponse{Graph: graph, BasedOnVersion: basedOnVersion}, nil
}

func NewReadTmpGraphTool() (tool.InvokableTool, error) {
	return einotool.InferTool(
		"read_tmp_graph",
		"读取当前tmp图。修改图前先读取，确认节点和邻接关系的现状。",
		ReadTmpGraphFunc,
	)
}
