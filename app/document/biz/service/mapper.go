package service

import (
	"strings"

	"github.com/MoScenix/industrial-fault-tree-ai/app/document/biz/model"
	document "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/document"
)

func toProtoDocument(doc model.Document) *document.Document {
	chunks := make([]*document.DocumentChunk, 0, len(doc.Chunks))
	for _, chunk := range doc.Chunks {
		chunks = append(chunks, &document.DocumentChunk{
			ChunkId: chunk.ChunkID,
			Text:    chunk.Text,
			Page:    chunk.Page,
			Order:   chunk.Order,
		})
	}

	return &document.Document{
		DocumentId:  doc.DocumentID,
		OwnerType:   toProtoOwnerType(doc.OwnerType),
		OwnerId:     doc.OwnerID,
		PdfId:       doc.PdfID,
		FileName:    doc.FileName,
		DisplayName: doc.DisplayName,
		ParseStatus: toProtoParseStatus(doc.ParseStatus),
		Summary:     doc.Summary,
		Chunks:      chunks,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}
}

func toProtoOwnerType(ownerType string) document.OwnerType {
	switch strings.ToUpper(ownerType) {
	case "PERSONAL":
		return document.OwnerType_PERSONAL
	case "PROJECT":
		return document.OwnerType_PROJECT
	default:
		return document.OwnerType_OWNER_TYPE_UNSPECIFIED
	}
}

func toProtoParseStatus(status string) document.ParseStatus {
	switch strings.ToUpper(status) {
	case "SUCCESS":
		return document.ParseStatus_SUCCESS
	case "FAILED":
		return document.ParseStatus_FAILED
	default:
		return document.ParseStatus_PARSE_STATUS_UNSPECIFIED
	}
}
