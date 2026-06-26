package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/airoa-org/yubi-app/backend/internal/infra/database/entity"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

const (
	outputPath   = "./internal/infra/database/schema/"
	outputUpPath = outputPath + "schema.up.sql"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "create schema error: ", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	db := bun.NewDB(nil, pgdialect.New())

	if err := createSchema(ctx, db, entity.Models, entity.IdxCreators); err != nil {
		return err
	}

	fmt.Println("Schema file generated successfully:")
	fmt.Println("  -", outputUpPath)

	return nil
}

func createSchema(ctx context.Context, db *bun.DB, models []any, idxCreators [][]entity.IndexQueryCreator) error {
	upQuery, err := createUpQuery(db, models, idxCreators)
	if err != nil {
		return err
	}

	return createSchemaFile(ctx, upQuery, outputUpPath)
}

func createSchemaFile(ctx context.Context, byteQuery []byte, outputPath string) error {
	select {
	case <-ctx.Done():
		return errors.New("failed to create schema file")
	default:
		if err := os.WriteFile(outputPath, byteQuery, 0644); err != nil {
			return err
		}
		return nil
	}
}

func createUpQuery(db *bun.DB, models []any, idxCreators [][]entity.IndexQueryCreator) ([]byte, error) {
	tableCreatorMap := buildTableCreatorMap(db)

	var result []byte
	for _, model := range models {
		tableName, query := buildTableQuery(db, model, tableCreatorMap)
		formattedQuery := formatCreateTableSQL(query)

		result = append(result, fmt.Sprintf("-- %s\n", tableName)...)
		result = append(result, formattedQuery...)
		result = append(result, ";\n\n"...)
	}

	indexBytes, err := indexesToByte(db, idxCreators)
	if err != nil {
		return nil, err
	}
	result = append(result, indexBytes...)
	return result, nil
}

func buildTableCreatorMap(db *bun.DB) map[string]entity.TableQueryCreator {
	tableCreatorMap := make(map[string]entity.TableQueryCreator)
	for _, tc := range entity.TableCreators {
		query := tc(db)
		tableName := query.GetTableName()
		tableCreatorMap[tableName] = tc
	}
	return tableCreatorMap
}

func buildTableQuery(db *bun.DB, model any, tableCreatorMap map[string]entity.TableQueryCreator) (tableName, query string) {
	defaultQuery := db.NewCreateTable().Model(model).IfNotExists()
	tableName = defaultQuery.GetTableName()

	if tc, ok := tableCreatorMap[tableName]; ok {
		query = tc(db).String()
	} else {
		query = defaultQuery.String()
	}
	return tableName, query
}

func formatCreateTableSQL(query string) string {
	formatted := strings.Replace(query, " (", "(\n  ", 1)
	formatted = strings.Replace(formatted, ", ", ",\n  ", -1)
	if lastIndex := strings.LastIndex(formatted, ")"); lastIndex != -1 {
		formatted = formatted[:lastIndex] + "\n" + formatted[lastIndex:]
	}
	return formatted
}

func indexesToByte(db *bun.DB, idxCreators [][]entity.IndexQueryCreator) ([]byte, error) {
	if len(idxCreators) == 0 {
		return nil, nil
	}

	var data []byte
	data = append(data, "-- indexes\n"...)

	for _, creators := range idxCreators {
		for _, idxCreator := range creators {
			idx := idxCreator(db)
			rawQuery, err := idx.AppendQuery(db.QueryGen(), nil)
			if err != nil {
				return nil, fmt.Errorf("failed to generate index query: %w", err)
			}
			data = append(data, rawQuery...)
			data = append(data, ";\n"...)
		}
	}

	data = append(data, "\n"...)
	return data, nil
}
