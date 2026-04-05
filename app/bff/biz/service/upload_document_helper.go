package service

import (
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

func sanitizeUploadFileName(fileHeader *multipart.FileHeader, fallback string) string {
	if fileHeader == nil {
		return fallback
	}
	name := filepath.Base(strings.TrimSpace(fileHeader.Filename))
	if name == "." || name == string(filepath.Separator) || name == "" {
		return fallback
	}
	return name
}

func saveUploadedFileToDocDir(c *app.RequestContext, fileHeader *multipart.FileHeader, rootDir, fileName string) (string, error) {
	if err := os.MkdirAll(rootDir, 0o755); err != nil {
		return "", err
	}
	dstPath := filepath.Join(rootDir, fileName)
	if err := c.SaveUploadedFile(fileHeader, dstPath); err != nil {
		return "", err
	}
	return dstPath, nil
}
