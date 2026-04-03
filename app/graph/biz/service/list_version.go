package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/dal/mysql"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/dal/redis"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/model"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/utils"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
)

type ListVersionService struct {
	ctx context.Context
} // NewListVersionService new ListVersionService
func NewListVersionService(ctx context.Context) *ListVersionService {
	return &ListVersionService{ctx: ctx}
}

// Run create note info
func (s *ListVersionService) Run(req *graph.ListVersionReq) (resp *graph.ListVersionResp, err error) {
	if mysql.DB == nil {
		return nil, errDBNotReady
	}

	item, err := model.NewGraphProQuery(s.ctx, mysql.DB, redis.RedisClient).GetGraphByID(uint(req.GraphId))
	if err != nil {
		return nil, err
	}

	currentVersion, err := utils.ReadCurrentVersion(item.ProjectDir)
	if err != nil {
		return nil, err
	}

	versionList, err := utils.ListVersions(item.ProjectDir, currentVersion)
	if err != nil {
		return nil, err
	}

	resp = &graph.ListVersionResp{}
	for _, v := range versionList {
		resp.VersionList = append(resp.VersionList, &graph.GraphVersionInfo{
			Version:     v.Version,
			VersionName: v.Version,
			IsCurrent:   v.IsCurrent,
			CreateTime:  v.CreateTime,
			UpdateTime:  v.UpdateTime,
		})
	}
	return resp, nil
}
