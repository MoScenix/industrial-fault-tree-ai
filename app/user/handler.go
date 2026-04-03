package main

import (
	"context"

	"github.com/MoScenix/industrial-fault-tree-ai/app/user/biz/service"
	user "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Update implements the UserServiceImpl interface.
func (s *UserServiceImpl) Update(ctx context.Context, req *user.UpdateReq) (resp *user.UpdateResp, err error) {
	resp, err = service.NewUpdateService(ctx).Run(req)

	return resp, err
}

// GetUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserReq) (resp *user.GetUserResp, err error) {
	resp, err = service.NewGetUserService(ctx).Run(req)

	return resp, err
}

// DeleteUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) DeleteUser(ctx context.Context, req *user.DeleteUserReq) (resp *user.DeleteUserResp, err error) {
	resp, err = service.NewDeleteUserService(ctx).Run(req)

	return resp, err
}

// AddUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) AddUser(ctx context.Context, req *user.AddUserReq) (resp *user.AddUserResp, err error) {
	resp, err = service.NewAddUserService(ctx).Run(req)

	return resp, err
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	resp, err = service.NewLoginService(ctx).Run(req)

	return resp, err
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	resp, err = service.NewRegisterService(ctx).Run(req)

	return resp, err
}

// ListUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) ListUser(ctx context.Context, req *user.ListUserReq) (resp *user.ListUserResp, err error) {
	resp, err = service.NewListUserService(ctx).Run(req)

	return resp, err
}
