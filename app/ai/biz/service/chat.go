package service

import (
	"context"

	ai "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/ai"
)

type ChatService struct {
	ctx context.Context
}

// NewChatService new ChatService
func NewChatService(ctx context.Context) *ChatService {
	return &ChatService{ctx: ctx}
}

func (s *ChatService) Run(req *ai.ChatReq, stream ai.AiService_ChatServer) (err error) {
	// TODO: inject chat system prompt, then orchestrate AI tools:
	// get_project_context, rag_search, read_tmp_graph.
	_ = req
	_ = stream
	return
}
