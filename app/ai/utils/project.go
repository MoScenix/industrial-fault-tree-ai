package utils

import (
	"encoding/json"
	"errors"
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
		Meta: GraphMeta{},
	}
}
