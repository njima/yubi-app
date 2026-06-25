package architecture_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const backendModule = "github.com/airoa-org/yubi-app/backend/"

func TestDomainDoesNotDependOnOuterLayers(t *testing.T) {
	assertNoForbiddenImports(t, "internal/domain", []string{
		"internal/app",
		"internal/authz",
		"internal/config",
		"internal/database",
		"internal/event",
		"internal/gen",
		"internal/infra",
		"internal/interfaces",
		"internal/log",
		"internal/pagination",
		"internal/repository",
		"internal/stack",
		"internal/usecase",
	})
}

func TestUsecaseDoesNotDependOnOuterLayers(t *testing.T) {
	assertNoForbiddenImports(t, "internal/usecase", []string{
		"internal/app",
		"internal/authz",
		"internal/config",
		"internal/database",
		"internal/gen",
		"internal/infra",
		"internal/interfaces",
		"internal/log",
		"internal/stack",
	})
}

func TestRepositoryInterfacesDoNotDependOnImplementations(t *testing.T) {
	assertNoForbiddenImports(t, "internal/repository", []string{
		"internal/app",
		"internal/authz",
		"internal/config",
		"internal/database",
		"internal/gen",
		"internal/infra",
		"internal/interfaces",
		"internal/log",
		"internal/stack",
		"internal/usecase",
	})
}

func TestAuthzDoesNotDependOnHTTPBoundary(t *testing.T) {
	assertNoForbiddenImports(t, "internal/authz", []string{
		"internal/ccontext",
		"internal/gen",
		"internal/interfaces",
	})
	assertNoForbiddenExternalImports(t, "internal/authz", []string{
		"github.com/gin-gonic/gin",
	})
}

func assertNoForbiddenImports(t *testing.T, packageDir string, forbiddenPrefixes []string) {
	t.Helper()

	backendRoot := filepath.Clean("../..")
	root := filepath.Join(backendRoot, packageDir)

	var violations []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		imports, err := parseImports(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(backendRoot, path)
		if err != nil {
			return err
		}
		for _, importPath := range imports {
			internalPath, ok := strings.CutPrefix(importPath, backendModule)
			if !ok {
				continue
			}
			for _, forbiddenPrefix := range forbiddenPrefixes {
				if internalPath == forbiddenPrefix || strings.HasPrefix(internalPath, forbiddenPrefix+"/") {
					violations = append(violations, rel+" imports "+internalPath)
					break
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", packageDir, err)
	}
	if len(violations) > 0 {
		t.Fatalf("%s must not import forbidden outer layers: %s", packageDir, strings.Join(violations, ", "))
	}
}

func assertNoForbiddenExternalImports(t *testing.T, packageDir string, forbiddenImports []string) {
	t.Helper()

	backendRoot := filepath.Clean("../..")
	root := filepath.Join(backendRoot, packageDir)

	forbidden := make(map[string]struct{}, len(forbiddenImports))
	for _, importPath := range forbiddenImports {
		forbidden[importPath] = struct{}{}
	}

	var violations []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		imports, err := parseImports(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(backendRoot, path)
		if err != nil {
			return err
		}
		for _, importPath := range imports {
			if _, ok := forbidden[importPath]; ok {
				violations = append(violations, rel+" imports "+importPath)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", packageDir, err)
	}
	if len(violations) > 0 {
		t.Fatalf("%s must not import forbidden HTTP boundary packages: %s", packageDir, strings.Join(violations, ", "))
	}
}

func parseImports(path string) ([]string, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	imports := make([]string, 0, len(file.Imports))
	for _, spec := range file.Imports {
		imports = append(imports, strings.Trim(spec.Path.Value, `"`))
	}
	return imports, nil
}
