package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MoScenix/industrial-fault-tree-ai/app/document/biz/model"
	"github.com/MoScenix/industrial-fault-tree-ai/app/document/conf"
	"github.com/ledongthuc/pdf"
)

const (
	defaultChunkSize    = 1000
	defaultChunkOverlap = 200
)

func ResolvePDFPath(pdfID, fileName string) (string, error) {
	docDir := filepath.Join(conf.GetConf().Document.PDFDir, pdfID)
	if fileName != "" {
		candidate := filepath.Join(docDir, filepath.Base(fileName))
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	entries, err := os.ReadDir(docDir)
	if err != nil {
		return "", err
	}

	candidates := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		candidates = append(candidates, filepath.Join(docDir, entry.Name()))
	}
	sort.Strings(candidates)
	if len(candidates) == 0 {
		return "", os.ErrNotExist
	}
	return candidates[0], nil
}

func ParsePDFFile(pdfID, fileName string) (string, error) {
	pdfPath, err := ResolvePDFPath(pdfID, fileName)
	if err != nil {
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
