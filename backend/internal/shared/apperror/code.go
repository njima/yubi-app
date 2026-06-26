package apperror

import "net/http"

// Domain identifies the bounded context or resource an error belongs to.
type Domain string

const (
	DomainEmpty        Domain = ""
	DomainOrganization Domain = "organization"
	DomainUser         Domain = "user"
	DomainRobot        Domain = "robot"
	DomainSite         Domain = "site"
	DomainLocation     Domain = "location"
	DomainTask         Domain = "task"
	DomainEpisode      Domain = "episode"
	DomainSubTask      Domain = "subtask"
	DomainAPIKey       Domain = "api_key"
	DomainAuth         Domain = "auth"
	DomainDatabase     Domain = "database"
	DomainValidation   Domain = "validation"
)

// Kind classifies an error into a broad category that maps to an HTTP status code.
type Kind int

const (
	KindEmpty Kind = iota
	KindNotFound
	KindBadRequest
	KindUnauthorized
	KindForbidden
	KindConflict
	KindInternal
	KindValidation
)

// HTTPStatus returns the HTTP status code that corresponds to this Kind.
func (k Kind) HTTPStatus() int {
	switch k {
	case KindNotFound:
		return http.StatusNotFound
	case KindBadRequest:
		return http.StatusBadRequest
	case KindUnauthorized:
		return http.StatusUnauthorized
	case KindForbidden:
		return http.StatusForbidden
	case KindConflict:
		return http.StatusConflict
	case KindInternal:
		return http.StatusInternalServerError
	case KindValidation:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// ResponseWrapErrorMessage returns a generic end-user message for this Kind.
func (k Kind) ResponseWrapErrorMessage() string {
	switch k {
	case KindNotFound:
		return "Resource not found"
	case KindBadRequest:
		return "Bad request"
	case KindUnauthorized:
		return "Unauthorized"
	case KindForbidden:
		return "Forbidden"
	case KindConflict:
		return "Resource conflict"
	case KindInternal:
		return "Internal server error"
	case KindValidation:
		return "Validation error"
	default:
		return "Unknown error"
	}
}

// Code is a structured error code that carries a Kind, Domain, machine-readable Code string,
// and an optional end-user-facing message.
type Code struct {
	Kind              Kind
	Domain            Domain
	Code              string
	MessageForEndUser string
}

var (
	CodeEmpty = Code{Kind: KindEmpty, Domain: DomainEmpty, Code: "", MessageForEndUser: ""}

	// NotFound errors
	CodeOrganizationNotFound   = Code{Kind: KindNotFound, Domain: DomainOrganization, Code: "organization_not_found", MessageForEndUser: "Organization not found"}
	CodeUserNotFound           = Code{Kind: KindNotFound, Domain: DomainUser, Code: "user_not_found", MessageForEndUser: "User not found"}
	CodeRobotNotFound          = Code{Kind: KindNotFound, Domain: DomainRobot, Code: "robot_not_found", MessageForEndUser: "Robot not found"}
	CodeSiteNotFound           = Code{Kind: KindNotFound, Domain: DomainSite, Code: "site_not_found", MessageForEndUser: "Site not found"}
	CodeLocationNotFound       = Code{Kind: KindNotFound, Domain: DomainLocation, Code: "location_not_found", MessageForEndUser: "Location not found"}
	CodeTaskNotFound           = Code{Kind: KindNotFound, Domain: DomainTask, Code: "task_not_found", MessageForEndUser: "Task not found"}
	CodeTaskVersionNotFound    = Code{Kind: KindNotFound, Domain: DomainTask, Code: "task_version_not_found", MessageForEndUser: "Task version not found"}
	CodeEpisodeNotFound        = Code{Kind: KindNotFound, Domain: DomainEpisode, Code: "episode_not_found", MessageForEndUser: "Episode not found"}
	CodeSubTaskNotFound        = Code{Kind: KindNotFound, Domain: DomainSubTask, Code: "subtask_not_found", MessageForEndUser: "SubTask not found"}
	CodeEpisodeSubTaskNotFound = Code{Kind: KindNotFound, Domain: DomainEpisode, Code: "episode_subtask_not_found", MessageForEndUser: "Episode SubTask not found"}
	CodeExecutionNotFound      = Code{Kind: KindNotFound, Domain: DomainEpisode, Code: "execution_not_found", MessageForEndUser: "Execution not found"}
	CodeAPIKeyNotFound         = Code{Kind: KindNotFound, Domain: DomainAPIKey, Code: "api_key_not_found", MessageForEndUser: "API key not found"}

	// Validation errors
	CodeValidationError = Code{Kind: KindValidation, Domain: DomainValidation, Code: "validation_error", MessageForEndUser: "Validation failed"}

	// BadRequest errors
	CodeBadRequest = Code{Kind: KindBadRequest, Domain: DomainEmpty, Code: "bad_request", MessageForEndUser: "Bad request"}

	// Unauthorized errors
	CodeUnauthorized = Code{Kind: KindUnauthorized, Domain: DomainAuth, Code: "unauthorized", MessageForEndUser: "Unauthorized"}

	// Forbidden errors
	CodeForbidden = Code{Kind: KindForbidden, Domain: DomainAuth, Code: "forbidden", MessageForEndUser: "Forbidden"}

	// Conflict errors
	CodeConflict               = Code{Kind: KindConflict, Domain: DomainEmpty, Code: "conflict", MessageForEndUser: "Resource conflict"}
	CodeTaskVersionNotApproved = Code{Kind: KindConflict, Domain: DomainTask, Code: "task_version_not_approved", MessageForEndUser: "Task version is not approved for data collection"}

	// Internal errors
	CodeInternal      = Code{Kind: KindInternal, Domain: DomainEmpty, Code: "internal_error", MessageForEndUser: "Internal server error"}
	CodeDatabaseError = Code{Kind: KindInternal, Domain: DomainDatabase, Code: "database_error", MessageForEndUser: "Database error"}
	CodeRedisError    = Code{Kind: KindInternal, Domain: DomainEmpty, Code: "redis_error", MessageForEndUser: "Redis error"}
)
