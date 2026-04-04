package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/utils"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
)

type RenameVersionService struct {
	ctx context.Context
} // NewRenameVersionService new RenameVersionService
func NewRenameVersionService(ctx context.Context) *RenameVersionService {
	return &RenameVersionService{ctx: ctx}
}

// Run create note info
func (s *RenameVersionService) Run(req *graph.RenameVersionReq) (resp *graph.RenameVersionResp, err error) {
	item, err := mustLoadAuthorizedGraph(s.ctx, req.GraphId)
	if err != nil {
		return nil, err
	}

	currentVersion, err := utils.ReadCurrentVersion(item.ProjectDir)
	if err != nil {
		return nil, err
	}

	if err := utils.RenameVersionDir(item.ProjectDir, req.Version, req.VersionName); err != nil {
		return nil, err
	}
	if currentVersion == req.Version {
		if err := utils.WriteCurrentVersion(item.ProjectDir, req.VersionName); err != nil {
			return nil, err
		}
	}

	return &graph.RenameVersionResp{Success: true}, nil
}
