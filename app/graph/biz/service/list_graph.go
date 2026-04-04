package service

import (
	"context"

	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
)

type ListGraphService struct {
	ctx context.Context
} // NewListGraphService new ListGraphService
func NewListGraphService(ctx context.Context) *ListGraphService {
	return &ListGraphService{ctx: ctx}
}

// Run create note info
func (s *ListGraphService) Run(req *graph.ListGraphReq) (resp *graph.ListGraphResp, err error) {
	userID, err := effectiveListUserID(s.ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	q, err := graphQuery(s.ctx)
	if err != nil {
		return nil, err
	}
	total, err := q.CountGraph(userID, req.GraphName)
	if err != nil {
		return nil, err
	}

	items, err := q.ListGraph(uint32(req.PageNum), userID, req.GraphName, uint32(req.PageSize))
	if err != nil {
		return nil, err
	}

	resp = &graph.ListGraphResp{
		Total: total,
	}
	for _, item := range items {
		resp.GraphList = append(resp.GraphList, toGraphInfo(item))
	}
	return resp, nil
}
