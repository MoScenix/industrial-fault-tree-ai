package utils

import (
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/MoScenix/industrial-fault-tree-ai/app/graph/conf"
	"github.com/google/uuid"
)

const (
	VersionsDirName   = "versions"
	TmpDirName        = "tmp"
	SuggestionDirName = "suggestions"
	CurrentFileName   = "current"
	DefaultVersion    = "v001"
)

type VersionInfo struct {
	Version    string
	IsCurrent  bool
	CreateTime string
	UpdateTime string
}

func NewProjectUUID() string {
	return uuid.New().String()
}

func ProjectDir(projectUUID string) string {
	return filepath.Join(conf.GetConf().Graph.RootDir, projectUUID)
}

func VersionsDir(projectDir string) string {
	return filepath.Join(projectDir, VersionsDirName)
}

func VersionDir(projectDir, version string) string {
	return filepath.Join(VersionsDir(projectDir), version)
}

func TmpDir(projectDir string) string {
	return filepath.Join(projectDir, TmpDirName)
}

func TmpVersionDir(projectDir, version string) string {
	return filepath.Join(TmpDir(projectDir), version)
}

func SuggestionsDir(projectDir string) string {
	return filepath.Join(projectDir, SuggestionDirName)
}

func SuggestionPath(projectDir, version string) string {
	return filepath.Join(SuggestionsDir(projectDir), version+".md")
}

func CurrentVersionPath(projectDir string) string {
	return filepath.Join(projectDir, CurrentFileName)
}

func EnsureProjectLayout(projectDir string) error {
	paths := []string{
		projectDir,
		VersionsDir(projectDir),
		VersionDir(projectDir, DefaultVersion),
		TmpDir(projectDir),
		SuggestionsDir(projectDir),
	}

	for _, p := range paths {
		if err := os.MkdirAll(p, os.ModePerm); err != nil {
			return err
		}
	}

	return WriteCurrentVersion(projectDir, DefaultVersion)
}

func RemoveProjectLayout(projectDir string) error {
	return os.RemoveAll(projectDir)
}

func ReadCurrentVersion(projectDir string) (string, error) {
	content, err := os.ReadFile(CurrentVersionPath(projectDir))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func WriteCurrentVersion(projectDir, version string) error {
	if err := os.MkdirAll(projectDir, os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(CurrentVersionPath(projectDir), []byte(version), 0644)
}

func HasTmp(projectDir string) bool {
	entries, err := os.ReadDir(TmpDir(projectDir))
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if entry.IsDir() {
			return true
		}
	}
	return false
}

func HasTmpVersion(projectDir, version string) bool {
	entries, err := os.ReadDir(TmpVersionDir(projectDir, version))
	return err == nil && len(entries) > 0
}

func HasTmpTree(projectDir, version string) bool {
	info, err := os.Stat(filepath.Join(TmpVersionDir(projectDir, version), "tree.json"))
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func EnsureTmpFromVersion(projectDir, version string) error {
	if HasTmpVersion(projectDir, version) {
		return nil
	}

	src := VersionDir(projectDir, version)
	dst := TmpVersionDir(projectDir, version)
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	if err := os.MkdirAll(dst, os.ModePerm); err != nil {
		return err
	}
	return CopyDir(src, dst)
}

func SaveTmpToVersion(projectDir, fromVersion, toVersion string) error {
	src := TmpVersionDir(projectDir, fromVersion)
	dst := VersionDir(projectDir, toVersion)
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	if err := os.MkdirAll(dst, os.ModePerm); err != nil {
		return err
	}
	return CopyDir(src, dst)
}

func ClearTmp(projectDir, version string) error {
	dst := TmpVersionDir(projectDir, version)
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	return nil
}

func CreateVersionFromCurrent(projectDir, currentVersion, newVersion string) error {
	src := VersionDir(projectDir, currentVersion)
	dst := VersionDir(projectDir, newVersion)
	if err := os.MkdirAll(VersionsDir(projectDir), os.ModePerm); err != nil {
		return err
	}
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	if err := os.MkdirAll(dst, os.ModePerm); err != nil {
		return err
	}
	return CopyDir(src, dst)
}

func DeleteVersionDir(projectDir, version string) error {
	return os.RemoveAll(VersionDir(projectDir, version))
}

func RenameVersionDir(projectDir, version, versionName string) error {
	if version == versionName || versionName == "" {
		return nil
	}
	return os.Rename(VersionDir(projectDir, version), VersionDir(projectDir, versionName))
}

func ListVersions(projectDir, currentVersion string) ([]VersionInfo, error) {
	entries, err := os.ReadDir(VersionsDir(projectDir))
	if err != nil {
		if os.IsNotExist(err) {
			return []VersionInfo{}, nil
		}
		return nil, err
	}

	versionList := make([]VersionInfo, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		versionList = append(versionList, VersionInfo{
			Version:    entry.Name(),
			IsCurrent:  entry.Name() == currentVersion,
			CreateTime: FormatTime(info.ModTime()),
			UpdateTime: FormatTime(info.ModTime()),
		})
	}

	sort.Slice(versionList, func(i, j int) bool {
		return versionList[i].Version < versionList[j].Version
	})

	return versionList, nil
}

func CopyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
				return err
			}
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
			continue
		}
		if err := copyFile(srcPath, dstPath); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
