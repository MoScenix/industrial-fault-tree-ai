package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/dal/mysql"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/dal/redis"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/model"
	"github.com/bytedance/gopkg/cloud/metainfo"
)

const (
	metaUserIDKey   = "x-user-id"
	metaUserRoleKey = "x-user-role"
	adminRole       = "admin"
)

var (
	errDBNotReady   = errors.New("database not initialized")
	errUnauthorized = errors.New("unauthorized")
	errForbidden    = errors.New("forbidden")
)

func currentUserID(ctx context.Context) (int64, bool) {
	userID, ok := metainfo.GetPersistentValue(ctx, metaUserIDKey)
	if !ok || userID == "" {
		return 0, false
	}
	parsed, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return 0, false
	}
	return parsed, true
}

func currentUserRole(ctx context.Context) string {
	role, _ := metainfo.GetPersistentValue(ctx, metaUserRoleKey)
	return role
}

func isAdmin(ctx context.Context) bool {
	return currentUserRole(ctx) == adminRole
}

func ensureLogin(ctx context.Context) (int64, error) {
	userID, ok := currentUserID(ctx)
	if !ok || userID <= 0 {
		return 0, errUnauthorized
	}
	return userID, nil
}

func graphQuery(ctx context.Context) (*model.GraphProQuery, error) {
	if mysql.DB == nil {
		return nil, errDBNotReady
	}
	return model.NewGraphProQuery(ctx, mysql.DB, redis.RedisClient), nil
}

func mustLoadAuthorizedGraph(ctx context.Context, graphID int64) (model.Graph, error) {
	userID, err := ensureLogin(ctx)
	if err != nil {
		return model.Graph{}, err
	}
	q, err := graphQuery(ctx)
	if err != nil {
		return model.Graph{}, err
	}
	item, err := q.GetGraphByID(uint(graphID))
	if err != nil {
		return model.Graph{}, err
	}
	if isAdmin(ctx) || item.UserID == uint(userID) {
		return item, nil
	}
	return model.Graph{}, errForbidden
}

func effectiveListUserID(ctx context.Context, requested int64) (uint, error) {
	userID, err := ensureLogin(ctx)
	if err != nil {
		return 0, err
	}
	if isAdmin(ctx) {
		if requested <= 0 {
			return 0, nil
		}
		return uint(requested), nil
	}
	return uint(userID), nil
}
