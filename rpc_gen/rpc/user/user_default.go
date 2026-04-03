package user

import (
	"context"

	user "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func Login(ctx context.Context, req *user.LoginReq, callOptions ...callopt.Option) (resp *user.LoginResp, err error) {
	resp, err = defaultClient.Login(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "Login call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func Register(ctx context.Context, req *user.RegisterReq, callOptions ...callopt.Option) (resp *user.RegisterResp, err error) {
	resp, err = defaultClient.Register(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "Register call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func Update(ctx context.Context, req *user.UpdateReq, callOptions ...callopt.Option) (resp *user.UpdateResp, err error) {
	resp, err = defaultClient.Update(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "Update call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func GetUser(ctx context.Context, req *user.GetUserReq, callOptions ...callopt.Option) (resp *user.GetUserResp, err error) {
	resp, err = defaultClient.GetUser(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "GetUser call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func DeleteUser(ctx context.Context, req *user.DeleteUserReq, callOptions ...callopt.Option) (resp *user.DeleteUserResp, err error) {
	resp, err = defaultClient.DeleteUser(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "DeleteUser call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func AddUser(ctx context.Context, req *user.AddUserReq, callOptions ...callopt.Option) (resp *user.AddUserResp, err error) {
	resp, err = defaultClient.AddUser(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "AddUser call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func ListUser(ctx context.Context, req *user.ListUserReq, callOptions ...callopt.Option) (resp *user.ListUserResp, err error) {
	resp, err = defaultClient.ListUser(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "ListUser call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
