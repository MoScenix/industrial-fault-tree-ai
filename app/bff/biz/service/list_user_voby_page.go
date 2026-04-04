package service

import (
	"context"

	user "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/user"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcuser "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type ListUserVOByPageService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListUserVOByPageService(Context context.Context, RequestContext *app.RequestContext) *ListUserVOByPageService {
	return &ListUserVOByPageService{RequestContext: RequestContext, Context: Context}
}

func (h *ListUserVOByPageService) Run(req *user.UserQueryRequest) (resp *user.BaseResponsePageUserVO, err error) {
	res, err := rpc.UserClient.ListUser(h.Context, &rpcuser.ListUserReq{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		UserName: req.UserName,
		Account:  req.UserAccount,
	})
	if err != nil {
		return &user.BaseResponsePageUserVO{
			Code:    1,
			Message: err.Error(),
		}, err
	}
	var userVOList []*user.UserVO
	for _, v := range res.UserList {
		userVOList = append(userVOList, &user.UserVO{
			Id:          v.Id,
			UserName:    v.UserName,
			UserAccount: v.UserAccount,
			UserAvatar:  v.UserAvatar,
			UserProfile: v.UserProfile,
			UserRole:    v.UserRole,
			CreateTime:  v.CreateTime,
		})
	}
	resp = &user.BaseResponsePageUserVO{
		Code:    0,
		Message: "success",
		Data: &user.PageUserVO{
			Records:    userVOList,
			PageNumber: req.PageNum,
			PageSize:   req.PageSize,
			TotalPage:  res.Total,
		},
	}
	return resp, nil
}
