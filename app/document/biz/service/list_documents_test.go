package service

import (
	"context"
	"testing"
	document "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/document"
)

func TestListDocuments_Run(t *testing.T) {
	ctx := context.Background()
	s := NewListDocumentsService(ctx)
	// init req and assert value

	req := &document.ListDocumentsReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
