package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name         string `gorm:"type:varchar(50);not null;index"`
	PasswordHash string `gorm:"type:varchar(100);not null"`
	UserAccount  string `gorm:"type:varchar(50);not null;uniqueIndex"`

	UserAvatar  string `gorm:"type:varchar(100);not null;default:''"`
	UserProfile string `gorm:"type:varchar(100);default:''"`
	UserRole    string `gorm:"type:varchar(10);not null;default:'user'"`
}

func (User) TableName() string {
	return "user"
}

type UserQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewUserQuery(ctx context.Context, db *gorm.DB) *UserQuery {
	return &UserQuery{
		ctx: ctx,
		db:  db,
	}
}
func (q *UserQuery) GetUserById(id int) (User, error) {
	user := User{}
	err := q.db.WithContext(q.ctx).Model(&User{}).Where("id = ?", id).First(&user).Error
	return user, err
}
func (q *UserQuery) GetUserByAccount(account string) (User, error) {
	user := User{}
	err := q.db.WithContext(q.ctx).Model(&User{}).Where("user_account = ?", account).First(&user).Error
	return user, err
}
func (q *UserQuery) UpdateUser(id uint, user User) error {
	err := q.db.WithContext(q.ctx).Model(&User{}).Where("id = ?", id).Updates(user).Error
	return err
}
func (q *UserQuery) CreateUser(user User) (User, error) {
	err := q.db.WithContext(q.ctx).Model(&User{}).Create(&user).Error
	return user, err
}
func (q *UserQuery) DeleteUser(id int) error {
	err := q.db.WithContext(q.ctx).Model(&User{}).Where("id = ?", id).Delete(&User{}).Error
	return err
}
func (q *UserQuery) ListUser(page uint32, username, useraccount string, page_size uint32) ([]User, error) {
	var users []User
	err := q.db.WithContext(q.ctx).Model(&User{}).Where("name like ? and user_account like ?", username+"%", useraccount+"%").Limit(int(page_size)).Offset(int(page_size * (page - 1))).Find(&users).Error
	return users, err
}
func (q *UserQuery) CountUser(username, useraccount string) (int64, error) {
	var count int64
	err := q.db.WithContext(q.ctx).Model(&User{}).Where("name like ? and user_account like ?", username+"%", useraccount+"%").Count(&count).Error
	return count, err
}

type UserProQuery struct {
	q      *UserQuery
	rdb    *redis.Client
	prefix string
}

func NewUserProQuery(ctx context.Context, db *gorm.DB, rdb *redis.Client) *UserProQuery {
	return &UserProQuery{
		q:      NewUserQuery(ctx, db),
		rdb:    rdb,
		prefix: "industrial-fault-tree-ai",
	}
}

func (p *UserProQuery) keyUser(id int) string {
	return fmt.Sprintf("%s_user_%d", p.prefix, id)
}

func (p *UserProQuery) keyAccount(account string) string {
	return fmt.Sprintf("%s_user_account_%s", p.prefix, account)
}
func (p *UserProQuery) GetUserById(id int) (User, error) {
	if p.rdb != nil {
		if val, err := p.rdb.Get(p.q.ctx, p.keyUser(id)).Result(); err == nil && val != "" {
			var u User
			if json.Unmarshal([]byte(val), &u) == nil {
				return u, nil
			}
		}
	}

	u, err := p.q.GetUserById(id)
	if err != nil {
		return User{}, err
	}

	if p.rdb != nil {
		if b, e := json.Marshal(u); e == nil {
			_ = p.rdb.Set(p.q.ctx, p.keyUser(id), b, time.Hour).Err()
			if u.UserAccount != "" {
				_ = p.rdb.Set(p.q.ctx, p.keyAccount(u.UserAccount), b, time.Hour).Err()
			}
		}
	}
	return u, nil
}
func (p *UserProQuery) GetUserByAccount(account string) (User, error) {
	if p.rdb != nil {
		if val, err := p.rdb.Get(p.q.ctx, p.keyAccount(account)).Result(); err == nil && val != "" {
			var u User
			if json.Unmarshal([]byte(val), &u) == nil {
				return u, nil
			}
		}
	}

	u, err := p.q.GetUserByAccount(account)
	if err != nil {
		return User{}, err
	}

	if p.rdb != nil {
		if b, e := json.Marshal(u); e == nil {
			_ = p.rdb.Set(p.q.ctx, p.keyAccount(account), b, time.Hour).Err()
			_ = p.rdb.Set(p.q.ctx, p.keyUser(int(u.ID)), b, time.Hour).Err()
		}
	}
	return u, nil
}

func (p *UserProQuery) UpdateUser(id uint, user User) error {
	if p.rdb != nil {
		var oldAcc string
		if old, err := p.q.GetUserById(int(id)); err == nil {
			oldAcc = old.UserAccount
		}

		_ = p.rdb.Del(p.q.ctx, p.keyUser(int(id))).Err()
		if oldAcc != "" {
			_ = p.rdb.Del(p.q.ctx, p.keyAccount(oldAcc)).Err()
		}
		if user.UserAccount != "" && user.UserAccount != oldAcc {
			_ = p.rdb.Del(p.q.ctx, p.keyAccount(user.UserAccount)).Err()
		}
	}
	return p.q.UpdateUser(id, user)
}
func (p *UserProQuery) CreateUser(user User) (User, error) {
	if p.rdb != nil && user.UserAccount != "" {
		_ = p.rdb.Del(p.q.ctx, p.keyAccount(user.UserAccount)).Err()
	}
	return p.q.CreateUser(user)
}

func (p *UserProQuery) DeleteUser(id int) error {
	if p.rdb != nil {
		var acc string
		if u, err := p.q.GetUserById(id); err == nil {
			acc = u.UserAccount
		}
		_ = p.rdb.Del(p.q.ctx, p.keyUser(id)).Err()
		if acc != "" {
			_ = p.rdb.Del(p.q.ctx, p.keyAccount(acc)).Err()
		}
	}
	return p.q.DeleteUser(id)
}
func (p *UserProQuery) ListUser(page uint32, username, useraccount string, pageSize uint32) ([]User, error) {
	return p.q.ListUser(page, username, useraccount, pageSize)
}
func (p *UserProQuery) CountUser(username, useraccount string) (int64, error) {
	return p.q.CountUser(username, useraccount)
}
