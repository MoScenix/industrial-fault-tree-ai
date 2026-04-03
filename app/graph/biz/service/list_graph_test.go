package service

import (
	"context"
	"testing"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
)

func TestListGraph_Run(t *testing.T) {
	ctx := context.Background()
	s := NewListGraphService(ctx)
	// init req and assert value

	req := &graph.ListGraphReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
