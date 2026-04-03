package graph

import (
	"context"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func AddGraph(ctx context.Context, req *graph.AddGraphReq, callOptions ...callopt.Option) (resp *graph.AddGraphResp, err error) {
	resp, err = defaultClient.AddGraph(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "AddGraph call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func DeleteGraph(ctx context.Context, req *graph.DeleteGraphReq, callOptions ...callopt.Option) (resp *graph.DeleteGraphResp, err error) {
	resp, err = defaultClient.DeleteGraph(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "DeleteGraph call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func UpdateGraph(ctx context.Context, req *graph.UpdateGraphReq, callOptions ...callopt.Option) (resp *graph.UpdateGraphResp, err error) {
	resp, err = defaultClient.UpdateGraph(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "UpdateGraph call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func GetGraph(ctx context.Context, req *graph.GetGraphReq, callOptions ...callopt.Option) (resp *graph.GetGraphResp, err error) {
	resp, err = defaultClient.GetGraph(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "GetGraph call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func ListGraph(ctx context.Context, req *graph.ListGraphReq, callOptions ...callopt.Option) (resp *graph.ListGraphResp, err error) {
	resp, err = defaultClient.ListGraph(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "ListGraph call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func StartEdit(ctx context.Context, req *graph.StartEditReq, callOptions ...callopt.Option) (resp *graph.StartEditResp, err error) {
	resp, err = defaultClient.StartEdit(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "StartEdit call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func Save(ctx context.Context, req *graph.SaveReq, callOptions ...callopt.Option) (resp *graph.SaveResp, err error) {
	resp, err = defaultClient.Save(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "Save call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func AddGraphMessage(ctx context.Context, req *graph.AddGraphMessageReq, callOptions ...callopt.Option) (resp *graph.AddGraphMessageResp, err error) {
	resp, err = defaultClient.AddGraphMessage(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "AddGraphMessage call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func ListGraphMessage(ctx context.Context, req *graph.ListGraphMessageReq, callOptions ...callopt.Option) (resp *graph.ListGraphMessageResp, err error) {
	resp, err = defaultClient.ListGraphMessage(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "ListGraphMessage call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func CreateVersion(ctx context.Context, req *graph.CreateVersionReq, callOptions ...callopt.Option) (resp *graph.CreateVersionResp, err error) {
	resp, err = defaultClient.CreateVersion(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "CreateVersion call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func DeleteVersion(ctx context.Context, req *graph.DeleteVersionReq, callOptions ...callopt.Option) (resp *graph.DeleteVersionResp, err error) {
	resp, err = defaultClient.DeleteVersion(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "DeleteVersion call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func RenameVersion(ctx context.Context, req *graph.RenameVersionReq, callOptions ...callopt.Option) (resp *graph.RenameVersionResp, err error) {
	resp, err = defaultClient.RenameVersion(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "RenameVersion call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func ListVersion(ctx context.Context, req *graph.ListVersionReq, callOptions ...callopt.Option) (resp *graph.ListVersionResp, err error) {
	resp, err = defaultClient.ListVersion(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "ListVersion call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
