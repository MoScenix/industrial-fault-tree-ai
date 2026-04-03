package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/dal/mysql"
	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/dal/redis"
	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/model"
	user "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/user"
)

type UpdateService struct {
	ctx context.Context
} // NewUpdateService new UpdateService
func NewUpdateService(ctx context.Context) *UpdateService {
	return &UpdateService{ctx: ctx}
}

// Run create note info
func (s *UpdateService) Run(req *user.UpdateReq) (resp *user.UpdateResp, err error) {
	// Finish your business logic.
	err = model.NewUserProQuery(s.ctx, mysql.DB, redis.RedisClient).UpdateUser(uint(req.Id), model.User{
		UserAvatar:  req.UserAvatar,
		Name:        req.UserName,
		UserProfile: req.UserProfile,
	})
	if err != nil {
		return nil, err
	}
	return &user.UpdateResp{}, nil
}
