package promptutil

import (
	"os"
	"path/filepath"
	"time"

	ai "github.com/MoScenix/industrial-fault-tree-ai/rpc_gen/kitex_gen/ai"
)

const (
	modifyPromptPath = "prompt/modify/modify.txt"
	logPromptPath    = "prompt/log/log.txt"
)

func LoadPrompt(mode ai.PromptMode) (string, error) {
	path := PromptPath(mode)
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func SavePrompt(mode ai.PromptMode, content string) (string, error) {
	path := PromptPath(mode)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", err
	}
	return time.Now().Format("2006-01-02 15:04:05"), nil
}

func PromptPath(mode ai.PromptMode) string {
	switch mode {
	case ai.PromptMode_LOG_MODE:
		return logPromptPath
	default:
		return modifyPromptPath
	}
}
