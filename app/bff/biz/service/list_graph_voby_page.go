package service

import (
	"context"

	graph "github.com/MoScenix/industrial-fault-tree-ai/app/bff/hertz_gen/bff/graph"
	"github.com/MoScenix/industrial-fault-tree-ai/app/bff/infra/rpc"
	rpcgraph "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/graph"
	"github.com/cloudwego/hertz/pkg/app"
)

type ListGraphVOByPageService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListGraphVOByPageService(Context context.Context, RequestContext *app.RequestContext) *ListGraphVOByPageService {
	return &ListGraphVOByPageService{RequestContext: RequestContext, Context: Context}
}

func (h *ListGraphVOByPageService) Run(req *graph.GraphQueryRequest) (resp *graph.BaseResponsePageGraphVO, err error) {
	if err := ensureLogin(h.Context); err != nil {
		return &graph.BaseResponsePageGraphVO{Code: 1, Message: graphAccessError(err).Error()}, nil
	}
	userID := req.UserId
	if !isAdmin(h.Context) {
		currentUserID, _ := getCurrentUserID(h.Context)
		userID = currentUserID
	} else if req.UserId == 0 {
		userID = 0
	}
	res, err := rpc.GraphClient.ListGraph(h.Context, &rpcgraph.ListGraphReq{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		GraphName: req.GraphName,
		UserId:    userID,
	})
	if err != nil {
		return &graph.BaseResponsePageGraphVO{Code: 1, Message: err.Error()}, err
	}
	records := make([]*graph.GraphVO, 0, len(res.GraphList))
	for _, item := range res.GraphList {
		records = append(records, &graph.GraphVO{
			Id:             item.Id,
			GraphName:      item.GraphName,
			Description:    item.Description,
			Cover:          item.Cover,
			UserId:         item.UserId,
			CurrentVersion: item.CurrentVersion,
			HasTmp:         item.HasTmp,
			CreateTime:     item.CreateTime,
			UpdateTime:     item.UpdateTime,
		})
	}
	return &graph.BaseResponsePageGraphVO{
		Code:    0,
		Message: "success",
		Data:    makeGraphVOPage(records, req.PageNum, req.PageSize, res.Total),
	}, nil
}
