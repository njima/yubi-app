package architecture_test

import (
	"os"
	"path/filepath"
	"strconv"
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

func TestUsecaseDoesNotDependOnRawDBConnOrTxRunner(t *testing.T) {
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
		if filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		text := string(content)
		if !strings.Contains(text, "repository.DBConn") && !strings.Contains(text, "repository.TxRunner") {
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
		t.Fatalf("usecase must depend on repository.DataAccess instead of raw DBConn/TxRunner; violations: %s", strings.Join(violations, ", "))
	}
}

func TestUsecaseTransactionCallbacksUseLocalConn(t *testing.T) {
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
		if filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(backendRoot, path)
		if err != nil {
			return err
		}
		for lineNumber, line := range strings.Split(string(content), "\n") {
			if !strings.Contains(line, "txData.Conn()") {
				continue
			}
			if strings.Contains(line, "conn := txData.Conn()") {
				continue
			}
			violations = append(violations, rel+":"+strconv.Itoa(lineNumber+1))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk backend/internal/usecase: %v", err)
	}
	if len(violations) > 0 {
		t.Fatalf("transaction callbacks should assign conn := txData.Conn() before repository calls; violations: %s", strings.Join(violations, ", "))
	}
}

func TestUsecaseInterfacesDoNotExposeRepositoryFilters(t *testing.T) {
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
		if filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(backendRoot, path)
		if err != nil {
			return err
		}

		inUsecaseInterface := false
		for lineNumber, line := range strings.Split(string(content), "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "type ") && strings.Contains(trimmed, "Usecase interface {") {
				inUsecaseInterface = true
				continue
			}
			if inUsecaseInterface && trimmed == "}" {
				inUsecaseInterface = false
				continue
			}
			if inUsecaseInterface && strings.Contains(line, "repository.") && strings.Contains(line, "Filter") {
				violations = append(violations, rel+":"+strconv.Itoa(lineNumber+1))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk backend/internal/usecase: %v", err)
	}
	if len(violations) > 0 {
		t.Fatalf("usecase interfaces must expose usecase input/filter types instead of repository filters; violations: %s", strings.Join(violations, ", "))
	}
}
