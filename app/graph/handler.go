package main

import (
	"context"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/service"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
)

// GraphServiceImpl implements the last service interface defined in the IDL.
type GraphServiceImpl struct{}

// AddGraph implements the GraphServiceImpl interface.
func (s *GraphServiceImpl) AddGraph(ctx context.Context, req *graph.AddGraphReq) (resp *graph.AddGraphResp, err error) {
	resp, err = service.NewAddGraphService(ctx).Run(req)

	return resp, err
}

// DeleteGraph implements the GraphServiceImpl interface.
func (s *GraphServiceImpl) DeleteGraph(ctx context.Context, req *graph.DeleteGraphReq) (resp *graph.DeleteGraphResp, err error) {
	resp, err = service.NewDeleteGraphService(ctx).Run(req)

	return resp, err
}

// UpdateGraph implements the GraphServiceImpl interface.
func (s *GraphServiceImpl) UpdateGraph(ctx context.Context, req *graph.UpdateGraphReq) (resp *graph.UpdateGraphResp, err error) {
	resp, err = service.NewUpdateGraphService(ctx).Run(req)

	return resp, err
}

// GetGraph implements the GraphServiceImpl interface.
func (s *GraphServiceImpl) GetGraph(ctx context.Context, req *graph.GetGraphReq) (resp *graph.GetGraphResp, err error) {
	resp, err = service.NewGetGraphService(ctx).Run(req)

	return resp, err
}

// ListGraph implements the GraphServiceImpl interface.
func (s *GraphServiceImpl) ListGraph(ctx context.Context, req *graph.ListGraphReq) (resp *graph.ListGraphResp, err error) {
	resp, err = service.NewListGraphService(ctx).Run(req)

	return resp, err
}

// StartEdit implements the GraphServiceImpl interface.
func (s *GraphServiceImpl) StartEdit(ctx context.Context, req *graph.StartEditReq) (resp *graph.StartEditResp, err error) {
	resp, err = service.NewStartEditService(ctx).Run(req)

	return resp, err
}

// Save implements the GraphServiceImpl interface.
func (s *GraphServiceImpl) Save(ctx context.Context, req *graph.SaveReq) (resp *graph.SaveResp, err error) {
	resp, err = service.NewSaveService(ctx).Run(req)

	return resp, err
}

// AddGraphMessage implements the GraphServiceImpl interface.
func (s *GraphServiceImpl) AddGraphMessage(ctx context.Context, req *graph.AddGraphMessageReq) (resp *graph.AddGraphMessageResp, err error) {
	resp, err = service.NewAddGraphMessageService(ctx).Run(req)

	return resp, err
}

// ListGraphMessage implements the GraphServiceImpl interface.
func (s *GraphServiceImpl) ListGraphMessage(ctx context.Context, req *graph.ListGraphMessageReq) (resp *graph.ListGraphMessageResp, err error) {
	resp, err = service.NewListGraphMessageService(ctx).Run(req)

	return resp, err
}

// CreateVersion implements the GraphServiceImpl interface.
func (s *GraphServiceImpl) CreateVersion(ctx context.Context, req *graph.CreateVersionReq) (resp *graph.CreateVersionResp, err error) {
	resp, err = service.NewCreateVersionService(ctx).Run(req)

	return resp, err
}

// DeleteVersion implements the GraphServiceImpl interface.
func (s *GraphServiceImpl) DeleteVersion(ctx context.Context, req *graph.DeleteVersionReq) (resp *graph.DeleteVersionResp, err error) {
	resp, err = service.NewDeleteVersionService(ctx).Run(req)

	return resp, err
}

// RenameVersion implements the GraphServiceImpl interface.
func (s *GraphServiceImpl) RenameVersion(ctx context.Context, req *graph.RenameVersionReq) (resp *graph.RenameVersionResp, err error) {
	resp, err = service.NewRenameVersionService(ctx).Run(req)

	return resp, err
}

// ListVersion implements the GraphServiceImpl interface.
func (s *GraphServiceImpl) ListVersion(ctx context.Context, req *graph.ListVersionReq) (resp *graph.ListVersionResp, err error) {
	resp, err = service.NewListVersionService(ctx).Run(req)

	return resp, err
}
