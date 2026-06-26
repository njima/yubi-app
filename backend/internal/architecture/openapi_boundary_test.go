package architecture_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const openAPIImport = `"github.com/airoa-org/yubi-app/backend/internal/gen/` + `openapi"`

func TestOpenAPIImportsStayAtHTTPBoundary(t *testing.T) {
	backendRoot := filepath.Clean("../..")
	internalRoot := filepath.Join(backendRoot, "internal")

	var violations []string
	err := filepath.WalkDir(internalRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == "gen" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !strings.Contains(string(content), openAPIImport) {
			return nil
		}
		rel, err := filepath.Rel(backendRoot, path)
		if err != nil {
			return err
		}
		if isAllowedOpenAPIImport(rel) {
			return nil
		}
		violations = append(violations, rel)
		return nil
	})
	if err != nil {
		t.Fatalf("walk backend/internal: %v", err)
	}
	if len(violations) > 0 {
		t.Fatalf("OpenAPI imports must stay at HTTP boundaries; violations: %s", strings.Join(violations, ", "))
	}
}

func isAllowedOpenAPIImport(path string) bool {
	return strings.HasPrefix(path, "internal/interfaces/http/")
}
