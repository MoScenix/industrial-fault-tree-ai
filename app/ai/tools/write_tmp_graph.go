package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"

	lutils "github.com/MoScenix/industrial-fault-tree-ai/app/ai/utils"
	"github.com/cloudwego/eino/components/tool"
	einotool "github.com/cloudwego/eino/components/tool/utils"
)

type WriteTmpGraphRequest struct {
	FilePath      string `json:"file_path" jsonschema:"description=要修改的图文件路径，必须来自 read_tmp_graph 返回值"`
	Operation     string `json:"operation" jsonschema:"description=操作类型，只允许 insert 或 delete"`
	Line          int    `json:"line,omitempty" jsonschema:"description=insert 时必填；表示在该行之前插入内容，允许范围 1..当前总行数+1"`
	StartLine     int    `json:"start_line,omitempty" jsonschema:"description=delete 时必填；表示删除区间起始行"`
	EndLine       int    `json:"end_line,omitempty" jsonschema:"description=delete 时必填；表示删除区间结束行，含本行"`
	Content       string `json:"content,omitempty" jsonschema:"description=insert 时必填；要插入的原始文本"`
	ChangeSummary string `json:"change_summary" jsonschema:"description=本次修改摘要，简短说明改了什么"`
}

type WriteTmpGraphResponse struct {
	Success        bool   `json:"success" jsonschema:"description=是否成功写回tmp图"`
	TmpPath        string `json:"tmp_path,omitempty" jsonschema:"description=tmp图实际写入的文件路径"`
	Operation      string `json:"operation,omitempty" jsonschema:"description=本次执行的操作类型"`
	LineCount      int    `json:"line_count,omitempty" jsonschema:"description=写回后的文件总行数"`
	BasedOnVersion string `json:"based_on_version,omitempty" jsonschema:"description=本次tmp图所基于的正式版本号"`
	ErrorMessage   string `json:"error_message,omitempty" jsonschema:"description=写图失败时的错误信息；success=false时优先查看这个字段"`
}

func logWriteTmpGraphResult(resp *WriteTmpGraphResponse) {
	if resp == nil {
		fmt.Printf("[tool:write_tmp_graph] result <nil>\n")
		return
	}
	payload, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("[tool:write_tmp_graph] result_marshal_error=%q\n", err.Error())
		return
	}
	fmt.Printf("[tool:write_tmp_graph] result=%s\n", string(payload))
}

func WriteTmpGraphFunc(ctx context.Context, req *WriteTmpGraphRequest) (*WriteTmpGraphResponse, error) {
	projectCtx, _ := ctx.Value(lutils.ProjectContextKey).(*lutils.ProjectContext)
	if projectCtx == nil || projectCtx.ProjectID == "" {
		fmt.Printf("[tool:write_tmp_graph] empty project context\n")
		resp := &WriteTmpGraphResponse{Success: false, ErrorMessage: "empty project context"}
		logWriteTmpGraphResult(resp)
		return resp, nil
	}
	if req == nil {
		fmt.Printf("[tool:write_tmp_graph] project=%s version=%s empty request\n",
			projectCtx.ProjectID, projectCtx.CurrentVersion)
		resp := &WriteTmpGraphResponse{Success: false, ErrorMessage: "empty request"}
		logWriteTmpGraphResult(resp)
		return resp, nil
	}

	tmpPath := lutils.TmpVersionTreePath(projectCtx.ProjectID, projectCtx.CurrentVersion)
	if filepath.Clean(req.FilePath) != filepath.Clean(tmpPath) {
		resp := &WriteTmpGraphResponse{
			Success:        false,
			TmpPath:        tmpPath,
			Operation:      req.Operation,
			BasedOnVersion: projectCtx.CurrentVersion,
			ErrorMessage:   "file_path must match read_tmp_graph returned tmp path",
		}
		logWriteTmpGraphResult(resp)
		return resp, nil
	}

	fmt.Printf("[tool:write_tmp_graph] project=%s version=%s op=%s line=%d start=%d end=%d content_len=%d change_summary=%q\n",
		projectCtx.ProjectID, projectCtx.CurrentVersion, req.Operation, req.Line, req.StartLine, req.EndLine, utf8.RuneCountInString(req.Content), req.ChangeSummary)

	lines, err := lutils.LoadTextLines(tmpPath)
	if err != nil {
		if !os.IsNotExist(err) {
			resp := &WriteTmpGraphResponse{
				Success:        false,
				TmpPath:        tmpPath,
				Operation:      req.Operation,
				BasedOnVersion: projectCtx.CurrentVersion,
				ErrorMessage:   err.Error(),
			}
			logWriteTmpGraphResult(resp)
			return resp, nil
		}
		lines, err = lutils.DefaultGraphLines(projectCtx.ProjectID)
		if err != nil {
			resp := &WriteTmpGraphResponse{
				Success:        false,
				TmpPath:        tmpPath,
				Operation:      req.Operation,
				BasedOnVersion: projectCtx.CurrentVersion,
				ErrorMessage:   err.Error(),
			}
			logWriteTmpGraphResult(resp)
			return resp, nil
		}
	}

	var nextLines []string
	switch req.Operation {
	case "insert":
		nextLines, err = lutils.InsertTextAtLine(lines, req.Line, req.Content)
	case "delete":
		nextLines, err = lutils.DeleteLineRange(lines, req.StartLine, req.EndLine)
	default:
		err = fmt.Errorf("unsupported operation: %s", req.Operation)
	}
	if err != nil {
		resp := &WriteTmpGraphResponse{
			Success:        false,
			TmpPath:        tmpPath,
			Operation:      req.Operation,
			BasedOnVersion: projectCtx.CurrentVersion,
			ErrorMessage:   err.Error(),
		}
		logWriteTmpGraphResult(resp)
		return resp, nil
	}

	if err := lutils.SaveTextLines(tmpPath, nextLines); err != nil {
		resp := &WriteTmpGraphResponse{
			Success:        false,
			TmpPath:        tmpPath,
			Operation:      req.Operation,
			BasedOnVersion: projectCtx.CurrentVersion,
			ErrorMessage:   err.Error(),
		}
		logWriteTmpGraphResult(resp)
		return resp, nil
	}

	fmt.Printf("[tool:write_tmp_graph] project=%s version=%s wrote=%s\n",
		projectCtx.ProjectID, projectCtx.CurrentVersion, tmpPath)
	resp := &WriteTmpGraphResponse{
		Success:        true,
		TmpPath:        tmpPath,
		Operation:      req.Operation,
		LineCount:      len(nextLines),
		BasedOnVersion: projectCtx.CurrentVersion,
	}
	logWriteTmpGraphResult(resp)
	return resp, nil
}

func NewWriteTmpGraphTool() (tool.InvokableTool, error) {
	return einotool.InferTool(
		"write_tmp_graph",
		"按行编辑 tmp 图文件。支持 insert 和 delete 两种操作。",
		WriteTmpGraphFunc,
	)
}
