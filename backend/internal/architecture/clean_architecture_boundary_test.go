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
		"internal/domain/authz",
		"internal/platform/config",
		"internal/infra/database",
		"internal/event",
		"internal/gen",
		"internal/infra",
		"internal/interfaces",
		"internal/log",
		"internal/usecase/pagination",
		"internal/repository",
		"internal/stack",
		"internal/usecase",
	})
}

func TestUsecaseDoesNotDependOnOuterLayers(t *testing.T) {
	assertNoForbiddenImports(t, "internal/usecase", []string{
		"internal/app",
		"internal/domain/authz",
		"internal/platform/config",
		"internal/infra/database",
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
		"internal/domain/authz",
		"internal/platform/config",
		"internal/infra/database",
		"internal/gen",
		"internal/infra",
		"internal/interfaces",
		"internal/log",
		"internal/stack",
		"internal/usecase",
	})
}

func TestAuthzDoesNotDependOnHTTPBoundary(t *testing.T) {
	assertNoForbiddenImports(t, "internal/domain/authz", []string{
		"internal/shared/requestctx",
		"internal/gen",
		"internal/interfaces",
	})
	assertNoForbiddenExternalImports(t, "internal/domain/authz", []string{
		"github.com/gin-gonic/gin",
	})
}

func TestHTTPControllersDoNotDependOnRepositoryLayer(t *testing.T) {
	assertNoForbiddenImports(t, "internal/interfaces/http/controller", []string{
		"internal/repository",
	})
}

func TestServerBootstrapDoesNotConstructUsecasesDirectly(t *testing.T) {
	backendRoot := filepath.Clean("../..")
	path := filepath.Join(backendRoot, "cmd", "server", "bootstrap.go")

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read cmd/server/bootstrap.go: %v", err)
	}
	if strings.Contains(string(content), "usecase.New") {
		t.Fatalf("cmd/server/bootstrap.go must delegate usecase wiring to cmd/server/usecases.go")
	}
}

func TestServerCompositionRootDoesNotLiveUnderInternalApp(t *testing.T) {
	backendRoot := filepath.Clean("../..")
	path := filepath.Join(backendRoot, "internal", "app")

	if _, err := os.Stat(path); err == nil {
		t.Fatalf("server composition root must live under cmd/server, not internal/app")
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat internal/app: %v", err)
	}
}

func TestRequestContextHelpersLiveUnderShared(t *testing.T) {
	backendRoot := filepath.Clean("../..")
	path := filepath.Join(backendRoot, "internal", "requestctx")

	if _, err := os.Stat(path); err == nil {
		t.Fatalf("request context helpers must live under internal/shared/requestctx, not internal/requestctx")
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat internal/requestctx: %v", err)
	}
}

func TestApplicationErrorsLiveUnderShared(t *testing.T) {
	backendRoot := filepath.Clean("../..")
	path := filepath.Join(backendRoot, "internal", "apperror")

	if _, err := os.Stat(path); err == nil {
		t.Fatalf("application errors must live under internal/shared/apperror, not internal/apperror")
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat internal/apperror: %v", err)
	}
}

func TestAuthorizationPolicyLivesUnderDomain(t *testing.T) {
	backendRoot := filepath.Clean("../..")
	path := filepath.Join(backendRoot, "internal", "authz")

	if _, err := os.Stat(path); err == nil {
		t.Fatalf("authorization policy must live under internal/domain/authz, not internal/authz")
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat internal/authz: %v", err)
	}
}

func TestRuntimeLogHelperLivesUnderPlatform(t *testing.T) {
	backendRoot := filepath.Clean("../..")
	path := filepath.Join(backendRoot, "internal", "log")

	if _, err := os.Stat(path); err == nil {
		t.Fatalf("runtime log helper must live under internal/platform/log, not internal/log")
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat internal/log: %v", err)
	}
}

func TestApplicationEventBusLivesUnderUsecase(t *testing.T) {
	backendRoot := filepath.Clean("../..")
	path := filepath.Join(backendRoot, "internal", "event")

	if _, err := os.Stat(path); err == nil {
		t.Fatalf("application event bus must live under internal/usecase/eventbus, not internal/event")
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat internal/event: %v", err)
	}
}

func TestDatabaseInfrastructureLivesUnderInfra(t *testing.T) {
	backendRoot := filepath.Clean("../..")
	path := filepath.Join(backendRoot, "internal", "database")

	if _, err := os.Stat(path); err == nil {
		t.Fatalf("database infrastructure must live under internal/infra/database, not internal/database")
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat internal/database: %v", err)
	}
}

func TestRuntimeConfigLivesUnderPlatform(t *testing.T) {
	backendRoot := filepath.Clean("../..")
	path := filepath.Join(backendRoot, "internal", "config")

	if _, err := os.Stat(path); err == nil {
		t.Fatalf("runtime config must live under internal/platform/config, not internal/config")
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat internal/config: %v", err)
	}
}

func TestPaginationRequestModelLivesUnderUsecase(t *testing.T) {
	backendRoot := filepath.Clean("../..")
	path := filepath.Join(backendRoot, "internal", "pagination")

	if _, err := os.Stat(path); err == nil {
		t.Fatalf("pagination request model must live under internal/usecase/pagination, not internal/pagination")
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat internal/pagination: %v", err)
	}
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
