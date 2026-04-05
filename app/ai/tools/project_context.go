package tools

import (
	"context"
	"fmt"

	lutils "github.com/MoScenix/industrial-fault-tree-ai/app/ai/utils"
	"github.com/cloudwego/eino/components/tool"
	einotool "github.com/cloudwego/eino/components/tool/utils"
)

// ProjectContextRequest is injected by the caller; AI does not provide project identity directly.
type ProjectContextRequest struct{}

// ProjectContextResponse describes the current AI execution context for a project.
type ProjectContextResponse struct {
	ProjectID       string
	DeviceName      string
	TopEvent        string
	CurrentVersion  string
	TmpVersionReady bool
	DocumentSummary string
}

func GetProjectContextFunc(ctx context.Context, req *ProjectContextRequest) (*ProjectContextResponse, error) {
	_ = req
	projectCtx, _ := ctx.Value(lutils.ProjectContextKey).(*lutils.ProjectContext)
	if projectCtx == nil {
		fmt.Printf("[tool:get_project_context] empty project context\n")
		return &ProjectContextResponse{}, nil
	}
	fmt.Printf("[tool:get_project_context] project=%s version=%s tmp_ready=%v\n",
		projectCtx.ProjectID, projectCtx.CurrentVersion, projectCtx.TmpVersionReady)
	return &ProjectContextResponse{
		ProjectID:       projectCtx.ProjectID,
		DeviceName:      projectCtx.DeviceName,
		TopEvent:        projectCtx.TopEvent,
		CurrentVersion:  projectCtx.CurrentVersion,
		TmpVersionReady: projectCtx.TmpVersionReady,
		DocumentSummary: projectCtx.DocumentSummary,
	}, nil
}

func NewGetProjectContextTool() (tool.InvokableTool, error) {
	return einotool.InferTool(
		"get_project_context",
		"读取当前项目上下文。先用它确认设备、顶事件、当前版本。",
		GetProjectContextFunc,
	)
}
