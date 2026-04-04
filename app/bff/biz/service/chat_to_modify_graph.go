package service

import (
	"context"
	"strconv"
	"strings"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcai "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/ai"
	rpcgraph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
	"github.com/cloudwego/hertz/pkg/app"
)

type ChatToModifyGraphService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewChatToModifyGraphService(Context context.Context, RequestContext *app.RequestContext) *ChatToModifyGraphService {
	return &ChatToModifyGraphService{RequestContext: RequestContext, Context: Context}
}

func (h *ChatToModifyGraphService) Run(req *graph.ChatToModifyGraphRequest) (resp *graph.ServerSentEventString, err error) {
	messageResp, err := rpc.GraphClient.ListGraphMessage(h.Context, &rpcgraph.ListGraphMessageReq{
		GraphId:  req.GraphId,
		PageSize: 20,
	})
	if err != nil {
		return &graph.ServerSentEventString{Message: err.Error()}, err
	}
	history := make([]*rpcai.HistoryItem, 0, len(messageResp.MessageList)+1)
	for i := len(messageResp.MessageList) - 1; i >= 0; i-- {
		m := messageResp.MessageList[i]
		history = append(history, &rpcai.HistoryItem{Role: m.Role, Content: m.Content})
	}
	history = append(history, &rpcai.HistoryItem{Role: "user", Content: req.Message})
	stream, err := rpc.AiClient.Chat(h.Context, &rpcai.ChatReq{
		ProjectId: strconv.FormatInt(req.GraphId, 10),
		History:   history,
	})
	if err != nil {
		return &graph.ServerSentEventString{Message: err.Error()}, err
	}
	var builder strings.Builder
	for {
		msg, recvErr := stream.Recv()
		if recvErr != nil {
			break
		}
		builder.WriteString(msg.Content)
	}
	answer := builder.String()
	userID, _ := getCurrentUserID(h.Context)
	if answer != "" {
		_, _ = rpc.GraphClient.AddGraphMessage(h.Context, &rpcgraph.AddGraphMessageReq{
			GraphId: req.GraphId, UserId: userID, Role: "user", Content: req.Message,
		})
		_, _ = rpc.GraphClient.AddGraphMessage(h.Context, &rpcgraph.AddGraphMessageReq{
			GraphId: req.GraphId, UserId: 0, Role: "assistant", Content: answer,
		})
	}
	return &graph.ServerSentEventString{D: answer, Message: "success"}, nil
}
