package validator

import "time"

// CheckConfig holds the configuration for URL checking.
type CheckConfig struct {
	MaxRedirects          int           `koanf:"max_redirects"`
	MaxSize               int64         `koanf:"max_size"`
	MaxURLLength          int           `koanf:"max_url_length"`
	CheckForFile          bool          `koanf:"check_for_file"`
	CheckIsReachable      bool          `koanf:"check_is_reachable"`
	CannotEndWithSlash    bool          `koanf:"cannot_end_with_slash"`
	HTTPClientTimeout     time.Duration `koanf:"http_client_timeout"`
	HTTPSRequired         bool          `koanf:"https_required"`
	ContentTypeMustBeHTML bool          `koanf:"content_type_must_be_html"`
}
