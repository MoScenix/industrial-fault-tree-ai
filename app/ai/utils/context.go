package utils

type contextKey string

const (
	// ProjectContextKey stores project-scoped execution context injected by upper services.
	ProjectContextKey contextKey = "project_context"
)

type ProjectContext struct {
	ProjectID       string
	DeviceName      string
	TopEvent        string
	CurrentVersion  string
	TmpVersionReady bool
	DocumentSummary string
}
