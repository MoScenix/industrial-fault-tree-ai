package service

import (
	"context"

	common "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/common"
	user "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/user"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/sessions"
)

type UserLogoutService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUserLogoutService(Context context.Context, RequestContext *app.RequestContext) *UserLogoutService {
	return &UserLogoutService{RequestContext: RequestContext, Context: Context}
}

func (h *UserLogoutService) Run(req *common.Empty) (resp *user.BaseResponseBoolean, err error) {
	sessions := sessions.Default(h.RequestContext)
	sessions.Clear()
	sessions.Save()
	return &user.BaseResponseBoolean{
		Code:    0,
		Message: "success",
		Data:    true,
	}, nil
}
