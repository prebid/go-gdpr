package vendorconsent

import (
	"encoding/base64"
	"strings"

	"github.com/prebid/go-gdpr/api"
	"github.com/prebid/go-gdpr/consent1"
	"github.com/prebid/go-gdpr/consent2"
)

// ParseString parses a Raw (unpadded) base64 URL encoded string.
func ParseString(consent string) (api.VendorConsents, error) {
	pieces := strings.Split(consent, ".")
	decoded, err := base64.RawURLEncoding.DecodeString(pieces[0])
	if err != nil {
		return nil, err
	}
	version := uint8(decoded[0] >> 2)
	if version == 2 {
		return consent2.Parse(decoded)
	}
	return consent1.Parse(decoded)
}
