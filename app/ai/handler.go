package main

import (
	"context"
	ai "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/ai"
	"github.com/MoScenix/industrial-fault-tree-ai/app/ai/biz/service"
)

// AiServiceImpl implements the last service interface defined in the IDL.
type AiServiceImpl struct{}

func (s *AiServiceImpl) Chat(req *ai.ChatReq, stream ai.AiService_ChatServer) (err error) {
	ctx := context.Background()
	err = service.NewChatService(ctx).Run(req, stream)
	return
}

// Validate implements the AiServiceImpl interface.
func (s *AiServiceImpl) Validate(ctx context.Context, req *ai.ValidateReq) (resp *ai.ValidateResp, err error) {
	resp, err = service.NewValidateService(ctx).Run(req)

	return resp, err
}

// UpdatePrompt implements the AiServiceImpl interface.
func (s *AiServiceImpl) UpdatePrompt(ctx context.Context, req *ai.UpdatePromptReq) (resp *ai.UpdatePromptResp, err error) {
	resp, err = service.NewUpdatePromptService(ctx).Run(req)

	return resp, err
}
