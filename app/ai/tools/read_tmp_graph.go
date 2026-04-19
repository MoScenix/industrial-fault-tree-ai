package tools

import (
	"context"
	"fmt"
	"os"

	lutils "github.com/MoScenix/industrial-fault-tree-ai/app/ai/utils"
	"github.com/cloudwego/eino/components/tool"
	einotool "github.com/cloudwego/eino/components/tool/utils"
)

// ReadTmpGraphRequest is empty because the working version always comes from the request context.
type ReadTmpGraphRequest struct{}

// ReadTmpGraphResponse returns the numbered tmp graph file content the AI is allowed to operate on.
type ReadTmpGraphResponse struct {
	Success          bool   `json:"success" jsonschema:"description=是否成功读取到图文件；true表示读图完成，false表示读图失败"`
	FilePath         string `json:"file_path,omitempty" jsonschema:"description=当前工作图文件路径；写回时必须使用这个路径"`
	NumberedContent  string `json:"numbered_content,omitempty" jsonschema:"description=带行号的图文件内容，格式为行号加竖线和内容"`
	LineCount        int    `json:"line_count,omitempty" jsonschema:"description=当前文件总行数"`
	RequestedVersion string `json:"requested_version,omitempty" jsonschema:"description=工具请求中传入的版本；为空表示未显式指定版本"`
	ResolvedVersion  string `json:"resolved_version,omitempty" jsonschema:"description=本次实际解析后读取的目标版本"`
	BasedOnVersion   string `json:"based_on_version,omitempty" jsonschema:"description=图内容基于哪个正式版本生成；如果读取的是tmp图，这里通常是tmp图对应的正式版本"`
	Source           string `json:"source,omitempty" jsonschema:"description=图的来源；tmp表示工作区临时图，version表示正式版本图，default表示未找到文件时返回的默认空图"`
	ErrorMessage     string `json:"error_message,omitempty" jsonschema:"description=读取失败时的错误信息；success=false时优先查看这个字段"`
}

func logReadTmpGraphResult(resp *ReadTmpGraphResponse) {
	if resp == nil {
		fmt.Printf("[tool:read_tmp_graph] result <nil>\n")
		return
	}
	fmt.Printf("[tool:read_tmp_graph] result success=%v source=%s requested_version=%s resolved_version=%s based_on=%s error=%q\n",
		resp.Success, resp.Source, resp.RequestedVersion, resp.ResolvedVersion, resp.BasedOnVersion, resp.ErrorMessage)
	fmt.Printf("[tool:read_tmp_graph] file=%s line_count=%d\n", resp.FilePath, resp.LineCount)
}

func ReadTmpGraphFunc(ctx context.Context, req *ReadTmpGraphRequest) (*ReadTmpGraphResponse, error) {
	projectCtx, _ := ctx.Value(lutils.ProjectContextKey).(*lutils.ProjectContext)
	if projectCtx == nil || projectCtx.ProjectID == "" {
		fmt.Printf("[tool:read_tmp_graph] empty project context\n")
		resp := &ReadTmpGraphResponse{
			Success:      false,
			ErrorMessage: "empty project context",
		}
		logReadTmpGraphResult(resp)
		return resp, nil
	}
	requestVersion := projectCtx.CurrentVersion
	resolvedVersion := projectCtx.CurrentVersion
	fmt.Printf("[tool:read_tmp_graph] project=%s ctx_version=%s\n",
		projectCtx.ProjectID, projectCtx.CurrentVersion)
	filePath := lutils.TmpVersionTreePath(projectCtx.ProjectID, projectCtx.CurrentVersion)

	lines, basedOnVersion, source, err := loadRequestedVersionLines(projectCtx.ProjectID, projectCtx.CurrentVersion)
	if err != nil {
		resp := &ReadTmpGraphResponse{
			Success:          false,
			RequestedVersion: requestVersion,
			ResolvedVersion:  resolvedVersion,
			ErrorMessage:     err.Error(),
		}
		logReadTmpGraphResult(resp)
		return resp, nil
	}

	resp := &ReadTmpGraphResponse{
		Success:          true,
		FilePath:         filePath,
		NumberedContent:  lutils.NumberedLines(lines),
		LineCount:        len(lines),
		RequestedVersion: requestVersion,
		ResolvedVersion:  resolvedVersion,
		BasedOnVersion:   basedOnVersion,
		Source:           source,
	}
	logReadTmpGraphResult(resp)
	return resp, nil
}

func NewReadTmpGraphTool() (tool.InvokableTool, error) {
	return einotool.InferTool(
		"read_tmp_graph",
		"读取请求版本对应的 tmp 图文件，返回带行号的 JSON 内容和文件路径。版本始终来自请求上下文，不由模型指定。",
		ReadTmpGraphFunc,
	)
}

func loadRequestedVersionLines(projectID, version string) ([]string, string, string, error) {
	tmpPath := lutils.TmpVersionTreePath(projectID, version)
	if _, err := os.Stat(tmpPath); err == nil {
		lines, lineErr := lutils.LoadTextLines(tmpPath)
		return lines, version, "tmp", lineErr
	}
	graph, err := lutils.LoadGraphFile(lutils.VersionTreePath(projectID, version))
	if err == nil {
		lines, lineErr := lutils.GraphFileLines(graph)
		return lines, version, "version", lineErr
	}
	if os.IsNotExist(err) {
		lines, lineErr := lutils.DefaultGraphLines(projectID)
		return lines, version, "default", lineErr
	}
	return nil, "", "", err
}
