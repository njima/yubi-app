package controller

import (
	"context"
	"testing"

	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
)

func TestUpdateUserById_RequiresRequestBody(t *testing.T) {
	c := &controller{}

	resp, err := c.UpdateUserById(context.Background(), openapi.UpdateUserByIdRequestObject{})

	if err == nil {
		t.Fatalf("expected error for nil request body")
	}
	if resp != nil {
		t.Fatalf("expected nil response for nil request body, got %T", resp)
	}
}
