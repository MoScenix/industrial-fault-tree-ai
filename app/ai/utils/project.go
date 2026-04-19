package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultGraphRootDir = "/graph"
	treeFileName        = "tree.json"
)

func GraphRootDir() string {
	if root := strings.TrimSpace(os.Getenv("GRAPH_ROOT_DIR")); root != "" {
		return root
	}
	return defaultGraphRootDir
}

func ProjectDir(projectID string) string {
	return filepath.Join(GraphRootDir(), projectID)
}

func TmpTreePath(projectID string) string {
	currentVersion, _ := ReadCurrentVersion(projectID)
	return TmpVersionTreePath(projectID, currentVersion)
}

func TmpVersionTreePath(projectID, version string) string {
	return filepath.Join(ProjectDir(projectID), "tmp", version, treeFileName)
}

func VersionTreePath(projectID, version string) string {
	return filepath.Join(ProjectDir(projectID), "versions", version, treeFileName)
}

func SuggestionPath(projectID, version string) string {
	return filepath.Join(ProjectDir(projectID), "suggestions", version+".md")
}

func CurrentVersionPath(projectID string) string {
	return filepath.Join(ProjectDir(projectID), "current")
}

func ReadCurrentVersion(projectID string) (string, error) {
	content, err := os.ReadFile(CurrentVersionPath(projectID))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func LoadGraphFile(path string) (*GraphFile, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var graph GraphFile
	if err := json.Unmarshal(content, &graph); err != nil {
		return nil, err
	}
	return &graph, nil
}

func SaveGraphFile(path string, graph *GraphFile) error {
	if graph == nil {
		return errors.New("graph is nil")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	content, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o644)
}

func GraphFileLines(graph *GraphFile) ([]string, error) {
	if graph == nil {
		return nil, errors.New("graph is nil")
	}
	content, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return nil, err
	}
	return splitLines(string(content)), nil
}

func LoadTextLines(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return splitLines(string(content)), nil
}

func SaveTextLines(path string, lines []string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(joinLines(lines)), 0o644)
}

func DefaultGraphLines(projectID string) ([]string, error) {
	return GraphFileLines(DefaultGraphFile(projectID))
}

func NumberedLines(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	var b strings.Builder
	for i, line := range lines {
		fmt.Fprintf(&b, "%d| %s", i+1, line)
		if i < len(lines)-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func InsertTextAtLine(lines []string, line int, content string) ([]string, error) {
	inserted := splitLines(content)
	if len(inserted) == 0 {
		return nil, errors.New("content is empty")
	}
	if line < 1 || line > len(lines)+1 {
		return nil, fmt.Errorf("line out of range: valid 1..%d", len(lines)+1)
	}
	result := make([]string, 0, len(lines)+len(inserted))
	result = append(result, lines[:line-1]...)
	result = append(result, inserted...)
	result = append(result, lines[line-1:]...)
	return result, nil
}

func DeleteLineRange(lines []string, startLine, endLine int) ([]string, error) {
	if len(lines) == 0 {
		return nil, errors.New("file is empty")
	}
	if startLine < 1 || endLine < startLine || endLine > len(lines) {
		return nil, fmt.Errorf("line range out of range: valid 1..%d", len(lines))
	}
	result := make([]string, 0, len(lines)-(endLine-startLine+1))
	result = append(result, lines[:startLine-1]...)
	result = append(result, lines[endLine:]...)
	return result, nil
}

func LoadWorkingGraph(projectID, version string) (*GraphFile, string, bool, error) {
	if version == "" {
		currentVersion, err := ReadCurrentVersion(projectID)
		if err != nil {
			return nil, "", false, err
		}
		tmpPath := TmpVersionTreePath(projectID, currentVersion)
		if _, err := os.Stat(tmpPath); err == nil {
			graph, loadErr := LoadGraphFile(tmpPath)
			return graph, graph.Meta.BasedOnVersion, true, loadErr
		}
		graph, loadErr := LoadGraphFile(VersionTreePath(projectID, currentVersion))
		return graph, currentVersion, false, loadErr
	}

	tmpPath := TmpVersionTreePath(projectID, version)
	if _, err := os.Stat(tmpPath); err == nil {
		graph, loadErr := LoadGraphFile(tmpPath)
		return graph, version, true, loadErr
	}
	graph, err := LoadGraphFile(VersionTreePath(projectID, version))
	return graph, version, false, err
}

func DefaultGraphFile(projectID string) *GraphFile {
	return &GraphFile{
		SchemaVersion: "fault-tree/v1",
		Tree: GraphTree{
			Name:      projectID,
			TopNodeID: "",
		},
		Nodes: []*GraphNode{},
		Meta:  GraphMeta{},
	}
}

func splitLines(content string) []string {
	normalized := strings.ReplaceAll(content, "\r\n", "\n")
	if normalized == "" {
		return []string{}
	}
	lines := strings.Split(normalized, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func joinLines(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}
