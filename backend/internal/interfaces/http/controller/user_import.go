package controller

import (
	"context"

	"github.com/airoa-org/yubi-app/backend/internal/gen/openapi"
	"github.com/airoa-org/yubi-app/backend/internal/shared/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/usecase"
)

func (c *controller) ValidateUserImport(ctx context.Context, request openapi.ValidateUserImportRequestObject) (openapi.ValidateUserImportResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	result, err := c.userImportUsecase.Validate(ctx, request.Body.CsvContent)
	if err != nil {
		return nil, err
	}

	return openapi.ValidateUserImport200JSONResponse(toOpenAPIUserValidationResponse(result)), nil
}

func (c *controller) ImportUsers(ctx context.Context, request openapi.ImportUsersRequestObject) (openapi.ImportUsersResponseObject, error) {
	if request.Body == nil {
		return nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "request body is required"))
	}

	result, err := c.userImportUsecase.Import(ctx, request.Body.CsvContent)
	if err != nil {
		return nil, err
	}

	return openapi.ImportUsers200JSONResponse(toOpenAPIUserImportResponse(result)), nil
}

func toOpenAPIUserValidationResponse(r usecase.UserImportValidationResult) openapi.UserImportValidationResponse {
	validRows := make([]openapi.UserImportValidRow, 0, len(r.ValidRows))
	for _, row := range r.ValidRows {
		roleStr := usecase.UserRoleString(row.Role)
		validRows = append(validRows, openapi.UserImportValidRow{
			RowNumber:   row.RowNumber,
			Email:       row.Email,
			DisplayName: row.DisplayName,
			Role:        roleStr,
		})
	}

	duplicateRows := make([]openapi.UserImportRowError, 0, len(r.DuplicateRows))
	for _, row := range r.DuplicateRows {
		duplicateRows = append(duplicateRows, toOpenAPIUserRowError(row))
	}

	errorRows := make([]openapi.UserImportRowError, 0, len(r.ErrorRows))
	for _, row := range r.ErrorRows {
		errorRows = append(errorRows, toOpenAPIUserRowError(row))
	}

	return openapi.UserImportValidationResponse{
		ValidRows:     validRows,
		DuplicateRows: duplicateRows,
		ErrorRows:     errorRows,
		Summary: openapi.UserImportSummary{
			ValidCount:     len(r.ValidRows),
			DuplicateCount: len(r.DuplicateRows),
			ErrorCount:     len(r.ErrorRows),
		},
	}
}

func toOpenAPIUserImportResponse(r usecase.UserImportResult) openapi.UserImportResponse {
	errors := make([]openapi.UserImportRowError, 0, len(r.Errors))
	for _, e := range r.Errors {
		errors = append(errors, toOpenAPIUserRowError(e))
	}

	return openapi.UserImportResponse{
		ImportedCount: r.ImportedCount,
		SkippedCount:  r.SkippedCount,
		ErrorCount:    r.ErrorCount,
		Errors:        errors,
	}
}

func toOpenAPIUserRowError(row usecase.UserImportRowError) openapi.UserImportRowError {
	e := openapi.UserImportRowError{
		RowNumber: row.RowNumber,
		Errors:    row.Errors,
	}
	if row.Email != "" {
		e.Email = &row.Email
	}
	return e
}
