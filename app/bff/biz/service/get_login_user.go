package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/biz/utils"
	common "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/common"
	user "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/user"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcuser "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type GetLoginUserService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetLoginUserService(Context context.Context, RequestContext *app.RequestContext) *GetLoginUserService {
	return &GetLoginUserService{RequestContext: RequestContext, Context: Context}
}

func (h *GetLoginUserService) Run(req *common.Empty) (resp *user.BaseResponseLoginUserVO, err error) {
	if h.Context.Value(utils.UserIdKey) == nil {
		return &user.BaseResponseLoginUserVO{
			Code:    0,
			Message: "未登录",
		}, nil
	}
	res, err := rpc.UserClient.GetUser(h.Context, &rpcuser.GetUserReq{
		Id: int64(h.Context.Value(utils.UserIdKey).(float64)),
	})
	if err != nil {
		return &user.BaseResponseLoginUserVO{
			Code:    1,
			Message: err.Error(),
		}, nil
	}
	return &user.BaseResponseLoginUserVO{
		Code:    0,
		Message: "success",
		Data: &user.LoginUserVO{
			Id:          res.Id,
			UserName:    res.UserName,
			UserAccount: res.UserAccount,
			UserAvatar:  res.UserAvatar,
			UserProfile: res.UserProfile,
			UserRole:    res.UserRole,
			CreateTime:  res.CreateTime,
			UpdateTime:  res.UpdateTime,
		},
	}, nil
}
