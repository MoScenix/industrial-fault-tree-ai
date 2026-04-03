package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/dal/mysql"
	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/dal/redis"
	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/model"
	user "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterService struct {
	ctx context.Context
} // NewRegisterService new RegisterService
func NewRegisterService(ctx context.Context) *RegisterService {
	return &RegisterService{ctx: ctx}
}

// Run create note info
func (s *RegisterService) Run(req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	// Finish your business logic.
	q := model.NewUserProQuery(s.ctx, mysql.DB, redis.RedisClient)
	PasswordHashed, err := bcrypt.GenerateFromPassword([]byte(req.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	newUser, err := q.CreateUser(model.User{
		UserAccount:  req.UserAccount,
		Name:         uuid.New().String(),
		UserRole:     "user",
		PasswordHash: string(PasswordHashed),
	})
	return &user.RegisterResp{
		UserId:   uint32(newUser.ID),
		UserRole: newUser.UserRole,
	}, nil
}
