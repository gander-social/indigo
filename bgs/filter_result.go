package bgs

// FilterResult represents the result of filtering an event
type FilterResult struct {
	Pass     bool   `json:"pass"`
	Reason   string `json:"reason"`
	HashOnly bool   `json:"hash_only,omitempty"`
}

// NewFilterResult creates a new filter result
func NewFilterResult(pass bool, reason string) *FilterResult {
	return &FilterResult{
		Pass:   pass,
		Reason: reason,
	}
}

// NewHashOnlyFilterResult creates a filter result that only allows hash
func NewHashOnlyFilterResult(reason string) *FilterResult {
	return &FilterResult{
		Pass:     false,
		Reason:   reason,
		HashOnly: true,
	}
}
