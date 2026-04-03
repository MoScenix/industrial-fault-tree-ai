package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/dal/mysql"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/dal/redis"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/model"
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
	if mysql.DB == nil {
		return nil, errDBNotReady
	}

	q := model.NewGraphProQuery(s.ctx, mysql.DB, redis.RedisClient)
	total, err := q.CountGraph(uint(req.UserId), req.GraphName)
	if err != nil {
		return nil, err
	}

	items, err := q.ListGraph(uint32(req.PageNum), uint(req.UserId), req.GraphName, uint32(req.PageSize))
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
