package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/biz/utils"
	user "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/user"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcuser "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetUserByIdService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetUserByIdService(Context context.Context, RequestContext *app.RequestContext) *GetUserByIdService {
	return &GetUserByIdService{RequestContext: RequestContext, Context: Context}
}

func (h *GetUserByIdService) Run(req *user.GetUserByIdRequest) (resp *user.BaseResponseUser, err error) {
	res, err := rpc.UserClient.GetUser(h.Context, &rpcuser.GetUserReq{
		Id: req.Id,
	})
	if err != nil {
		return &user.BaseResponseUser{
			Code:    1,
			Message: err.Error(),
		}, nil
	}
	return &user.BaseResponseUser{
		Code:    0,
		Message: "success",
		Data: &user.User{
			Id:          res.Id,
			UserName:    res.UserName,
			UserAccount: res.UserAccount,
			UserAvatar:  res.UserAvatar,
			UserProfile: res.UserProfile,
			UserRole:    res.UserRole,
			CreateTime:  res.CreateTime,
			UpdateTime:  res.UpdateTime,
			IsDelete:    utils.BoolToInt32(res.IsDelete),
		},
	}, nil
}
