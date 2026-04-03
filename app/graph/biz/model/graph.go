package model

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/conf"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Graph struct {
	gorm.Model
	GraphName   string `gorm:"type:varchar(100);index" json:"graphName"`
	Description string `gorm:"type:text" json:"description"`
	Cover       string `gorm:"type:varchar(255)" json:"cover"`
	UserID      uint   `gorm:"type:int;index" json:"userId"`
	ProjectUUID string `gorm:"type:varchar(64);uniqueIndex" json:"projectUUID"`
	ProjectDir  string `gorm:"type:varchar(255)" json:"projectDir"`
}

type GraphQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewGraphQuery(ctx context.Context, db *gorm.DB) *GraphQuery {
	return &GraphQuery{
		ctx: ctx,
		db:  db,
	}
}

func (q *GraphQuery) CreateGraph(graph Graph) (Graph, error) {
	if graph.ProjectUUID != "" && graph.ProjectDir == "" {
		graph.ProjectDir = BuildProjectDir(graph.ProjectUUID)
	}

	err := q.db.WithContext(q.ctx).
		Model(&Graph{}).
		Create(&graph).Error
	return graph, err
}

func (q *GraphQuery) GetGraphByID(id uint) (Graph, error) {
	graph := Graph{}
	err := q.db.WithContext(q.ctx).
		Model(&Graph{}).
		Where("id = ?", id).
		First(&graph).Error
	return graph, err
}

func (q *GraphQuery) UpdateGraph(id uint, graph Graph) error {
	if graph.ProjectUUID != "" && graph.ProjectDir == "" {
		graph.ProjectDir = BuildProjectDir(graph.ProjectUUID)
	}

	return q.db.WithContext(q.ctx).
		Model(&Graph{}).
		Where("id = ?", id).
		Updates(graph).Error
}

func (q *GraphQuery) DeleteGraph(id uint) error {
	return q.db.WithContext(q.ctx).
		Model(&Graph{}).
		Where("id = ?", id).
		Delete(&Graph{}).Error
}

func (q *GraphQuery) ListGraph(page uint32, userID uint, graphName string, pageSize uint32) ([]Graph, error) {
	var graphs []Graph

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	tx := q.db.WithContext(q.ctx).Model(&Graph{})
	if userID != 0 {
		tx = tx.Where("user_id = ?", userID)
	}
	if graphName != "" {
		tx = tx.Where("graph_name LIKE ?", graphName+"%")
	}

	err := tx.Order("id desc").
		Limit(int(pageSize)).
		Offset(int(pageSize * (page - 1))).
		Find(&graphs).Error

	return graphs, err
}

func (q *GraphQuery) CountGraph(userID uint, graphName string) (int64, error) {
	var count int64

	tx := q.db.WithContext(q.ctx).Model(&Graph{})
	if userID != 0 {
		tx = tx.Where("user_id = ?", userID)
	}
	if graphName != "" {
		tx = tx.Where("graph_name LIKE ?", graphName+"%")
	}

	err := tx.Count(&count).Error
	return count, err
}

type GraphProQuery struct {
	q      *GraphQuery
	rdb    *redis.Client
	prefix string
}

func NewGraphProQuery(ctx context.Context, db *gorm.DB, rdb *redis.Client) *GraphProQuery {
	return &GraphProQuery{
		q:      NewGraphQuery(ctx, db),
		rdb:    rdb,
		prefix: "industrial-fault-tree-ai",
	}
}

func (p *GraphProQuery) keyGraph(id uint) string {
	return fmt.Sprintf("%s_graph_%d", p.prefix, id)
}

func (p *GraphProQuery) GetGraphByID(id uint) (Graph, error) {
	if p.rdb != nil {
		if val, err := p.rdb.Get(p.q.ctx, p.keyGraph(id)).Result(); err == nil && val != "" {
			var item Graph
			if json.Unmarshal([]byte(val), &item) == nil {
				return item, nil
			}
		}
	}

	item, err := p.q.GetGraphByID(id)
	if err != nil {
		return Graph{}, err
	}

	if p.rdb != nil {
		if b, e := json.Marshal(item); e == nil {
			_ = p.rdb.Set(p.q.ctx, p.keyGraph(id), b, time.Hour).Err()
		}
	}
	return item, nil
}

func (p *GraphProQuery) CreateGraph(graph Graph) (Graph, error) {
	created, err := p.q.CreateGraph(graph)
	if err != nil {
		return Graph{}, err
	}
	if p.rdb != nil {
		_ = p.rdb.Del(p.q.ctx, p.keyGraph(created.ID)).Err()
	}
	return created, nil
}

func (p *GraphProQuery) UpdateGraph(id uint, graph Graph) error {
	if p.rdb != nil {
		_ = p.rdb.Del(p.q.ctx, p.keyGraph(id)).Err()
	}
	return p.q.UpdateGraph(id, graph)
}

func (p *GraphProQuery) DeleteGraph(id uint) error {
	if p.rdb != nil {
		_ = p.rdb.Del(p.q.ctx, p.keyGraph(id)).Err()
	}
	return p.q.DeleteGraph(id)
}

func (p *GraphProQuery) ListGraph(page uint32, userID uint, graphName string, pageSize uint32) ([]Graph, error) {
	return p.q.ListGraph(page, userID, graphName, pageSize)
}

func (p *GraphProQuery) CountGraph(userID uint, graphName string) (int64, error) {
	return p.q.CountGraph(userID, graphName)
}

func BuildProjectDir(projectUUID string) string {
	return filepath.Join(conf.GetConf().Graph.RootDir, projectUUID)
}
