package validator

type CheckErrorMSG int

type CheckError struct {
	Msg string
}

func (e CheckError) Error() string {
	return e.Msg
}

const (
	ErrInvalidSchema CheckErrorMSG = iota
	ErrFileNotAllowed
	ErrMaxRedirectsExceeded
	ErrSizeExceeded
	ErrURLLengthExceeded
	ErrLoginKeywordsFound
	ErrXSSDetected
	ErrUnreachable
	ErrMustHaveTLD
	ErrURLCannotContainCreds
	ErrURLCannotEndWithSlash
	ErrHTTPClientTimeout
	ErrHTTPSRequired
	ErrDNSRecordNotFound
	ErrSeekerNotInitialized
	ErrIsNotAbsoluteURL
	ErrContentTypeNotAllowed
	ErrContentTypeIsEmpty
	ErrDNSNameValueNull
)
