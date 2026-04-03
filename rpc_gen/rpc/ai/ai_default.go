package ai

import (
	"context"
	ai "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/ai"
	"github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/ai/aiservice"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func Chat(ctx context.Context, Req *ai.ChatReq, callOptions ...callopt.Option) (stream aiservice.AiService_ChatClient, err error) {
	stream, err = defaultClient.Chat(ctx, Req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "Chat call failed,err =%+v", err)
		return nil, err
	}
	return stream, nil
}

func Validate(ctx context.Context, req *ai.ValidateReq, callOptions ...callopt.Option) (resp *ai.ValidateResp, err error) {
	resp, err = defaultClient.Validate(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "Validate call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func UpdatePrompt(ctx context.Context, req *ai.UpdatePromptReq, callOptions ...callopt.Option) (resp *ai.UpdatePromptResp, err error) {
	resp, err = defaultClient.UpdatePrompt(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "UpdatePrompt call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
