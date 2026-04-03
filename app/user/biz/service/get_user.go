package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/dal/mysql"
	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/dal/redis"
	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/model"
	user "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/user"
)

type GetUserService struct {
	ctx context.Context
} // NewGetUserService new GetUserService
func NewGetUserService(ctx context.Context) *GetUserService {
	return &GetUserService{ctx: ctx}
}

// Run create note info
func (s *GetUserService) Run(req *user.GetUserReq) (resp *user.GetUserResp, err error) {
	// Finish your business logic.
	q := model.NewUserProQuery(s.ctx, mysql.DB, redis.RedisClient)
	User, err := q.GetUserById(int(req.Id))
	if err != nil {
		return nil, err
	}
	return &user.GetUserResp{
		Id:          int64(User.ID),
		UserName:    User.Name,
		UserAccount: User.UserAccount,
		UserAvatar:  User.UserAvatar,
		UserProfile: User.UserProfile,
		UserRole:    User.UserRole,
		IsDelete:    User.DeletedAt.Valid,
		CreateTime:  User.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdateTime:  User.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
