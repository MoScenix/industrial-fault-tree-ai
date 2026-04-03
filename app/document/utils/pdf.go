package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MoScenix/industrial-fault-tree-ai/app/document/biz/model"
	"github.com/MoScenix/industrial-fault-tree-ai/app/document/conf"
	"github.com/ledongthuc/pdf"
)

const (
	defaultChunkSize    = 1000
	defaultChunkOverlap = 200
)

func ResolvePDFPath(pdfID string) string {
	return filepath.Join(conf.GetConf().Document.PDFDir, fmt.Sprintf("%s.pdf", pdfID))
}

func ParsePDFFile(pdfID string) (string, error) {
	pdfPath := ResolvePDFPath(pdfID)
	if _, err := os.Stat(pdfPath); err != nil {
		return "", err
	}

	f, r, err := pdf.Open(pdfPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	plainText, err := r.GetPlainText()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(plainText); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

func BuildChunks(documentID, content string) []model.DocumentChunk {
	if content == "" {
		return nil
	}

	runes := []rune(content)
	chunks := make([]model.DocumentChunk, 0, len(runes)/defaultChunkSize+1)
	start := 0
	order := int64(1)

	for start < len(runes) {
		end := start + defaultChunkSize
		if end > len(runes) {
			end = len(runes)
		}

		chunkText := strings.TrimSpace(string(runes[start:end]))
		if chunkText != "" {
			chunks = append(chunks, model.DocumentChunk{
				ChunkID:    fmt.Sprintf("%s_%03d", documentID, order),
				DocumentID: documentID,
				Text:       chunkText,
				Page:       0,
				Order:      order,
			})
			order++
		}

		if end == len(runes) {
			break
		}

		start = end - defaultChunkOverlap
		if start < 0 {
			start = 0
		}
	}

	return chunks
}
