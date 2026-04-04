package service

import (
	"context"
	"os"
	"time"

	"github.com/MoScenix/industrial-fault-tree-ai/app/ai/promptutil"
	ai "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/ai"
)

type GetPromptService struct {
	ctx context.Context
}

func NewGetPromptService(ctx context.Context) *GetPromptService {
	return &GetPromptService{ctx: ctx}
}

func (s *GetPromptService) Run(req *ai.GetPromptReq) (resp *ai.GetPromptResp, err error) {
	content, err := promptutil.LoadPrompt(req.Mode)
	if err != nil {
		return nil, err
	}
	path := promptutil.PromptPath(req.Mode)
	updatedAt := ""
	if stat, statErr := os.Stat(path); statErr == nil {
		updatedAt = stat.ModTime().Format(time.DateTime)
	}
	return &ai.GetPromptResp{
		Success:   true,
		Mode:      req.Mode,
		Content:   content,
		UpdatedAt: updatedAt,
	}, nil
}
