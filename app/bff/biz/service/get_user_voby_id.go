package service

import (
	"context"

	user "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/user"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcuser "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetUserVOByIdService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetUserVOByIdService(Context context.Context, RequestContext *app.RequestContext) *GetUserVOByIdService {
	return &GetUserVOByIdService{RequestContext: RequestContext, Context: Context}
}

func (h *GetUserVOByIdService) Run(req *user.GetUserVOByIdRequest) (resp *user.BaseResponseUserVO, err error) {
	res, err := rpc.UserClient.GetUser(h.Context, &rpcuser.GetUserReq{
		Id: req.Id,
	})
	if err != nil {
		return &user.BaseResponseUserVO{
			Code:    1,
			Message: err.Error(),
		}, nil
	}
	return &user.BaseResponseUserVO{
		Code:    0,
		Message: "success",
		Data: &user.UserVO{
			Id:          res.Id,
			UserName:    res.UserName,
			UserAccount: res.UserAccount,
			UserAvatar:  res.UserAvatar,
			UserProfile: res.UserProfile,
			UserRole:    res.UserRole,
			CreateTime:  res.CreateTime,
		},
	}, nil
}
