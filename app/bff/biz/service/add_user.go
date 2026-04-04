package service

import (
	"context"

	user "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/user"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcuser "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type AddUserService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAddUserService(Context context.Context, RequestContext *app.RequestContext) *AddUserService {
	return &AddUserService{RequestContext: RequestContext, Context: Context}
}

func (h *AddUserService) Run(req *user.UserAddRequest) (resp *user.BaseResponseLong, err error) {
	_, err = rpc.UserClient.AddUser(h.Context, &rpcuser.AddUserReq{
		UserAccount:  req.UserAccount,
		UserAvatar:   req.UserAvatar,
		UserName:     req.UserName,
		UserProfile:  req.UserProfile,
		UserRole:     req.UserRole,
		UserPassword: req.UserName,
	})
	if err != nil {
		return &user.BaseResponseLong{
			Code:    1,
			Message: err.Error(),
		}, err
	}
	return &user.BaseResponseLong{
		Code:    0,
		Message: "success",
	}, nil
}
