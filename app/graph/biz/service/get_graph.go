package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/dal/mysql"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/dal/redis"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/model"
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
	if mysql.DB == nil {
		return nil, errDBNotReady
	}

	item, err := model.NewGraphProQuery(s.ctx, mysql.DB, redis.RedisClient).GetGraphByID(uint(req.Id))
	if err != nil {
		return nil, err
	}

	return &graph.GetGraphResp{
		Graph: toGraphInfo(item),
	}, nil
}
