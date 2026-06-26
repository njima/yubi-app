package controller

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

func (c *controller) ValidateTaskImport(ctx context.Context, request openapi.ValidateTaskImportRequestObject) (openapi.ValidateTaskImportResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	result, err := c.taskImportUsecase.Validate(ctx, body.CsvContent)
	if err != nil {
		return nil, err
	}

	return openapi.ValidateTaskImport200JSONResponse(toOpenAPIValidationResponse(result)), nil
}

func (c *controller) ImportTasks(ctx context.Context, request openapi.ImportTasksRequestObject) (openapi.ImportTasksResponseObject, error) {
	body, err := requiredBody(request.Body)
	if err != nil {
		return nil, err
	}

	result, err := c.taskImportUsecase.Import(ctx, body.CsvContent)
	if err != nil {
		return nil, err
	}

	return openapi.ImportTasks200JSONResponse(toOpenAPIImportResponse(result)), nil
}

func toOpenAPIValidationResponse(r usecase.TaskImportValidationResult) openapi.TaskImportValidationResponse {
	validRows := make([]openapi.TaskImportRow, 0, len(r.ValidRows))
	for _, row := range r.ValidRows {
		validRows = append(validRows, toOpenAPIImportRow(row))
	}

	duplicateRows := make([]openapi.TaskImportRowError, 0, len(r.DuplicateRows))
	for _, row := range r.DuplicateRows {
		duplicateRows = append(duplicateRows, toOpenAPIRowError(row))
	}

	errorRows := make([]openapi.TaskImportRowError, 0, len(r.ErrorRows))
	for _, row := range r.ErrorRows {
		errorRows = append(errorRows, toOpenAPIRowError(row))
	}

	return openapi.TaskImportValidationResponse{
		ValidRows:     validRows,
		DuplicateRows: duplicateRows,
		ErrorRows:     errorRows,
		Summary: openapi.TaskImportSummary{
			ValidCount:     len(r.ValidRows),
			DuplicateCount: len(r.DuplicateRows),
			ErrorCount:     len(r.ErrorRows),
		},
	}
}

func toOpenAPIImportResponse(r usecase.TaskImportResult) openapi.TaskImportResponse {
	errors := make([]openapi.TaskImportRowError, 0, len(r.Errors))
	for _, e := range r.Errors {
		errors = append(errors, toOpenAPIRowError(e))
	}

	return openapi.TaskImportResponse{
		ImportedCount: r.ImportedCount,
		SkippedCount:  r.SkippedCount,
		ErrorCount:    r.ErrorCount,
		Errors:        errors,
	}
}

func toOpenAPIImportRow(row usecase.TaskImportRow) openapi.TaskImportRow {
	r := openapi.TaskImportRow{
		RowNumber:  row.RowNumber,
		Name:       row.Name,
		ManualUrl:  row.ManualURL,
		Priority:   row.Priority,
		Difficulty: row.Difficulty,
		Deadline:   row.Deadline,
	}
	if row.Description != "" {
		r.Description = &row.Description
	}
	if row.Status != "" {
		r.Status = &row.Status
	}
	if row.RobotType != "" {
		r.RobotType = &row.RobotType
	}
	if row.Tags != "" {
		r.Tags = &row.Tags
	}
	return r
}

func toOpenAPIRowError(row usecase.TaskImportRowError) openapi.TaskImportRowError {
	r := openapi.TaskImportRowError{
		RowNumber: row.RowNumber,
		Errors:    row.Errors,
	}
	if row.Name != "" {
		r.Name = &row.Name
	}
	return r
}
