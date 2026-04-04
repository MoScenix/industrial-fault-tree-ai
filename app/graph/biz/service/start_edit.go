package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/utils"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
)

type StartEditService struct {
	ctx context.Context
} // NewStartEditService new StartEditService
func NewStartEditService(ctx context.Context) *StartEditService {
	return &StartEditService{ctx: ctx}
}

// Run create note info
func (s *StartEditService) Run(req *graph.StartEditReq) (resp *graph.StartEditResp, err error) {
	item, err := mustLoadAuthorizedGraph(s.ctx, req.GraphId)
	if err != nil {
		return nil, err
	}

	basedOnVersion := req.Version
	if basedOnVersion == "" {
		basedOnVersion, err = utils.ReadCurrentVersion(item.ProjectDir)
		if err != nil {
			return nil, err
		}
	}

	if err := utils.EnsureTmpFromVersion(item.ProjectDir, basedOnVersion); err != nil {
		return nil, err
	}

	return &graph.StartEditResp{
		Success:        true,
		TmpReady:       true,
		BasedOnVersion: basedOnVersion,
		Message:        "tmp ready",
	}, nil
}
