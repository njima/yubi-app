package apperror

import (
	"net/http"
	"testing"
)

func TestKind_HTTPStatus(t *testing.T) {
	tests := []struct {
		name string
		kind Kind
		want int
	}{
		{name: "KindNotFound → 404", kind: KindNotFound, want: http.StatusNotFound},
		{name: "KindBadRequest → 400", kind: KindBadRequest, want: http.StatusBadRequest},
		{name: "KindUnauthorized → 401", kind: KindUnauthorized, want: http.StatusUnauthorized},
		{name: "KindForbidden → 403", kind: KindForbidden, want: http.StatusForbidden},
		{name: "KindConflict → 409", kind: KindConflict, want: http.StatusConflict},
		{name: "KindInternal → 500", kind: KindInternal, want: http.StatusInternalServerError},
		{name: "KindValidation → 400", kind: KindValidation, want: http.StatusBadRequest},
		{name: "KindEmpty → 500 (default)", kind: KindEmpty, want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.kind.HTTPStatus()
			if got != tt.want {
				t.Errorf("Kind.HTTPStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSameKind(t *testing.T) {
	tests := []struct {
		name string
		err  error
		kind Kind
		want bool
	}{
		{
			name: "direct error with matching Kind",
			err:  NewError(NewMessage(CodeNotFound, "not found")),
			kind: KindNotFound,
			want: true,
		},
		{
			name: "direct error with non-matching Kind",
			err:  NewError(NewMessage(CodeNotFound, "not found")),
			kind: KindBadRequest,
			want: false,
		},
		{
			name: "wrapped error with matching Kind in inner layer",
			err:  WrapWithMessage(NewError(NewMessage(CodeNotFound, "not found")), NewMessage(CodeEmpty, "")),
			kind: KindNotFound,
			want: true,
		},
		{
			name: "nil error returns false",
			err:  nil,
			kind: KindNotFound,
			want: false,
		},
		{
			name: "non-apperror error returns false",
			err:  &nonAppError{msg: "plain error"},
			kind: KindNotFound,
			want: false,
		},
		{
			name: "KindConflict error matches KindConflict",
			err:  NewError(NewMessage(CodeConflict, "conflict")),
			kind: KindConflict,
			want: true,
		},
		{
			name: "KindValidation error matches KindValidation",
			err:  NewError(NewMessage(CodeValidationError, "validation")),
			kind: KindValidation,
			want: true,
		},
		{
			name: "KindForbidden error does not match KindConflict",
			err:  NewError(NewMessage(CodeForbidden, "forbidden")),
			kind: KindConflict,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SameKind(tt.err, tt.kind)
			if got != tt.want {
				t.Errorf("SameKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

// CodeNotFound is a convenience alias for test clarity
var CodeNotFound = CodeOrganizationNotFound

func TestNewErrorResponse(t *testing.T) {
	tests := []struct {
		name     string
		codes    []Code
		wantCode int
		wantMsg  string
	}{
		{
			name:     "not found error",
			codes:    []Code{CodeOrganizationNotFound},
			wantCode: http.StatusNotFound,
			wantMsg:  "Organization not found", // MessageForEndUser takes priority when non-empty
		},
		{
			name:     "validation error",
			codes:    []Code{CodeValidationError},
			wantCode: http.StatusBadRequest,
			wantMsg:  "Validation failed", // MessageForEndUser takes priority when non-empty
		},
		{
			name:     "forbidden error",
			codes:    []Code{CodeForbidden},
			wantCode: http.StatusForbidden,
			wantMsg:  "Forbidden",
		},
		{
			name:     "unauthorized error",
			codes:    []Code{CodeUnauthorized},
			wantCode: http.StatusUnauthorized,
			wantMsg:  "Unauthorized",
		},
		{
			name:     "conflict error",
			codes:    []Code{CodeConflict},
			wantCode: http.StatusConflict,
			wantMsg:  "Resource conflict",
		},
		{
			name:     "internal error",
			codes:    []Code{CodeInternal},
			wantCode: http.StatusInternalServerError,
			wantMsg:  "Internal server error",
		},
		{
			name:     "empty codes falls back to internal error",
			codes:    []Code{},
			wantCode: http.StatusInternalServerError,
			wantMsg:  "Internal server error",
		},
		{
			name:     "multiple codes uses first code",
			codes:    []Code{CodeOrganizationNotFound, CodeValidationError},
			wantCode: http.StatusNotFound,
			wantMsg:  "Organization not found", // MessageForEndUser of first code
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewErrorResponse(tt.codes)
			if got.Code != tt.wantCode {
				t.Errorf("NewErrorResponse() Code = %v, want %v", got.Code, tt.wantCode)
			}
			if got.Message != tt.wantMsg {
				t.Errorf("NewErrorResponse() Message = %v, want %v", got.Message, tt.wantMsg)
			}
		})
	}
}

// nonAppError is a plain error type for testing non-apperror cases.
type nonAppError struct {
	msg string
}

func (e *nonAppError) Error() string { return e.msg }
