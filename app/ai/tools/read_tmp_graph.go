package tools

import (
	"context"
	"encoding/json"
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
	Success          bool              `json:"success" jsonschema:"description=是否成功读取到图；true表示读图完成，false表示读图失败"`
	Graph            *lutils.GraphFile `json:"graph,omitempty" jsonschema:"description=读取到的完整图内容；success=true时通常有值"`
	RequestedVersion string            `json:"requested_version,omitempty" jsonschema:"description=工具请求中传入的版本；为空表示未显式指定版本"`
	ResolvedVersion  string            `json:"resolved_version,omitempty" jsonschema:"description=本次实际解析后读取的目标版本"`
	BasedOnVersion   string            `json:"based_on_version,omitempty" jsonschema:"description=图内容基于哪个正式版本生成；如果读取的是tmp图，这里通常是tmp图对应的正式版本"`
	Source           string            `json:"source,omitempty" jsonschema:"description=图的来源；tmp表示工作区临时图，version表示正式版本图，default表示未找到文件时返回的默认空图"`
	ErrorMessage     string            `json:"error_message,omitempty" jsonschema:"description=读取失败时的错误信息；success=false时优先查看这个字段"`
}

func logReadTmpGraphResult(resp *ReadTmpGraphResponse) {
	if resp == nil {
		fmt.Printf("[tool:read_tmp_graph] result <nil>\n")
		return
	}
	fmt.Printf("[tool:read_tmp_graph] result success=%v source=%s requested_version=%s resolved_version=%s based_on=%s error=%q\n",
		resp.Success, resp.Source, resp.RequestedVersion, resp.ResolvedVersion, resp.BasedOnVersion, resp.ErrorMessage)
	if resp.Graph == nil {
		fmt.Printf("[tool:read_tmp_graph] graph=<nil>\n")
		return
	}
	graphJSON, err := json.Marshal(resp.Graph)
	if err != nil {
		fmt.Printf("[tool:read_tmp_graph] graph_marshal_error=%q\n", err.Error())
		return
	}
	fmt.Printf("[tool:read_tmp_graph] graph=%s\n", string(graphJSON))
}

func ReadTmpGraphFunc(ctx context.Context, req *ReadTmpGraphRequest) (*ReadTmpGraphResponse, error) {
	projectCtx, _ := ctx.Value(lutils.ProjectContextKey).(*lutils.ProjectContext)
	if projectCtx == nil || projectCtx.ProjectID == "" {
		fmt.Printf("[tool:read_tmp_graph] empty project context\n")
		resp := &ReadTmpGraphResponse{
			Success:      false,
			Graph:        lutils.DefaultGraphFile(""),
			ErrorMessage: "empty project context",
		}
		logReadTmpGraphResult(resp)
		return resp, nil
	}
	requestVersion := ""
	if req != nil {
		requestVersion = req.Version
	}
	resolvedVersion := requestVersion
	if resolvedVersion == "" {
		resolvedVersion = projectCtx.CurrentVersion
	}
	fmt.Printf("[tool:read_tmp_graph] project=%s ctx_version=%s req_version=%s\n",
		projectCtx.ProjectID, projectCtx.CurrentVersion, requestVersion)

	if req != nil && req.Version != "" {
		graph, basedOnVersion, _, err := lutils.LoadWorkingGraph(projectCtx.ProjectID, req.Version)
		if err != nil {
			if os.IsNotExist(err) {
				resp := &ReadTmpGraphResponse{
					Success:          true,
					Graph:            lutils.DefaultGraphFile(projectCtx.ProjectID),
					RequestedVersion: requestVersion,
					ResolvedVersion:  resolvedVersion,
					BasedOnVersion:   req.Version,
					Source:           "default",
				}
				logReadTmpGraphResult(resp)
				return resp, nil
			}
			resp := &ReadTmpGraphResponse{
				Success:          false,
				RequestedVersion: requestVersion,
				ResolvedVersion:  resolvedVersion,
				ErrorMessage:     err.Error(),
			}
			logReadTmpGraphResult(resp)
			return resp, nil
		}
		fmt.Printf("[tool:read_tmp_graph] loaded project=%s version=%s based_on=%s\n",
			projectCtx.ProjectID, req.Version, basedOnVersion)
		resp := &ReadTmpGraphResponse{
			Success:          true,
			Graph:            graph,
			RequestedVersion: requestVersion,
			ResolvedVersion:  resolvedVersion,
			BasedOnVersion:   basedOnVersion,
			Source:           "tmp",
		}
		logReadTmpGraphResult(resp)
		return resp, nil
	}

	graph, basedOnVersion, fromTmp, err := lutils.LoadWorkingGraph(projectCtx.ProjectID, projectCtx.CurrentVersion)
	if err != nil {
		if os.IsNotExist(err) {
			resp := &ReadTmpGraphResponse{
				Success:          true,
				Graph:            lutils.DefaultGraphFile(projectCtx.ProjectID),
				RequestedVersion: requestVersion,
				ResolvedVersion:  resolvedVersion,
				BasedOnVersion:   projectCtx.CurrentVersion,
				Source:           "default",
			}
			logReadTmpGraphResult(resp)
			return resp, nil
		}
		resp := &ReadTmpGraphResponse{
			Success:          false,
			RequestedVersion: requestVersion,
			ResolvedVersion:  resolvedVersion,
			ErrorMessage:     err.Error(),
		}
		logReadTmpGraphResult(resp)
		return resp, nil
	}
	fmt.Printf("[tool:read_tmp_graph] loaded project=%s version=%s based_on=%s\n",
		projectCtx.ProjectID, projectCtx.CurrentVersion, basedOnVersion)
	source := "version"
	if fromTmp {
		source = "tmp"
	}
	resp := &ReadTmpGraphResponse{
		Success:          true,
		Graph:            graph,
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
		"读取当前tmp图。修改图前先读取，确认节点和邻接关系的现状。",
		ReadTmpGraphFunc,
	)
}
