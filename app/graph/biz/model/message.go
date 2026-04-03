package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	GraphID uint   `gorm:"type:int;index" json:"graphId"`
	UserID  uint   `gorm:"type:int;index" json:"userId"`
	Role    string `gorm:"type:varchar(50)" json:"role"`
	Content string `gorm:"type:text" json:"content"`
}

type MessageQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewMessageQuery(ctx context.Context, db *gorm.DB) *MessageQuery {
	return &MessageQuery{
		ctx: ctx,
		db:  db,
	}
}

func (q *MessageQuery) CreateMessage(msg Message) (Message, error) {
	err := q.db.WithContext(q.ctx).
		Model(&Message{}).
		Create(&msg).Error
	return msg, err
}

func (q *MessageQuery) GetMessageByID(id uint) (Message, error) {
	msg := Message{}
	err := q.db.WithContext(q.ctx).
		Model(&Message{}).
		Where("id = ?", id).
		First(&msg).Error
	return msg, err
}

func (q *MessageQuery) DeleteMessageByID(id uint) error {
	return q.db.WithContext(q.ctx).
		Model(&Message{}).
		Where("id = ?", id).
		Delete(&Message{}).Error
}

func (q *MessageQuery) ListMessagesByGraphID(graphID uint, limit int, lastCreateTime *time.Time) ([]Message, error) {
	var msgs []Message

	tx := q.db.WithContext(q.ctx).
		Model(&Message{}).
		Where("graph_id = ?", graphID)

	if lastCreateTime != nil && !lastCreateTime.IsZero() {
		tx = tx.Where("created_at < ?", *lastCreateTime)
	}

	err := tx.Order("created_at desc").
		Limit(limit).
		Find(&msgs).Error

	return msgs, err
}

func (q *MessageQuery) Count(graphID uint) (int64, error) {
	var count int64
	err := q.db.WithContext(q.ctx).
		Model(&Message{}).
		Where("graph_id = ?", graphID).
		Count(&count).Error
	return count, err
}

type MessageProQuery struct {
	q      *MessageQuery
	rdb    *redis.Client
	prefix string
}

func NewMessageProQuery(ctx context.Context, db *gorm.DB, rdb *redis.Client) *MessageProQuery {
	return &MessageProQuery{
		q:      NewMessageQuery(ctx, db),
		rdb:    rdb,
		prefix: "industrial-fault-tree-ai",
	}
}

func (p *MessageProQuery) keyGraphMessages(graphID uint, limit int, lastCreateTime *time.Time) string {
	cursor := ""
	if lastCreateTime != nil && !lastCreateTime.IsZero() {
		cursor = lastCreateTime.Format(time.RFC3339Nano)
	}
	return fmt.Sprintf("%s_graph_messages_%d_%d_%s", p.prefix, graphID, limit, cursor)
}

func (p *MessageProQuery) invalidateGraphMessages(graphID uint) {
	if p.rdb == nil {
		return
	}
	pattern := fmt.Sprintf("%s_graph_messages_%d_*", p.prefix, graphID)
	keys, err := p.rdb.Keys(p.q.ctx, pattern).Result()
	if err != nil || len(keys) == 0 {
		return
	}
	_ = p.rdb.Del(p.q.ctx, keys...).Err()
}

func (p *MessageProQuery) CreateMessage(msg Message) (Message, error) {
	created, err := p.q.CreateMessage(msg)
	if err != nil {
		return Message{}, err
	}
	p.invalidateGraphMessages(created.GraphID)
	return created, nil
}

func (p *MessageProQuery) GetMessageByID(id uint) (Message, error) {
	return p.q.GetMessageByID(id)
}

func (p *MessageProQuery) DeleteMessageByID(id uint) error {
	msg, err := p.q.GetMessageByID(id)
	if err == nil {
		p.invalidateGraphMessages(msg.GraphID)
	}
	return p.q.DeleteMessageByID(id)
}

func (p *MessageProQuery) ListMessagesByGraphID(graphID uint, limit int, lastCreateTime *time.Time) ([]Message, error) {
	if p.rdb != nil {
		if val, err := p.rdb.Get(p.q.ctx, p.keyGraphMessages(graphID, limit, lastCreateTime)).Result(); err == nil && val != "" {
			var msgs []Message
			if json.Unmarshal([]byte(val), &msgs) == nil {
				return msgs, nil
			}
		}
	}

	msgs, err := p.q.ListMessagesByGraphID(graphID, limit, lastCreateTime)
	if err != nil {
		return nil, err
	}

	if p.rdb != nil {
		if b, e := json.Marshal(msgs); e == nil {
			_ = p.rdb.Set(p.q.ctx, p.keyGraphMessages(graphID, limit, lastCreateTime), b, time.Hour).Err()
		}
	}
	return msgs, nil
}

func (p *MessageProQuery) Count(graphID uint) (int64, error) {
	return p.q.Count(graphID)
}
