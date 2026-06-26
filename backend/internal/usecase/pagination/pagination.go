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

// Normalized holds pagination parameters ready for repository calls.
type Normalized struct {
	Page   int
	Limit  int
	Offset int
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

// Normalize validates integer pagination parameters and computes the offset.
func Normalize(page, limit int) Normalized {
	return NormalizeWithDefault(page, limit, DefaultLimit)
}

// NormalizeWithDefault validates integer pagination parameters using defaultLimit
// when limit is not positive, then computes the offset.
func NormalizeWithDefault(page, limit, defaultLimit int) Normalized {
	if defaultLimit <= 0 {
		defaultLimit = DefaultLimit
	}

	p := Params{
		Page:  1,
		Limit: defaultLimit,
	}

	if page > 0 {
		p.Page = page
	}
	if limit > 0 {
		p.Limit = limit
		if p.Limit > MaxLimit {
			p.Limit = MaxLimit
		}
	}

	return Normalized{
		Page:   p.Page,
		Limit:  p.Limit,
		Offset: (p.Page - 1) * p.Limit,
	}
}
