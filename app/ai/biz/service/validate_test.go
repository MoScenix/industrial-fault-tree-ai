package service

import (
	"context"
	"testing"
)

func TestValidate_Run(t *testing.T) {
	_ = NewValidateService(context.Background())
	t.Skip("skeleton module: validate business logic is not implemented yet")
}
