package architecture_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const bunImport = "github.com/uptrace/bun"

func TestUsecaseDoesNotImportBun(t *testing.T) {
	backendRoot := filepath.Clean("../..")
	usecaseRoot := filepath.Join(backendRoot, "internal", "usecase")

	var violations []string
	err := filepath.WalkDir(usecaseRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !strings.Contains(string(content), bunImport) {
			return nil
		}
		rel, err := filepath.Rel(backendRoot, path)
		if err != nil {
			return err
		}
		violations = append(violations, rel)
		return nil
	})
	if err != nil {
		t.Fatalf("walk backend/internal/usecase: %v", err)
	}
	if len(violations) > 0 {
		t.Fatalf("usecase must not import bun directly; violations: %s", strings.Join(violations, ", "))
	}
}
