package service

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/dal/mysql"
	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/dal/redis"
	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/model"
	user "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/user"
)

type ListUserService struct {
	ctx context.Context
} // NewListUserService new ListUserService
func NewListUserService(ctx context.Context) *ListUserService {
	return &ListUserService{ctx: ctx}
}

// Run create note info
func (s *ListUserService) Run(req *user.ListUserReq) (resp *user.ListUserResp, err error) {
	// Finish your business logic.
	q := model.NewUserProQuery(s.ctx, mysql.DB, redis.RedisClient)
	res, err := q.ListUser(uint32(req.PageNum), req.UserName, req.Account, uint32(req.PageSize))
	if err != nil {
		return nil, err
	}
	total, err := q.CountUser(req.UserName, req.Account)
	if err != nil {
		return nil, err
	}
	resp = &user.ListUserResp{
		Total: int64((total + req.PageSize - 1) / req.PageSize),
	}
	for _, v := range res {
		resp.UserList = append(resp.UserList, &user.GetUserResp{
			Id:           int64(v.ID),
			UserAccount:  v.UserAccount,
			UserAvatar:   v.UserAvatar,
			UserProfile:  v.UserProfile,
			UserRole:     v.UserRole,
			CreateTime:   v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdateTime:   v.UpdatedAt.Format("2006-01-02 15:04:05"),
			UserName:     v.Name,
			IsDelete:     v.DeletedAt.Valid,
			UserPassword: v.PasswordHash,
		})
	}
	return resp, nil
}
