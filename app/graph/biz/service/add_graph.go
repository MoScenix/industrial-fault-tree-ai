package service

import (
	"context"

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
	userID, err := ensureLogin(s.ctx)
	if err != nil {
		return nil, err
	}
	q, err := graphQuery(s.ctx)
	if err != nil {
		return nil, err
	}

	projectUUID := utils.NewProjectUUID()
	projectDir := utils.ProjectDir(projectUUID)
	graphName := req.GraphName
	if graphName == "" {
		graphName = "未命名图项目"
	}

	res, err := q.CreateGraph(model.Graph{
		GraphName:   graphName,
		Description: req.Description,
		Cover:       req.Cover,
		UserID:      uint(userID),
		ProjectUUID: projectUUID,
		ProjectDir:  projectDir,
	})
	if err != nil {
		return nil, err
	}

	if err := utils.EnsureProjectLayout(projectDir); err != nil {
		_ = q.DeleteGraph(res.ID)
		return nil, err
	}

	return &graph.AddGraphResp{
		Id: int64(res.ID),
	}, nil
}
