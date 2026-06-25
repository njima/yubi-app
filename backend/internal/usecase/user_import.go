package usecase

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/mail"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/rs/zerolog"

	"github.com/airoa-org/yubi-app/backend/internal/apperror"
	"github.com/airoa-org/yubi-app/backend/internal/ccontext"
	"github.com/airoa-org/yubi-app/backend/internal/domain/model"
	"github.com/airoa-org/yubi-app/backend/internal/repository"
)

const (
	maxUserImportRows    = 50
	maxUserImportFileLen = 5 * 1024 * 1024 // 5MB
)

var userImportHeaders = []string{"email", "display_name", "role"}

var userRoleMap = map[string]model.UserRole{
	"admin":         model.UserRoleAdmin,
	"data_engineer": model.UserRoleDataEngineer,
	"manager":       model.UserRoleManager,
	"operator":      model.UserRoleOperator,
	"viewer":        model.UserRoleViewer,
}

var userRoleStringMap = func() map[model.UserRole]string {
	m := make(map[model.UserRole]string, len(userRoleMap))
	for k, v := range userRoleMap {
		m[v] = k
	}
	return m
}()

func UserRoleString(role model.UserRole) string {
	if s, ok := userRoleStringMap[role]; ok {
		return s
	}
	return "viewer"
}

type UserImportUsecase interface {
	Validate(ctx context.Context, csvContent string) (UserImportValidationResult, error)
	Import(ctx context.Context, csvContent string) (UserImportResult, error)
}

type UserImportValidationResult struct {
	ValidRows     []UserImportValidRow
	DuplicateRows []UserImportRowError
	ErrorRows     []UserImportRowError
}

type UserImportResult struct {
	ImportedCount int
	SkippedCount  int
	ErrorCount    int
	Errors        []UserImportRowError
}

type UserImportValidRow struct {
	RowNumber   int
	Email       string
	DisplayName string
	Role        model.UserRole
}

type UserImportRowError struct {
	RowNumber int
	Email     string
	Errors    []string
}

type userImport struct {
	userRepo repository.User
	db       repository.DBConn
	logger   zerolog.Logger
}

func NewUserImport(
	userRepo repository.User,
	db repository.DBConn,
	logger zerolog.Logger,
) UserImportUsecase {
	return &userImport{
		userRepo: userRepo,
		db:       db,
		logger:   logger,
	}
}

func (u *userImport) Validate(ctx context.Context, csvContent string) (UserImportValidationResult, error) {
	rows, parseErrors, err := u.parseCSV(csvContent)
	if err != nil {
		return UserImportValidationResult{}, err
	}
	result, err := u.validateInternal(ctx, rows, parseErrors)
	return result, err
}

func (u *userImport) Import(ctx context.Context, csvContent string) (UserImportResult, error) {
	orgID, err := ccontext.OrganizationID(ctx)
	if err != nil {
		return UserImportResult{}, err
	}

	rows, parseErrors, err := u.parseCSV(csvContent)
	if err != nil {
		return UserImportResult{}, err
	}

	validation, err := u.validateInternal(ctx, rows, parseErrors)
	if err != nil {
		return UserImportResult{}, err
	}

	importErrors := append([]UserImportRowError{}, validation.ErrorRows...)
	importedCount := 0
	skippedCount := len(validation.DuplicateRows)

	for _, row := range validation.ValidRows {
		nu, err := model.InitUser(orgID, row.DisplayName, row.Email, row.Role)
		if err != nil {
			u.logger.Error().Err(err).Str("email", row.Email).Int("row", row.RowNumber).Msg("failed to initialize user model")
			importErrors = append(importErrors, UserImportRowError{
				RowNumber: row.RowNumber,
				Email:     row.Email,
				Errors:    []string{"failed to create user"},
			})
			continue
		}

		if _, err := u.userRepo.Create(ctx, u.db, nu); err != nil {
			u.logger.Error().Err(err).Str("email", row.Email).Int("row", row.RowNumber).Msg("failed to create user in database")
			importErrors = append(importErrors, UserImportRowError{
				RowNumber: row.RowNumber,
				Email:     row.Email,
				Errors:    []string{"failed to save user"},
			})
			continue
		}

		importedCount++
	}

	sort.Slice(importErrors, func(i, j int) bool {
		return importErrors[i].RowNumber < importErrors[j].RowNumber
	})

	return UserImportResult{
		ImportedCount: importedCount,
		SkippedCount:  skippedCount,
		ErrorCount:    len(importErrors),
		Errors:        importErrors,
	}, nil
}

func (u *userImport) validateInternal(
	ctx context.Context,
	rows []userImportRow,
	parseErrors []UserImportRowError,
) (UserImportValidationResult, error) {
	emails := make([]string, 0, len(rows))
	for _, r := range rows {
		if r.email != "" {
			emails = append(emails, strings.ToLower(r.email))
		}
	}

	existingEmails, err := u.userRepo.ExistsByEmails(ctx, u.db, emails)
	if err != nil {
		return UserImportValidationResult{}, err
	}

	var result UserImportValidationResult
	seenEmails := make(map[string]bool)

	for _, row := range rows {
		errs := validateUserImportRow(row)
		if len(errs) > 0 {
			result.ErrorRows = append(result.ErrorRows, UserImportRowError{
				RowNumber: row.rowNumber,
				Email:     row.email,
				Errors:    errs,
			})
			continue
		}

		emailLower := strings.ToLower(row.email)
		if existingEmails[emailLower] || seenEmails[emailLower] {
			result.DuplicateRows = append(result.DuplicateRows, UserImportRowError{
				RowNumber: row.rowNumber,
				Email:     row.email,
				Errors:    []string{fmt.Sprintf("email %q already exists", row.email)},
			})
			continue
		}

		seenEmails[emailLower] = true
		result.ValidRows = append(result.ValidRows, UserImportValidRow{
			RowNumber:   row.rowNumber,
			Email:       row.email,
			DisplayName: row.displayName,
			Role:        row.role,
		})
	}

	result.ErrorRows = append(result.ErrorRows, parseErrors...)
	sort.Slice(result.ErrorRows, func(i, j int) bool {
		return result.ErrorRows[i].RowNumber < result.ErrorRows[j].RowNumber
	})
	return result, nil
}

type userImportRow struct {
	rowNumber   int
	email       string
	displayName string
	role        model.UserRole
}

func (u *userImport) parseCSV(csvContent string) ([]userImportRow, []UserImportRowError, error) {
	if len(csvContent) > maxUserImportFileLen {
		return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "CSV content exceeds maximum size of 5MB"))
	}

	csvContent = strings.TrimPrefix(csvContent, "\xef\xbb\xbf")

	if !utf8.ValidString(csvContent) {
		return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "CSV content is not valid UTF-8"))
	}

	reader := csv.NewReader(strings.NewReader(csvContent))
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	header, err := reader.Read()
	if err == io.EOF {
		return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "CSV must have a header row and at least one data row"))
	}
	if err != nil {
		return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "failed to parse CSV: %v", err))
	}
	if len(header) != len(userImportHeaders) {
		return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest,
			"CSV header must have exactly %d columns: %s", len(userImportHeaders), strings.Join(userImportHeaders, ", ")))
	}
	for i, h := range header {
		if strings.TrimSpace(strings.ToLower(h)) != userImportHeaders[i] {
			return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest,
				"CSV column %d must be %q, got %q", i+1, userImportHeaders[i], h))
		}
	}

	var rows []userImportRow
	var parseErrors []UserImportRowError
	rowNum := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "failed to parse CSV: %v", err))
		}
		rowNum++
		if rowNum > maxUserImportRows {
			return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest,
				"CSV exceeds maximum of %d rows", maxUserImportRows))
		}
		if len(record) != len(userImportHeaders) {
			parseErrors = append(parseErrors, UserImportRowError{
				RowNumber: rowNum,
				Errors:    []string{fmt.Sprintf("expected %d columns, got %d", len(userImportHeaders), len(record))},
			})
			continue
		}

		roleStr := strings.ToLower(strings.TrimSpace(record[2]))
		role := model.UserRoleViewer
		if roleStr != "" {
			if r, ok := userRoleMap[roleStr]; ok {
				role = r
			} else {
				parseErrors = append(parseErrors, UserImportRowError{
					RowNumber: rowNum,
					Email:     strings.TrimSpace(record[0]),
					Errors:    []string{fmt.Sprintf("role %q is invalid", record[2])},
				})
				continue
			}
		}

		rows = append(rows, userImportRow{
			rowNumber:   rowNum,
			email:       strings.TrimSpace(record[0]),
			displayName: strings.TrimSpace(record[1]),
			role:        role,
		})
	}

	if rowNum == 0 {
		return nil, nil, apperror.NewError(apperror.NewMessage(apperror.CodeBadRequest, "CSV must have a header row and at least one data row"))
	}

	return rows, parseErrors, nil
}

func validateUserImportRow(row userImportRow) []string {
	var errs []string

	if row.email == "" {
		errs = append(errs, "email is required")
	} else if addr, err := mail.ParseAddress(row.email); err != nil || addr.Name != "" {
		errs = append(errs, "email is invalid")
	}

	if row.displayName == "" {
		errs = append(errs, "display_name is required")
	} else if len([]rune(row.displayName)) > 60 {
		errs = append(errs, "display_name must be 60 characters or less")
	}

	return errs
}
