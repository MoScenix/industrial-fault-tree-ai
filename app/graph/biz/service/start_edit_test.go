package service

import (
	"context"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
	"testing"
)

func TestStartEdit_Run(t *testing.T) {
	ctx := context.Background()
	s := NewStartEditService(ctx)
	// init req and assert value

	req := &graph.StartEditReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
