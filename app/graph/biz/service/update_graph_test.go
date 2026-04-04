package service

import (
	"context"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
	"testing"
)

func TestUpdateGraph_Run(t *testing.T) {
	ctx := context.Background()
	s := NewUpdateGraphService(ctx)
	// init req and assert value

	req := &graph.UpdateGraphReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
