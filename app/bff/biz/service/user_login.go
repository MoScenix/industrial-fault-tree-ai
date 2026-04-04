package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/biz/utils"
	user "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/user"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcuser "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/sessions"
)

type UserLoginService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUserLoginService(Context context.Context, RequestContext *app.RequestContext) *UserLoginService {
	return &UserLoginService{RequestContext: RequestContext, Context: Context}
}

func (h *UserLoginService) Run(req *user.UserLoginRequest) (resp *user.BaseResponseLoginUserVO, err error) {
	res1, err := rpc.UserClient.Login(h.Context, &rpcuser.LoginReq{
		UserAccount:  req.UserAccount,
		UserPassword: req.UserPassword,
	})
	if err != nil {
		return &user.BaseResponseLoginUserVO{
			Code:    1,
			Message: err.Error(),
		}, nil
	}
	session := sessions.Default(h.RequestContext)
	session.Set(utils.UserIdKey, res1.UserId)
	session.Set(utils.UserRoleKey, res1.UserRole)
	session.Save()
	res, err := rpc.UserClient.GetUser(h.Context, &rpcuser.GetUserReq{
		Id: int64(res1.UserId),
	})
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
