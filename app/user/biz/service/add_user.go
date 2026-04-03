package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/dal/mysql"
	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/dal/redis"
	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/model"
	user "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/user"
	"golang.org/x/crypto/bcrypt"
)

type AddUserService struct {
	ctx context.Context
} // NewAddUserService new AddUserService
func NewAddUserService(ctx context.Context) *AddUserService {
	return &AddUserService{ctx: ctx}
}

// Run create note info
func (s *AddUserService) Run(req *user.AddUserReq) (resp *user.AddUserResp, err error) {
	// Finish your business logic.
	q := model.NewUserProQuery(s.ctx, mysql.DB, redis.RedisClient)
	PasswordHashed, err := bcrypt.GenerateFromPassword([]byte(req.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	_, err = q.CreateUser(model.User{
		UserAccount:  req.UserAccount,
		Name:         req.UserName,
		UserRole:     req.UserRole,
		PasswordHash: string(PasswordHashed),
		UserAvatar:   req.UserAvatar,
		UserProfile:  req.UserProfile,
	})
	return &user.AddUserResp{}, nil
}
