package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/utils"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
)

type DeleteVersionService struct {
	ctx context.Context
} // NewDeleteVersionService new DeleteVersionService
func NewDeleteVersionService(ctx context.Context) *DeleteVersionService {
	return &DeleteVersionService{ctx: ctx}
}

// Run create note info
func (s *DeleteVersionService) Run(req *graph.DeleteVersionReq) (resp *graph.DeleteVersionResp, err error) {
	item, err := mustLoadAuthorizedGraph(s.ctx, req.GraphId)
	if err != nil {
		return nil, err
	}

	currentVersion, err := utils.ReadCurrentVersion(item.ProjectDir)
	if err != nil {
		return nil, err
	}
	if req.Version == currentVersion {
		return &graph.DeleteVersionResp{Success: false}, nil
	}

	if err := utils.DeleteVersionDir(item.ProjectDir, req.Version); err != nil {
		return nil, err
	}

	return &graph.DeleteVersionResp{Success: true}, nil
}
