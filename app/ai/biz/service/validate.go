package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	lagent "github.com/MoScenix/industrial-fault-tree-ai/app/ai/agent"
	"github.com/MoScenix/industrial-fault-tree-ai/app/ai/promptutil"
	lutils "github.com/MoScenix/industrial-fault-tree-ai/app/ai/utils"
	"github.com/cloudwego/eino/schema"
	ai "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/ai"
)

type ValidateService struct {
	ctx context.Context
}

// NewValidateService creates a validate service.
func NewValidateService(ctx context.Context) *ValidateService {
	return &ValidateService{ctx: ctx}
}

// Run validates a project graph version. Business logic is intentionally deferred.
func (s *ValidateService) Run(req *ai.ValidateReq) (resp *ai.ValidateResp, err error) {
	currentVersion, _ := lutils.ReadCurrentVersion(req.ProjectId)
	targetVersion := req.Version
	if targetVersion == "" {
		targetVersion = currentVersion
	}
	_, statErr := os.Stat(lutils.TmpVersionTreePath(req.ProjectId, targetVersion))
	fmt.Printf("[ai:validate] project=%s version=%s tmp_ready=%v\n",
		req.ProjectId, targetVersion, statErr == nil)
	s.ctx = context.WithValue(s.ctx, lutils.ProjectContextKey, &lutils.ProjectContext{
		ProjectID:       req.ProjectId,
		CurrentVersion:  targetVersion,
		TmpVersionReady: statErr == nil,
	})

	prompt, err := promptutil.LoadPrompt(ai.PromptMode_LOG_MODE)
	if err != nil {
		return nil, err
	}
	agent, err := lagent.NewReActAgent(s.ctx, ai.PromptMode_LOG_MODE)
	if err != nil {
		return nil, err
	}

	result, err := agent.Generate(s.ctx, []*schema.Message{
		schema.SystemMessage(prompt),
		schema.UserMessage(fmt.Sprintf("请为项目 %s 的版本 %s 生成一份校验建议，输出 Markdown，总结当前图的风险、缺失信息和改进建议。", req.ProjectId, targetVersion)),
	})
	if err != nil {
		return nil, err
	}

	suggestionPath := lutils.SuggestionPath(req.ProjectId, targetVersion)
	if err := os.MkdirAll(filepath.Dir(suggestionPath), 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(suggestionPath, []byte(result.Content), 0o644); err != nil {
		return nil, err
	}

	return &ai.ValidateResp{
		Valid:   true,
		Summary: result.Content,
		Issues:  []*ai.ValidationIssue{},
	}, nil
}
