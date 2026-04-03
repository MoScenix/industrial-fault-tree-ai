package service

import (
	"context"
	"testing"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
)

func TestCreateVersion_Run(t *testing.T) {
	ctx := context.Background()
	s := NewCreateVersionService(ctx)
	// init req and assert value

	req := &graph.CreateVersionReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
