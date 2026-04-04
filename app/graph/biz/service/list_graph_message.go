package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/dal/mysql"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/dal/redis"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/model"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/utils"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
)

type ListGraphMessageService struct {
	ctx context.Context
} // NewListGraphMessageService new ListGraphMessageService
func NewListGraphMessageService(ctx context.Context) *ListGraphMessageService {
	return &ListGraphMessageService{ctx: ctx}
}

// Run create note info
func (s *ListGraphMessageService) Run(req *graph.ListGraphMessageReq) (resp *graph.ListGraphMessageResp, err error) {
	if mysql.DB == nil {
		return nil, errDBNotReady
	}
	if _, err := mustLoadAuthorizedGraph(s.ctx, req.GraphId); err != nil {
		return nil, err
	}

	q := model.NewMessageProQuery(s.ctx, mysql.DB, redis.RedisClient)
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}
	total, err := q.Count(uint(req.GraphId))
	if err != nil {
		return nil, err
	}

	cursor, err := utils.ParseTimePtr(req.LastCreateTime)
	if err != nil {
		return nil, err
	}

	items, err := q.ListMessagesByGraphID(uint(req.GraphId), pageSize, cursor)
	if err != nil {
		return nil, err
	}

	resp = &graph.ListGraphMessageResp{
		Total: total,
	}
	for _, item := range items {
		resp.MessageList = append(resp.MessageList, toGraphMessage(item))
	}
	resp.HasMore = len(items) == pageSize && int64(len(items)) < total
	return resp, nil
}
