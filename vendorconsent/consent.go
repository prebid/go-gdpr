package vendorconsent

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/prebid/go-gdpr/api"
	tcf1 "github.com/prebid/go-gdpr/vendorconsent/tcf1"
	tcf2 "github.com/prebid/go-gdpr/vendorconsent/tcf2"
)

const (
	consentStringTCF2Separator = '.'
)

var (
	errEmptyDecodedConsent = fmt.Errorf("decoded consent cannot be empty")
)

// ParseString parses a Raw (unpadded) base64 URL encoded string.
func ParseString(consent string) (api.VendorConsents, error) {
	// split TCF 2.0 segments
	if index := strings.IndexByte(consent, consentStringTCF2Separator); index != -1 {
		consent = consent[:index]
	}
	decoded, err := base64.RawURLEncoding.DecodeString(consent)
	if err != nil {
		return nil, err
	}
	version, err := ParseVersion(decoded)
	if err != nil {
		return nil, err
	}
	if version == 2 {
		return tcf2.Parse(decoded)
	}
	return tcf1.Parse(decoded)
}

// ParseVersion parses version from base64-decoded consent string
func ParseVersion(decodedConsent []byte) (uint8, error) {
	if len(decodedConsent) == 0 {
		return 0, errEmptyDecodedConsent
	}
	// read version from first 6 bits
	return decodedConsent[0] >> 2, nil
}

// Backwards compatibility

// VendorConsents old interface
type VendorConsents interface {
	api.VendorConsents
}

// Parse method is the equivalent of call tcf1.Parse(...)
func Parse(data []byte) (api.VendorConsents, error) {
	return tcf1.Parse(data)
}
