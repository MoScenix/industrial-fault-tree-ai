package utils

type contextKey string

const (
	// ProjectContextKey stores project-scoped execution context injected by upper services.
	ProjectContextKey contextKey = "project_context"
)
