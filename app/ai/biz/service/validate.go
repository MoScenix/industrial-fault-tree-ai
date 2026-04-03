package service

import (
	"context"

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
	_ = req
	return &ai.ValidateResp{}, nil
}
