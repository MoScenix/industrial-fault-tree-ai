package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/biz/dal/mysql"
	lutils "github.com/MoScenix/industrial-fault-tree-ai/app/bff/biz/utils"
	graphbff "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
)

var errGraphNotFound = errors.New("graph not found")
var errUnauthorizedGraphAccess = errors.New("no permission to access this graph")
var errUnauthorized = errors.New("unauthorized")

type graphRecord struct {
	ID          uint      `gorm:"column:id"`
	GraphName   string    `gorm:"column:graph_name"`
	Description string    `gorm:"column:description"`
	Cover       string    `gorm:"column:cover"`
	UserID      uint      `gorm:"column:user_id"`
	ProjectUUID string    `gorm:"column:project_uuid"`
	ProjectDir  string    `gorm:"column:project_dir"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (graphRecord) TableName() string {
	return "graphs"
}

func loadGraphRecord(id int64) (*graphRecord, error) {
	if mysql.DB == nil {
		return nil, errDBNotReady
	}
	var item graphRecord
	if err := mysql.DB.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func currentUserIDFromContext(ctx any) (int64, bool) {
	switch v := ctx.(type) {
	case int64:
		return v, true
	case int32:
		return int64(v), true
	case int:
		return int64(v), true
	case float64:
		return int64(v), true
	case float32:
		return int64(v), true
	case uint64:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint:
		return int64(v), true
	default:
		return 0, false
	}
}

func getCurrentUserID(ctx interface{ Value(key any) any }) (int64, bool) {
	return currentUserIDFromContext(ctx.Value(lutils.UserIdKey))
}

func getCurrentUserRole(ctx interface{ Value(key any) any }) string {
	role, _ := ctx.Value(lutils.UserRoleKey).(string)
	return role
}

func isAdmin(ctx interface{ Value(key any) any }) bool {
	return getCurrentUserRole(ctx) == lutils.AdminRole
}

func ensureLogin(ctx interface{ Value(key any) any }) error {
	if _, ok := getCurrentUserID(ctx); !ok {
		return errUnauthorized
	}
	return nil
}

func loadAuthorizedGraphRecord(ctx interface{ Value(key any) any }, graphID int64) (*graphRecord, error) {
	if err := ensureLogin(ctx); err != nil {
		return nil, err
	}
	item, err := loadGraphRecord(graphID)
	if err != nil {
		return nil, err
	}
	if isAdmin(ctx) {
		return item, nil
	}
	userID, _ := getCurrentUserID(ctx)
	if item.UserID != uint(userID) {
		return nil, errUnauthorizedGraphAccess
	}
	return item, nil
}

func graphAccessError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, errUnauthorized):
		return fmt.Errorf("请先登录")
	case errors.Is(err, errUnauthorizedGraphAccess):
		return fmt.Errorf("无权操作该项目")
	default:
		return err
	}
}

func currentVersion(projectDir string) string {
	content, err := os.ReadFile(filepath.Join(projectDir, "current"))
	if err != nil {
		return ""
	}
	return string(content)
}

func hasTmp(projectDir string) bool {
	entries, err := os.ReadDir(filepath.Join(projectDir, "tmp"))
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if entry.IsDir() {
			return true
		}
	}
	return false
}

func tmpVersionDir(projectDir, version string) string {
	return filepath.Join(projectDir, "tmp", version)
}

func versionDir(projectDir, version string) string {
	return filepath.Join(projectDir, "versions", version)
}

func treePath(projectDir, version string, isTmp bool) string {
	if isTmp {
		return filepath.Join(tmpVersionDir(projectDir, version), "tree.json")
	}
	return filepath.Join(versionDir(projectDir, version), "tree.json")
}

func suggestionPath(projectDir, version string) string {
	return filepath.Join(projectDir, "suggestions", version+".md")
}

func readOptionalFile(path string) (string, time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", time.Time{}, nil
		}
		return "", time.Time{}, err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return "", time.Time{}, err
	}
	return string(content), info.ModTime(), nil
}

func toGraphVO(item *graphRecord) *graphbff.GraphVO {
	return &graphbff.GraphVO{
		Id:             int64(item.ID),
		GraphName:      item.GraphName,
		Description:    item.Description,
		Cover:          item.Cover,
		UserId:         int64(item.UserID),
		CurrentVersion: currentVersion(item.ProjectDir),
		HasTmp:         hasTmp(item.ProjectDir),
		CreateTime:     item.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdateTime:     item.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func makeGraphVOPage(records []*graphbff.GraphVO, pageNum, pageSize, total int64) *graphbff.PageGraphVO {
	totalPage := int64(0)
	if pageSize > 0 {
		totalPage = int64(math.Ceil(float64(total) / float64(pageSize)))
	}
	return &graphbff.PageGraphVO{
		Records:    records,
		PageNumber: pageNum,
		PageSize:   pageSize,
		TotalPage:  totalPage,
		TotalRow:   total,
	}
}

func makeGraphVersionPage(records []*graphbff.GraphVersionVO, total int64) *graphbff.PageGraphVersionVO {
	return &graphbff.PageGraphVersionVO{
		Records:    records,
		PageNumber: 1,
		PageSize:   total,
		TotalPage:  1,
		TotalRow:   total,
	}
}

func makeGraphMessagePage(records []*graphbff.GraphMessageVO, pageSize, total int64) *graphbff.PageGraphMessageVO {
	return &graphbff.PageGraphMessageVO{
		Records:    records,
		PageNumber: 1,
		PageSize:   pageSize,
		TotalPage:  1,
		TotalRow:   total,
	}
}

func normalizeJSONString(content string) string {
	var obj interface{}
	if err := json.Unmarshal([]byte(content), &obj); err != nil {
		return content
	}
	pretty, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return content
	}
	return string(pretty)
}

func newObjectID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return time.Now().Format("20060102150405")
	}
	return hex.EncodeToString(buf)
}
