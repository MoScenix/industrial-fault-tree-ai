package agent

import (
	"context"
	"os"

	"github.com/MoScenix/industrial-fault-tree-ai/app/ai/tools"
	ai "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/ai"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
)

func of[T any](v T) *T {
	return &v
}

func NewChatModel(ctx context.Context) (*qwen.ChatModel, error) {
	modelName := os.Getenv("MODEL_NAME")
	if modelName == "" {
		modelName = "qwen-max"
	}
	baseURL := os.Getenv("MODEL_BASE_URL")
	if baseURL == "" {
		baseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	}

	return qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		BaseURL:     baseURL,
		APIKey:      os.Getenv("DASHSCOPE_API_KEY"),
		Model:       modelName,
		MaxTokens:   of(2048),
		Temperature: of(float32(0.7)),
		TopP:        of(float32(0.7)),
	})
}

func NewReActAgent(ctx context.Context, mode ai.PromptMode) (*react.Agent, error) {
	cm, err := NewChatModel(ctx)
	if err != nil {
		return nil, err
	}

	baseTools, err := toolsByMode(mode)
	if err != nil {
		return nil, err
	}

	return react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: cm,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: baseTools,
		},
		MaxStep: 20,
	})
}

func NewAgent(ctx context.Context, mode ai.PromptMode) (*react.Agent, error) {
	return NewReActAgent(ctx, mode)
}

func toolsByMode(mode ai.PromptMode) ([]tool.BaseTool, error) {
	invokableTools := make([]tool.InvokableTool, 0, 4)

	getProjectContextTool, err := tools.NewGetProjectContextTool()
	if err != nil {
		return nil, err
	}
	invokableTools = append(invokableTools, getProjectContextTool)

	ragSearchTool, err := tools.NewRAGSearchTool()
	if err != nil {
		return nil, err
	}
	invokableTools = append(invokableTools, ragSearchTool)

	readTmpGraphTool, err := tools.NewReadTmpGraphTool()
	if err != nil {
		return nil, err
	}
	invokableTools = append(invokableTools, readTmpGraphTool)

	if mode == ai.PromptMode_MODIFY_MODE {
		writeTmpGraphTool, err := tools.NewWriteTmpGraphTool()
		if err != nil {
			return nil, err
		}
		invokableTools = append(invokableTools, writeTmpGraphTool)
	}

	baseTools := make([]tool.BaseTool, 0, len(invokableTools))
	for _, t := range invokableTools {
		baseTools = append(baseTools, t)
	}
	return baseTools, nil
}
