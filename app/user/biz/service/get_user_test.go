package service

import (
	"context"
	"testing"
	user "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/user"
)

func TestGetUser_Run(t *testing.T) {
	ctx := context.Background()
	s := NewGetUserService(ctx)
	// init req and assert value

	req := &user.GetUserReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
