package service

import (
	"context"
	"os"
	"path/filepath"
	"strconv"

	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/biz/utils"
	user "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/user"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcuser "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateUserService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateUserService(Context context.Context, RequestContext *app.RequestContext) *UpdateUserService {
	return &UpdateUserService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateUserService) Run(req *user.UserUpdateRequest) (resp *user.BaseResponseBoolean, err error) {
	if req.Id != 0 {
		if int64(h.Context.Value(utils.UserIdKey).(float64)) != req.Id {
			return &user.BaseResponseBoolean{
				Code:    2,
				Message: "用户id不一致",
				Data:    false,
			}, nil
		}
	} else {
		req.Id = int64(h.Context.Value(utils.UserIdKey).(float64))
	}
	avatar, err := h.RequestContext.FormFile("avatar")
	if avatar != nil && err == nil {
		req.UserAvatar = "/static/avatar/" + strconv.FormatInt(req.Id, 10) + ".jpg"
		os.MkdirAll(filepath.Dir(req.UserAvatar), os.ModePerm)
		h.RequestContext.SaveUploadedFile(avatar, req.UserAvatar)
	}
	_, err = rpc.UserClient.Update(h.Context, &rpcuser.UpdateReq{
		Id:          req.Id,
		UserName:    req.UserName,
		UserAvatar:  req.UserAvatar,
		UserProfile: req.UserProfile,
		UserRole:    req.UserRole,
	})
	if err != nil {
		return &user.BaseResponseBoolean{
			Code:    1,
			Message: "更新用户信息失败",
			Data:    false,
		}, err
	}
	return &user.BaseResponseBoolean{
		Code:    0,
		Message: "success",
		Data:    true,
	}, err
}
