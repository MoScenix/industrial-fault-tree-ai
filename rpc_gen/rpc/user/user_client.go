package user

import (
	"context"

	user "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/user"
	"github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/user/userservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() userservice.Client
	Service() string
	Login(ctx context.Context, req *user.LoginReq, callOptions ...callopt.Option) (r *user.LoginResp, err error)
	Register(ctx context.Context, req *user.RegisterReq, callOptions ...callopt.Option) (r *user.RegisterResp, err error)
	Update(ctx context.Context, req *user.UpdateReq, callOptions ...callopt.Option) (r *user.UpdateResp, err error)
	GetUser(ctx context.Context, req *user.GetUserReq, callOptions ...callopt.Option) (r *user.GetUserResp, err error)
	DeleteUser(ctx context.Context, req *user.DeleteUserReq, callOptions ...callopt.Option) (r *user.DeleteUserResp, err error)
	AddUser(ctx context.Context, req *user.AddUserReq, callOptions ...callopt.Option) (r *user.AddUserResp, err error)
	ListUser(ctx context.Context, req *user.ListUserReq, callOptions ...callopt.Option) (r *user.ListUserResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := userservice.NewClient(dstService, opts...)
	if err != nil {
		return nil, err
	}
	return &clientImpl{
		service:     dstService,
		kitexClient: kitexClient,
	}, nil
}

type clientImpl struct {
	service     string
	kitexClient userservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() userservice.Client {
	return c.kitexClient
}

func (c *clientImpl) Login(ctx context.Context, req *user.LoginReq, callOptions ...callopt.Option) (r *user.LoginResp, err error) {
	return c.kitexClient.Login(ctx, req, callOptions...)
}

func (c *clientImpl) Register(ctx context.Context, req *user.RegisterReq, callOptions ...callopt.Option) (r *user.RegisterResp, err error) {
	return c.kitexClient.Register(ctx, req, callOptions...)
}

func (c *clientImpl) Update(ctx context.Context, req *user.UpdateReq, callOptions ...callopt.Option) (r *user.UpdateResp, err error) {
	return c.kitexClient.Update(ctx, req, callOptions...)
}

func (c *clientImpl) GetUser(ctx context.Context, req *user.GetUserReq, callOptions ...callopt.Option) (r *user.GetUserResp, err error) {
	return c.kitexClient.GetUser(ctx, req, callOptions...)
}

func (c *clientImpl) DeleteUser(ctx context.Context, req *user.DeleteUserReq, callOptions ...callopt.Option) (r *user.DeleteUserResp, err error) {
	return c.kitexClient.DeleteUser(ctx, req, callOptions...)
}

func (c *clientImpl) AddUser(ctx context.Context, req *user.AddUserReq, callOptions ...callopt.Option) (r *user.AddUserResp, err error) {
	return c.kitexClient.AddUser(ctx, req, callOptions...)
}

func (c *clientImpl) ListUser(ctx context.Context, req *user.ListUserReq, callOptions ...callopt.Option) (r *user.ListUserResp, err error) {
	return c.kitexClient.ListUser(ctx, req, callOptions...)
}
