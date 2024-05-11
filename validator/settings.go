package validator

import (
	"time"
)

var CheckErrorMessages = map[CheckErrorMSG]string{
	ErrInvalidSchema:         "Invalid URL schema. Only 'http' or 'https' are allowed.",
	ErrFileNotAllowed:        "Files are not allowed.",
	ErrMaxRedirectsExceeded:  "Maximum number of redirects exceeded.",
	ErrSizeExceeded:          "URL size exceeds the maximum allowed size.",
	ErrURLLengthExceeded:     "URL length exceeds the maximum allowed length.",
	ErrUnreachable:           "URL is unreachable.",
	ErrMustHaveTLD:           "URL must have a valid top-level domain (TLD).",
	ErrURLCannotContainCreds: "URL should not contain credentials.",
	ErrURLCannotEndWithSlash: "URL should not end with a slash.",
	ErrHTTPClientTimeout:     "Website is taking too long to respond",
	ErrHTTPSRequired:         "HTTPS is required.",
	ErrDNSRecordNotFound:     "DNS record not found.",
	ErrSeekerNotInitialized:  "Seeker is not initialized.",
	ErrIsNotAbsoluteURL:      "URL is not absolute.",
	ErrContentTypeNotAllowed: "Content type is not allowed.",
	ErrContentTypeIsEmpty:    "Content type is empty.",
	ErrDNSNameValueNull:      "DNS name value is null.",
}

var DefaultCheckConfig = CheckConfig{
	MaxRedirects:          2,       // Set your desired default values
	MaxSize:               4194304, // 4MB
	MaxURLLength:          2048,    // 2048 characters
	CheckForFile:          true,
	CheckIsReachable:      true,
	CannotEndWithSlash:    true,
	HTTPSRequired:         false,
	HTTPClientTimeout:     10 * time.Second,
	ContentTypeMustBeHTML: true,
}
