package service

import (
	"context"
	"testing"
	ai "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/ai"
)

func TestGetPrompt_Run(t *testing.T) {
	ctx := context.Background()
	s := NewGetPromptService(ctx)
	// init req and assert value

	req := &ai.GetPromptReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
