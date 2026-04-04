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

type UserRegisterService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUserRegisterService(Context context.Context, RequestContext *app.RequestContext) *UserRegisterService {
	return &UserRegisterService{RequestContext: RequestContext, Context: Context}
}

func (h *UserRegisterService) Run(req *user.UserRegisterRequest) (resp *user.BaseResponseLong, err error) {
	if req.UserPassword != req.CheckPassword {
		return &user.BaseResponseLong{
			Code:    2,
			Message: "两次密码不一致",
		}, nil
	}
	res, err := rpc.UserClient.Register(h.Context, &rpcuser.RegisterReq{
		UserAccount:  req.UserAccount,
		UserPassword: req.UserPassword,
	})

	if err != nil {
		return &user.BaseResponseLong{
			Code:    1,
			Message: err.Error(),
		}, nil
	}
	session := sessions.Default(h.RequestContext)
	session.Set(utils.UserIdKey, res.UserId)
	session.Set(utils.UserRoleKey, res.UserRole)
	session.Save()
	return &user.BaseResponseLong{
		Code:    0,
		Message: "success",
		Data:    int64(res.UserId),
	}, nil
}
