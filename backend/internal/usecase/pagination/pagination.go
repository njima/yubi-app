package pagination

const (
	DefaultLimit = 20
	MaxLimit     = 1000
)

// Params holds validated pagination parameters.
type Params struct {
	Page  int
	Limit int
}

// Parse extracts and validates pagination parameters from nullable query params.
// It applies default values and caps the limit to MaxLimit.
func Parse(page, limit *int) Params {
	p := Params{
		Page:  1,
		Limit: DefaultLimit,
	}

	if page != nil && *page > 0 {
		p.Page = *page
	}
	if limit != nil && *limit > 0 {
		p.Limit = *limit
		if p.Limit > MaxLimit {
			p.Limit = MaxLimit
		}
	}

	return p
}
