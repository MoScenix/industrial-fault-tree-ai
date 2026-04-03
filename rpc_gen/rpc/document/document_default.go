package document

import (
	"context"
	document "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/document"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func ParsePersonalPDF(ctx context.Context, req *document.ParsePersonalPDFReq, callOptions ...callopt.Option) (resp *document.ParsePDFResp, err error) {
	resp, err = defaultClient.ParsePersonalPDF(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "ParsePersonalPDF call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func ParseProjectPDF(ctx context.Context, req *document.ParseProjectPDFReq, callOptions ...callopt.Option) (resp *document.ParsePDFResp, err error) {
	resp, err = defaultClient.ParseProjectPDF(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "ParseProjectPDF call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func GetDocument(ctx context.Context, req *document.GetDocumentReq, callOptions ...callopt.Option) (resp *document.GetDocumentResp, err error) {
	resp, err = defaultClient.GetDocument(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "GetDocument call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func ListDocuments(ctx context.Context, req *document.ListDocumentsReq, callOptions ...callopt.Option) (resp *document.ListDocumentsResp, err error) {
	resp, err = defaultClient.ListDocuments(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "ListDocuments call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func SearchDocuments(ctx context.Context, req *document.SearchDocumentsReq, callOptions ...callopt.Option) (resp *document.SearchDocumentsResp, err error) {
	resp, err = defaultClient.SearchDocuments(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "SearchDocuments call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
