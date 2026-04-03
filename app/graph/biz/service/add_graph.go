package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/dal/mysql"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/dal/redis"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/model"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/utils"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
)

type AddGraphService struct {
	ctx context.Context
} // NewAddGraphService new AddGraphService
func NewAddGraphService(ctx context.Context) *AddGraphService {
	return &AddGraphService{ctx: ctx}
}

// Run create note info
func (s *AddGraphService) Run(req *graph.AddGraphReq) (resp *graph.AddGraphResp, err error) {
	if mysql.DB == nil {
		return nil, errDBNotReady
	}

	projectUUID := utils.NewProjectUUID()
	projectDir := utils.ProjectDir(projectUUID)
	graphName := req.GraphName
	if graphName == "" {
		graphName = "未命名图项目"
	}

	res, err := model.NewGraphProQuery(s.ctx, mysql.DB, redis.RedisClient).CreateGraph(model.Graph{
		GraphName:   graphName,
		Description: req.Description,
		Cover:       req.Cover,
		UserID:      uint(req.UserId),
		ProjectUUID: projectUUID,
		ProjectDir:  projectDir,
	})
	if err != nil {
		return nil, err
	}

	if err := utils.EnsureProjectLayout(projectDir); err != nil {
		_ = model.NewGraphProQuery(s.ctx, mysql.DB, redis.RedisClient).DeleteGraph(res.ID)
		return nil, err
	}

	return &graph.AddGraphResp{
		Id: int64(res.ID),
	}, nil
}
