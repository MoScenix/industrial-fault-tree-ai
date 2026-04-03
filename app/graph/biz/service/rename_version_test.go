package service

import (
	"context"
	"testing"
	graph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
)

func TestRenameVersion_Run(t *testing.T) {
	ctx := context.Background()
	s := NewRenameVersionService(ctx)
	// init req and assert value

	req := &graph.RenameVersionReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
