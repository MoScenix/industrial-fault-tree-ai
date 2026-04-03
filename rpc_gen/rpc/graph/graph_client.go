package graph

import (
	"context"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"

	"github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph/graphservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() graphservice.Client
	Service() string
	AddGraph(ctx context.Context, Req *graph.AddGraphReq, callOptions ...callopt.Option) (r *graph.AddGraphResp, err error)
	DeleteGraph(ctx context.Context, Req *graph.DeleteGraphReq, callOptions ...callopt.Option) (r *graph.DeleteGraphResp, err error)
	UpdateGraph(ctx context.Context, Req *graph.UpdateGraphReq, callOptions ...callopt.Option) (r *graph.UpdateGraphResp, err error)
	GetGraph(ctx context.Context, Req *graph.GetGraphReq, callOptions ...callopt.Option) (r *graph.GetGraphResp, err error)
	ListGraph(ctx context.Context, Req *graph.ListGraphReq, callOptions ...callopt.Option) (r *graph.ListGraphResp, err error)
	StartEdit(ctx context.Context, Req *graph.StartEditReq, callOptions ...callopt.Option) (r *graph.StartEditResp, err error)
	Save(ctx context.Context, Req *graph.SaveReq, callOptions ...callopt.Option) (r *graph.SaveResp, err error)
	AddGraphMessage(ctx context.Context, Req *graph.AddGraphMessageReq, callOptions ...callopt.Option) (r *graph.AddGraphMessageResp, err error)
	ListGraphMessage(ctx context.Context, Req *graph.ListGraphMessageReq, callOptions ...callopt.Option) (r *graph.ListGraphMessageResp, err error)
	CreateVersion(ctx context.Context, Req *graph.CreateVersionReq, callOptions ...callopt.Option) (r *graph.CreateVersionResp, err error)
	DeleteVersion(ctx context.Context, Req *graph.DeleteVersionReq, callOptions ...callopt.Option) (r *graph.DeleteVersionResp, err error)
	RenameVersion(ctx context.Context, Req *graph.RenameVersionReq, callOptions ...callopt.Option) (r *graph.RenameVersionResp, err error)
	ListVersion(ctx context.Context, Req *graph.ListVersionReq, callOptions ...callopt.Option) (r *graph.ListVersionResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := graphservice.NewClient(dstService, opts...)
	if err != nil {
		return nil, err
	}
	cli := &clientImpl{
		service:     dstService,
		kitexClient: kitexClient,
	}

	return cli, nil
}

type clientImpl struct {
	service     string
	kitexClient graphservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() graphservice.Client {
	return c.kitexClient
}

func (c *clientImpl) AddGraph(ctx context.Context, Req *graph.AddGraphReq, callOptions ...callopt.Option) (r *graph.AddGraphResp, err error) {
	return c.kitexClient.AddGraph(ctx, Req, callOptions...)
}

func (c *clientImpl) DeleteGraph(ctx context.Context, Req *graph.DeleteGraphReq, callOptions ...callopt.Option) (r *graph.DeleteGraphResp, err error) {
	return c.kitexClient.DeleteGraph(ctx, Req, callOptions...)
}

func (c *clientImpl) UpdateGraph(ctx context.Context, Req *graph.UpdateGraphReq, callOptions ...callopt.Option) (r *graph.UpdateGraphResp, err error) {
	return c.kitexClient.UpdateGraph(ctx, Req, callOptions...)
}

func (c *clientImpl) GetGraph(ctx context.Context, Req *graph.GetGraphReq, callOptions ...callopt.Option) (r *graph.GetGraphResp, err error) {
	return c.kitexClient.GetGraph(ctx, Req, callOptions...)
}

func (c *clientImpl) ListGraph(ctx context.Context, Req *graph.ListGraphReq, callOptions ...callopt.Option) (r *graph.ListGraphResp, err error) {
	return c.kitexClient.ListGraph(ctx, Req, callOptions...)
}

func (c *clientImpl) StartEdit(ctx context.Context, Req *graph.StartEditReq, callOptions ...callopt.Option) (r *graph.StartEditResp, err error) {
	return c.kitexClient.StartEdit(ctx, Req, callOptions...)
}

func (c *clientImpl) Save(ctx context.Context, Req *graph.SaveReq, callOptions ...callopt.Option) (r *graph.SaveResp, err error) {
	return c.kitexClient.Save(ctx, Req, callOptions...)
}

func (c *clientImpl) AddGraphMessage(ctx context.Context, Req *graph.AddGraphMessageReq, callOptions ...callopt.Option) (r *graph.AddGraphMessageResp, err error) {
	return c.kitexClient.AddGraphMessage(ctx, Req, callOptions...)
}

func (c *clientImpl) ListGraphMessage(ctx context.Context, Req *graph.ListGraphMessageReq, callOptions ...callopt.Option) (r *graph.ListGraphMessageResp, err error) {
	return c.kitexClient.ListGraphMessage(ctx, Req, callOptions...)
}

func (c *clientImpl) CreateVersion(ctx context.Context, Req *graph.CreateVersionReq, callOptions ...callopt.Option) (r *graph.CreateVersionResp, err error) {
	return c.kitexClient.CreateVersion(ctx, Req, callOptions...)
}

func (c *clientImpl) DeleteVersion(ctx context.Context, Req *graph.DeleteVersionReq, callOptions ...callopt.Option) (r *graph.DeleteVersionResp, err error) {
	return c.kitexClient.DeleteVersion(ctx, Req, callOptions...)
}

func (c *clientImpl) RenameVersion(ctx context.Context, Req *graph.RenameVersionReq, callOptions ...callopt.Option) (r *graph.RenameVersionResp, err error) {
	return c.kitexClient.RenameVersion(ctx, Req, callOptions...)
}

func (c *clientImpl) ListVersion(ctx context.Context, Req *graph.ListVersionReq, callOptions ...callopt.Option) (r *graph.ListVersionResp, err error) {
	return c.kitexClient.ListVersion(ctx, Req, callOptions...)
}
