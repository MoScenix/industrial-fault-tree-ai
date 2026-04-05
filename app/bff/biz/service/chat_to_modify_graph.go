package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcai "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/ai"
	rpcgraph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/sse"
)

type ChatToModifyGraphService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewChatToModifyGraphService(Context context.Context, RequestContext *app.RequestContext) *ChatToModifyGraphService {
	return &ChatToModifyGraphService{RequestContext: RequestContext, Context: Context}
}

func (h *ChatToModifyGraphService) Run(req *graph.ChatToModifyGraphRequest) (resp *graph.ServerSentEventString, err error) {
	w := sse.NewWriter(h.RequestContext)
	defer w.Close()

	item, err := loadAuthorizedGraphRecord(h.Context, req.GraphId)
	if err != nil {
		return sendChatErr(w, graphAccessError(err))
	}
	version := req.Version
	if version == "" {
		version = currentVersion(item.ProjectDir)
	}
	fmt.Printf("[bff:chat] graph_id=%d project=%s version=%s\n",
		req.GraphId, projectIDFromDir(item.ProjectDir), version)
	messageResp, err := rpc.GraphClient.ListGraphMessage(h.Context, &rpcgraph.ListGraphMessageReq{
		GraphId:  req.GraphId,
		PageSize: 20,
	})
	if err != nil {
		return sendChatErr(w, err)
	}
	history := make([]*rpcai.HistoryItem, 0, len(messageResp.MessageList)+1)
	for i := len(messageResp.MessageList) - 1; i >= 0; i-- {
		m := messageResp.MessageList[i]
		history = append(history, &rpcai.HistoryItem{Role: m.Role, Content: m.Content})
	}
	history = append(history, &rpcai.HistoryItem{Role: "user", Content: req.Message})
	stream, err := rpc.AiClient.Chat(h.Context, &rpcai.ChatReq{
		ProjectId: projectIDFromDir(item.ProjectDir),
		History:   history,
		Version:   version,
	})
	if err != nil {
		return sendChatErr(w, err)
	}
	defer stream.Close()
	var builder strings.Builder
	for {
		msg, recvErr := stream.Recv()
		if recvErr != nil {
			if recvErr == io.EOF {
				break
			}
			return sendChatErr(w, recvErr)
		}
		if msg == nil || msg.Content == "" {
			continue
		}
		builder.WriteString(msg.Content)
		payload, marshalErr := json.Marshal(graph.ServerSentEventString{D: msg.Content})
		if marshalErr != nil {
			return sendChatErr(w, marshalErr)
		}
		if writeErr := w.WriteEvent("", "message", payload); writeErr != nil {
			return &graph.ServerSentEventString{Message: writeErr.Error()}, writeErr
		}
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
	_ = w.WriteEvent("", "done", []byte("1"))
	return &graph.ServerSentEventString{D: answer, Message: "success"}, nil
}

func sendChatErr(w *sse.Writer, err error) (*graph.ServerSentEventString, error) {
	if err == nil {
		return &graph.ServerSentEventString{Message: "success"}, nil
	}
	payload, marshalErr := json.Marshal(graph.ServerSentEventString{Message: err.Error()})
	if marshalErr == nil {
		_ = w.WriteEvent("", "business-error", payload)
	}
	return &graph.ServerSentEventString{Message: err.Error()}, err
}
