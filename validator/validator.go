package validator

import (
	"errors"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type Validator struct {
	config CheckConfig
}

func NewValidator(config CheckConfig) *Validator {
	return &Validator{config}
}

func DefaultValidator() *Validator {
	return &Validator{DefaultCheckConfig}
}

// CheckURL validates a URL based on the given configuration and returns custom errors.
func (v Validator) ValidateURL(inputURL string) (isValid bool, err error) {
	isValid = false
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return
	}

	if v.config.MaxURLLength != 0 && len(inputURL) > v.config.MaxURLLength {
		err = &CheckError{CheckErrorMessages[ErrURLLengthExceeded]}
		return
	}

	if !parsedURL.IsAbs() {
		err = &CheckError{CheckErrorMessages[ErrIsNotAbsoluteURL]}
		return
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		err = &CheckError{CheckErrorMessages[ErrInvalidSchema]}
		return
	}

	if v.config.HTTPSRequired && parsedURL.Scheme != "https" {
		err = &CheckError{CheckErrorMessages[ErrHTTPSRequired]}
		return
	}

	if !HasValidTLD(parsedURL.Host) {
		err = &CheckError{CheckErrorMessages[ErrMustHaveTLD]}
		return
	}

	// Check for username and password in the URL http://username:password@yoursitename.com
	if parsedURL.User != nil {
		err = &CheckError{CheckErrorMessages[ErrURLCannotContainCreds]}
		return
	}

	if v.config.CannotEndWithSlash && strings.HasPrefix(parsedURL.Path, "/") {
		err = &CheckError{CheckErrorMessages[ErrURLCannotEndWithSlash]}
		return
	}

	return true, nil
}

// Validate a Website via visiting the URL and checking the response according to the configuration.
func (v Validator) ValidateSite(inputURL string) (isValid bool, err error) {
	seeker := NewSeeker(v.config, nil)

	if seeker == nil {
		err = &CheckError{CheckErrorMessages[ErrSeekerNotInitialized]}
		return
	}

	_err := seeker.Seek(inputURL)

	if _err != nil {
		err = _err
		return
	}

	if seeker.Logs.ContentLength > v.config.MaxSize {
		err = &CheckError{CheckErrorMessages[ErrSizeExceeded]}
		return
	}

	if seeker.Logs.Redirects != nil && len(*seeker.Logs.Redirects) > v.config.MaxRedirects {
		err = &CheckError{CheckErrorMessages[ErrMaxRedirectsExceeded]}
		return
	}

	log.Printf("Seeker Took: %v\n", time.Since(seeker.Logs.StartAt))
	log.Printf("First Byte took: %v\n", seeker.Logs.FirstByteDuration)
	log.Printf("ValidateURL took: %v\n", seeker.Logs.Duration)

	for _, redirect := range *seeker.Logs.Redirects {
		log.Printf("Seeker Redirect: %v\n", redirect.URL)
	}

	if seeker.Logs.Content != nil {
		log.Printf("Seeker Content: %v\n", *seeker.Logs.Content)
	}

	if seeker.Logs.ContentType == "" {
		err = &CheckError{CheckErrorMessages[ErrContentTypeIsEmpty]}
		return
	}

	// check if the content type include Logs."text/html"
	if seeker.Config.ContentTypeMustBeHTML && !strings.Contains(seeker.Logs.ContentType, "text/html") {
		err = &CheckError{CheckErrorMessages[ErrContentTypeNotAllowed]}
		return
	}

	return true, nil
}

func (v Validator) ValidateOwnershipOverDNSTxtRecord(inputURL string, DNSName string, DNSValue string, DNSServer string) (isValid bool, err error) {
	isValid = false

	if DNSName == "" || DNSValue == "" {
		err = errors.New(CheckErrorMessages[ErrDNSNameValueNull])
		return
	}

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return
	}

	if !HasValidTLD(parsedURL.Host) {
		err = errors.New(CheckErrorMessages[ErrMustHaveTLD])
		return
	}

	recordKey := strings.Join([]string{DNSName, parsedURL.Host}, ".")

	records, err := fetchDNSRecords(DNSServer, recordKey, dns.TypeTXT)
	if err != nil {
		return
	}

	for _, record := range records {
		txtRecord, ok := record.(*dns.TXT)
		if ok && DNSValue == strings.Join(txtRecord.Txt, "") {
			isValid = true
			break
		}
	}

	if !isValid {
		err = errors.New(CheckErrorMessages[ErrDNSRecordNotFound])
		return
	}

	return
}

func (v Validator) seekIsNeeded() bool {
	return v.config.CheckIsReachable || v.config.CheckForFile || v.config.MaxSize > 0 || v.config.MaxRedirects > 0
}

func HasValidTLD(host string) bool {
	// Define a regular expression to match common TLD patterns.
	tldPattern := `(\.[a-zA-Z]{2,63})$`
	matched, err := regexp.MatchString(tldPattern, host)
	return matched && err == nil
}
