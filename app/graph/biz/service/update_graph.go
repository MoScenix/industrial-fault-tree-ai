package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/model"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
)

type UpdateGraphService struct {
	ctx context.Context
} // NewUpdateGraphService new UpdateGraphService
func NewUpdateGraphService(ctx context.Context) *UpdateGraphService {
	return &UpdateGraphService{ctx: ctx}
}

// Run create note info
func (s *UpdateGraphService) Run(req *graph.UpdateGraphReq) (resp *graph.UpdateGraphResp, err error) {
	current, err := mustLoadAuthorizedGraph(s.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	q, err := graphQuery(s.ctx)
	if err != nil {
		return nil, err
	}

	update := model.Graph{}
	if req.GraphName != "" {
		update.GraphName = req.GraphName
		current.GraphName = req.GraphName
	}
	if req.Description != "" {
		update.Description = req.Description
		current.Description = req.Description
	}
	if req.Cover != "" {
		update.Cover = req.Cover
		current.Cover = req.Cover
	}

	if err := q.UpdateGraph(uint(req.Id), update); err != nil {
		return nil, err
	}
	return &graph.UpdateGraphResp{
		Success: true,
	}, nil
}
