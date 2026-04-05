package service

import (
	"context"
	"fmt"
	"io"
	"os"

	lagent "github.com/MoScenix/industrial-fault-tree-ai/app/ai/agent"
	"github.com/MoScenix/industrial-fault-tree-ai/app/ai/promptutil"
	lutils "github.com/MoScenix/industrial-fault-tree-ai/app/ai/utils"
	ai "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/ai"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/kitex/pkg/klog"
)

type ChatService struct {
	ctx context.Context
}

// NewChatService new ChatService
func NewChatService(ctx context.Context) *ChatService {
	return &ChatService{ctx: ctx}
}

func (s *ChatService) Run(req *ai.ChatReq, stream ai.AiService_ChatServer) (err error) {
	currentVersion, _ := lutils.ReadCurrentVersion(req.ProjectId)
	targetVersion := req.Version
	if targetVersion == "" {
		targetVersion = currentVersion
	}
	_, statErr := os.Stat(lutils.TmpVersionTreePath(req.ProjectId, targetVersion))
	fmt.Printf("[ai:chat] project=%s version=%s tmp_ready=%v history=%d\n",
		req.ProjectId, targetVersion, statErr == nil, len(req.History))

	s.ctx = context.WithValue(s.ctx, lutils.ProjectContextKey, &lutils.ProjectContext{
		ProjectID:       req.ProjectId,
		CurrentVersion:  targetVersion,
		TmpVersionReady: statErr == nil,
	})

	prompt, err := promptutil.LoadPrompt(ai.PromptMode_MODIFY_MODE)
	if err != nil {
		return err
	}
	agent, err := lagent.NewReActAgent(s.ctx, ai.PromptMode_MODIFY_MODE)
	if err != nil {
		return err
	}

	messages := make([]*schema.Message, 0, len(req.History)+1)
	messages = append(messages, schema.SystemMessage(prompt))
	for _, history := range req.History {
		switch history.Role {
		case "assistant":
			messages = append(messages, schema.AssistantMessage(history.Content, nil))
		default:
			messages = append(messages, schema.UserMessage(history.Content))
		}
	}

	reader, err := agent.Stream(s.ctx, messages)
	if err != nil {
		return err
	}
	defer reader.Close()

	for {
		msg, recvErr := reader.Recv()
		if recvErr != nil {
			if recvErr == io.EOF {
				break
			}
			klog.Error(recvErr)
			return recvErr
		}
		if msg == nil || msg.Content == "" {
			continue
		}
		if err := stream.Send(&ai.ChatResp{Content: msg.Content}); err != nil {
			return err
		}
	}
	return nil
}
