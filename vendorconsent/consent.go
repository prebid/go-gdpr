package vendorconsent

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/prebid/go-gdpr/api"
	tcf1 "github.com/prebid/go-gdpr/vendorconsent/tcf1"
	tcf2 "github.com/prebid/go-gdpr/vendorconsent/tcf2"
)

// ParseString parses a Raw (unpadded) base64 URL encoded string.
func ParseString(consent string) (api.VendorConsents, error) {
	pieces := strings.Split(consent, ".")
	decoded, err := base64.RawURLEncoding.DecodeString(pieces[0])
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
		return 0, fmt.Errorf("decoded consent cannot be empty")
	}
	// read version from first 6 bits
	return decodedConsent[0] >> 2, nil
}

// Backwards compatibility

type VendorConsents interface {
	api.VendorConsents
}

func Parse(data []byte) (api.VendorConsents, error) {
	return tcf1.Parse(data)
}
