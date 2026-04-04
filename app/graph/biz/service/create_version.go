package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/utils"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
)

type CreateVersionService struct {
	ctx context.Context
} // NewCreateVersionService new CreateVersionService
func NewCreateVersionService(ctx context.Context) *CreateVersionService {
	return &CreateVersionService{ctx: ctx}
}

// Run create note info
func (s *CreateVersionService) Run(req *graph.CreateVersionReq) (resp *graph.CreateVersionResp, err error) {
	item, err := mustLoadAuthorizedGraph(s.ctx, req.GraphId)
	if err != nil {
		return nil, err
	}

	currentVersion, err := utils.ReadCurrentVersion(item.ProjectDir)
	if err != nil {
		return nil, err
	}

	version := req.VersionName
	if version == "" {
		version = utils.FormatTimeNow()
	}

	if err := utils.CreateVersionFromCurrent(item.ProjectDir, currentVersion, version); err != nil {
		return nil, err
	}

	return &graph.CreateVersionResp{
		Success: true,
		Version: version,
	}, nil
}
