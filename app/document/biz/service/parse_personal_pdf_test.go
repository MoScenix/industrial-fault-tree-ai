package service

import (
	"context"
	"testing"
	document "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/document"
)

func TestParsePersonalPDF_Run(t *testing.T) {
	ctx := context.Background()
	s := NewParsePersonalPDFService(ctx)
	// init req and assert value

	req := &document.ParsePersonalPDFReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
