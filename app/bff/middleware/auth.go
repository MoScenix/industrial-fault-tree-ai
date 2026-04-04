package middleware

import (
	"context"
	"strings"

	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/biz/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/sessions"
)

func GlobalAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		session := sessions.Default(c)
		userId := session.Get(utils.UserIdKey)
		userRole := session.Get(utils.UserRoleKey)
		if userId == nil || userRole == nil {
			c.Next(ctx)
			return
		}
		ctx = context.WithValue(ctx, utils.UserIdKey, userId)
		ctx = context.WithValue(ctx, utils.UserRoleKey, userRole)
		c.Next(ctx)
	}
}

func Auth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		path := c.Path()
		if ctx.Value(utils.UserIdKey) == nil || ctx.Value(utils.UserRoleKey) != utils.AdminRole {
			if strings.Contains(string(path), "/admin") {
				c.Redirect(consts.StatusFound, []byte("/user/login"))
				c.Abort()
				return
			}
		}
		if ctx.Value(utils.UserIdKey) == nil {
			if strings.Contains(string(path), "/app") {
				c.Redirect(consts.StatusFound, []byte("/user/login"))
				c.Abort()
				return
			}
		}
		c.Next(ctx)
	}
}
