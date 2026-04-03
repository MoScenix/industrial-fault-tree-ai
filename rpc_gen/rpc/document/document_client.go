package document

import (
	"context"
	document "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/document"

	"github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/document/documentservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() documentservice.Client
	Service() string
	ParsePersonalPDF(ctx context.Context, Req *document.ParsePersonalPDFReq, callOptions ...callopt.Option) (r *document.ParsePDFResp, err error)
	ParseProjectPDF(ctx context.Context, Req *document.ParseProjectPDFReq, callOptions ...callopt.Option) (r *document.ParsePDFResp, err error)
	GetDocument(ctx context.Context, Req *document.GetDocumentReq, callOptions ...callopt.Option) (r *document.GetDocumentResp, err error)
	ListDocuments(ctx context.Context, Req *document.ListDocumentsReq, callOptions ...callopt.Option) (r *document.ListDocumentsResp, err error)
	SearchDocuments(ctx context.Context, Req *document.SearchDocumentsReq, callOptions ...callopt.Option) (r *document.SearchDocumentsResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := documentservice.NewClient(dstService, opts...)
	if err != nil {
		return nil, err
	}
	cli := &clientImpl{
		service:     dstService,
		kitexClient: kitexClient,
	}

	return cli, nil
}

type clientImpl struct {
	service     string
	kitexClient documentservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() documentservice.Client {
	return c.kitexClient
}

func (c *clientImpl) ParsePersonalPDF(ctx context.Context, Req *document.ParsePersonalPDFReq, callOptions ...callopt.Option) (r *document.ParsePDFResp, err error) {
	return c.kitexClient.ParsePersonalPDF(ctx, Req, callOptions...)
}

func (c *clientImpl) ParseProjectPDF(ctx context.Context, Req *document.ParseProjectPDFReq, callOptions ...callopt.Option) (r *document.ParsePDFResp, err error) {
	return c.kitexClient.ParseProjectPDF(ctx, Req, callOptions...)
}

func (c *clientImpl) GetDocument(ctx context.Context, Req *document.GetDocumentReq, callOptions ...callopt.Option) (r *document.GetDocumentResp, err error) {
	return c.kitexClient.GetDocument(ctx, Req, callOptions...)
}

func (c *clientImpl) ListDocuments(ctx context.Context, Req *document.ListDocumentsReq, callOptions ...callopt.Option) (r *document.ListDocumentsResp, err error) {
	return c.kitexClient.ListDocuments(ctx, Req, callOptions...)
}

func (c *clientImpl) SearchDocuments(ctx context.Context, Req *document.SearchDocumentsReq, callOptions ...callopt.Option) (r *document.SearchDocumentsResp, err error) {
	return c.kitexClient.SearchDocuments(ctx, Req, callOptions...)
}
