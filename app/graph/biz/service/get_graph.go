package service

import (
	"context"

	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
)

type GetGraphService struct {
	ctx context.Context
} // NewGetGraphService new GetGraphService
func NewGetGraphService(ctx context.Context) *GetGraphService {
	return &GetGraphService{ctx: ctx}
}

// Run create note info
func (s *GetGraphService) Run(req *graph.GetGraphReq) (resp *graph.GetGraphResp, err error) {
	item, err := mustLoadAuthorizedGraph(s.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &graph.GetGraphResp{
		Graph: toGraphInfo(item),
	}, nil
}
