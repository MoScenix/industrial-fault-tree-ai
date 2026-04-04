package service

import (
	"context"

	user "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/user"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcuser "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type DeleteUserService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDeleteUserService(Context context.Context, RequestContext *app.RequestContext) *DeleteUserService {
	return &DeleteUserService{RequestContext: RequestContext, Context: Context}
}

func (h *DeleteUserService) Run(req *user.DeleteRequest) (resp *user.BaseResponseBoolean, err error) {
	_, err = rpc.UserClient.DeleteUser(h.Context, &rpcuser.DeleteUserReq{
		UserId: req.Id,
	})
	if err != nil {
		return &user.BaseResponseBoolean{
			Code:    1,
			Message: err.Error(),
			Data:    false,
		}, nil
	}
	return &user.BaseResponseBoolean{
		Code:    0,
		Message: "success",
		Data:    true,
	}, nil
}
