package validator

import "time"

// CheckConfig holds the configuration for URL checking.
type CheckConfig struct {
	MaxRedirects          int
	MaxSize               int64
	MaxURLLength          int
	CheckForFile          bool
	CheckIsReachable      bool
	CannotEndWithSlash    bool
	HTTPClientTimeout     time.Duration
	HTTPSRequired         bool
	ContentTypeMustBeHTML bool
}
