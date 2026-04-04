package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/utils"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
)

type DeleteGraphService struct {
	ctx context.Context
} // NewDeleteGraphService new DeleteGraphService
func NewDeleteGraphService(ctx context.Context) *DeleteGraphService {
	return &DeleteGraphService{ctx: ctx}
}

// Run create note info
func (s *DeleteGraphService) Run(req *graph.DeleteGraphReq) (resp *graph.DeleteGraphResp, err error) {
	item, err := mustLoadAuthorizedGraph(s.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	q, err := graphQuery(s.ctx)
	if err != nil {
		return nil, err
	}

	if err := q.DeleteGraph(item.ID); err != nil {
		return nil, err
	}
	if err := utils.RemoveProjectLayout(item.ProjectDir); err != nil {
		return nil, err
	}

	return &graph.DeleteGraphResp{
		Success: true,
	}, nil
}
