package service

import (
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/biz/model"
	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/utils"
	graphpb "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
)

func toGraphInfo(item model.Graph) *graphpb.GraphInfo {
	currentVersion, _ := utils.ReadCurrentVersion(item.ProjectDir)
	return &graphpb.GraphInfo{
		Id:             int64(item.ID),
		GraphName:      item.GraphName,
		Description:    item.Description,
		Cover:          item.Cover,
		UserId:         int64(item.UserID),
		CurrentVersion: currentVersion,
		HasTmp:         utils.HasTmp(item.ProjectDir),
		CreateTime:     utils.FormatTime(item.CreatedAt),
		UpdateTime:     utils.FormatTime(item.UpdatedAt),
		ProjectDir:     item.ProjectDir,
	}
}

func toGraphMessage(item model.Message) *graphpb.GraphMessage {
	return &graphpb.GraphMessage{
		Id:         int64(item.ID),
		GraphId:    int64(item.GraphID),
		UserId:     int64(item.UserID),
		Role:       item.Role,
		Content:    item.Content,
		CreateTime: utils.FormatTime(item.CreatedAt),
		UpdateTime: utils.FormatTime(item.UpdatedAt),
		IsDelete:   0,
	}
}
