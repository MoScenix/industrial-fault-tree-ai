package service

import (
	"context"

	ai "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/ai"
)

type UpdatePromptService struct {
	ctx context.Context
}

// NewUpdatePromptService creates a prompt update service.
func NewUpdatePromptService(ctx context.Context) *UpdatePromptService {
	return &UpdatePromptService{ctx: ctx}
}

// Run updates one of the managed prompts. Business logic is intentionally deferred.
func (s *UpdatePromptService) Run(req *ai.UpdatePromptReq) (resp *ai.UpdatePromptResp, err error) {
	_ = req
	return &ai.UpdatePromptResp{}, nil
}
