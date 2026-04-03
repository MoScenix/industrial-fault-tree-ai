package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/dal/mysql"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/dal/redis"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/model"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
)

type AddGraphMessageService struct {
	ctx context.Context
}

func NewAddGraphMessageService(ctx context.Context) *AddGraphMessageService {
	return &AddGraphMessageService{ctx: ctx}
}

func (s *AddGraphMessageService) Run(req *graph.AddGraphMessageReq) (*graph.AddGraphMessageResp, error) {
	if mysql.DB == nil {
		return nil, errDBNotReady
	}

	item, err := model.NewMessageProQuery(s.ctx, mysql.DB, redis.RedisClient).CreateMessage(model.Message{
		GraphID: uint(req.GraphId),
		UserID:  uint(req.UserId),
		Role:    req.Role,
		Content: req.Content,
	})
	if err != nil {
		return nil, err
	}

	return &graph.AddGraphMessageResp{
		Id: int64(item.ID),
	}, nil
}
