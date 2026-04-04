package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/dal/mysql"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/dal/redis"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/model"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/utils"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
)

type SaveOperationService struct {
	ctx context.Context
} // NewSaveService new SaveService
func NewSaveService(ctx context.Context) *SaveOperationService {
	return &SaveOperationService{ctx: ctx}
}

// Run create note info
func (s *SaveOperationService) Run(req *graph.SaveReq) (resp *graph.SaveResp, err error) {
	if mysql.DB == nil {
		return nil, errDBNotReady
	}

	q := model.NewGraphProQuery(s.ctx, mysql.DB, redis.RedisClient)
	item, err := q.GetGraphByID(uint(req.GraphId))
	if err != nil {
		return nil, err
	}

	fromVersion := req.FromVersion
	if fromVersion == "" {
		fromVersion, err = utils.ReadCurrentVersion(item.ProjectDir)
		if err != nil {
			return nil, err
		}
	}
	toVersion := req.ToVersion
	if toVersion == "" {
		toVersion = fromVersion
	}
	if !utils.HasTmpVersion(item.ProjectDir, fromVersion) {
		return &graph.SaveResp{
			Success:     false,
			FromVersion: fromVersion,
			ToVersion:   toVersion,
			Message:     "tmp not ready",
		}, nil
	}
	if err := utils.SaveTmpToVersion(item.ProjectDir, fromVersion, toVersion); err != nil {
		return nil, err
	}
	if err := utils.WriteCurrentVersion(item.ProjectDir, toVersion); err != nil {
		return nil, err
	}
	if err := utils.ClearTmp(item.ProjectDir, fromVersion); err != nil {
		return nil, err
	}

	return &graph.SaveResp{
		Success:     true,
		FromVersion: fromVersion,
		ToVersion:   toVersion,
		Message:     "save success",
	}, nil
}
