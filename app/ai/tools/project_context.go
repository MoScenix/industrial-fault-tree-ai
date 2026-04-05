package tools

import (
	"context"
	"encoding/json"
	"fmt"

	lutils "github.com/MoScenix/industrial-fault-tree-ai/app/ai/utils"
	"github.com/cloudwego/eino/components/tool"
	einotool "github.com/cloudwego/eino/components/tool/utils"
)

// ProjectContextRequest is injected by the caller; AI does not provide project identity directly.
type ProjectContextRequest struct{}

// ProjectContextResponse describes the current AI execution context for a project.
type ProjectContextResponse struct {
	ProjectID       string `json:"project_id,omitempty" jsonschema:"description=当前对话绑定的项目ID"`
	DeviceName      string `json:"device_name,omitempty" jsonschema:"description=当前项目关联的设备名称"`
	TopEvent        string `json:"top_event,omitempty" jsonschema:"description=当前项目的顶事件描述"`
	CurrentVersion  string `json:"current_version,omitempty" jsonschema:"description=当前项目的正式当前版本号"`
	TmpVersionReady bool   `json:"tmp_version_ready" jsonschema:"description=当前版本是否已经存在可编辑的tmp图"`
	DocumentSummary string `json:"document_summary,omitempty" jsonschema:"description=当前项目文档摘要或补充上下文"`
}

func logProjectContextResult(resp *ProjectContextResponse) {
	if resp == nil {
		fmt.Printf("[tool:get_project_context] result <nil>\n")
		return
	}
	payload, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("[tool:get_project_context] result_marshal_error=%q\n", err.Error())
		return
	}
	fmt.Printf("[tool:get_project_context] result=%s\n", string(payload))
}

func GetProjectContextFunc(ctx context.Context, req *ProjectContextRequest) (*ProjectContextResponse, error) {
	_ = req
	projectCtx, _ := ctx.Value(lutils.ProjectContextKey).(*lutils.ProjectContext)
	if projectCtx == nil {
		fmt.Printf("[tool:get_project_context] empty project context\n")
		resp := &ProjectContextResponse{}
		logProjectContextResult(resp)
		return resp, nil
	}
	fmt.Printf("[tool:get_project_context] project=%s version=%s tmp_ready=%v\n",
		projectCtx.ProjectID, projectCtx.CurrentVersion, projectCtx.TmpVersionReady)
	resp := &ProjectContextResponse{
		ProjectID:       projectCtx.ProjectID,
		DeviceName:      projectCtx.DeviceName,
		TopEvent:        projectCtx.TopEvent,
		CurrentVersion:  projectCtx.CurrentVersion,
		TmpVersionReady: projectCtx.TmpVersionReady,
		DocumentSummary: projectCtx.DocumentSummary,
	}
	logProjectContextResult(resp)
	return resp, nil
}

func NewGetProjectContextTool() (tool.InvokableTool, error) {
	return einotool.InferTool(
		"get_project_context",
		"读取当前项目上下文。先用它确认设备、顶事件、当前版本。",
		GetProjectContextFunc,
	)
}
