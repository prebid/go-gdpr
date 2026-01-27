package vendorconsent

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/prebid/go-gdpr/api"
	"github.com/prebid/go-gdpr/bitutils"
	"github.com/prebid/go-gdpr/consentconstants"
)

const (
	consentStringTCF2Separator = '.'
	consentStringTCF2Prefix    = 'C'
)

// ParseString parses the TCF 2.0 vendor string base64 encoded
func ParseString(consent string) (api.VendorConsents, error) {
	if consent == "" {
		return nil, consentconstants.ErrEmptyDecodedConsent
	}

	buff := []byte(consent)

	// Decode core string (first segment, required)
	firstDot := strings.IndexByte(consent, consentStringTCF2Separator)
	coreEnd := len(consent)
	if firstDot != -1 {
		coreEnd = firstDot
	}
	writePos, err := base64.RawURLEncoding.Decode(buff, buff[:coreEnd])
	if err != nil {
		return nil, fmt.Errorf("failed to decode core segment: %w", err)
	}

	// Decode additional segments, looking for Disclosed Vendors
	disclosedOffset := 0
	if firstDot != -1 {
		pos := firstDot
		for pos < len(consent) {
			pos++ // skip '.'
			nextDot := strings.IndexByte(consent[pos:], consentStringTCF2Separator)
			segEnd := len(consent)
			if nextDot != -1 {
				segEnd = pos + nextDot
			}

			n, err := base64.RawURLEncoding.Decode(buff[writePos:], buff[pos:segEnd])
			if err != nil {
				return nil, fmt.Errorf("failed to decode segment: %w", err)
			}

			if n > 0 && (buff[writePos]>>5) == 1 { // segment type 1 = Disclosed Vendors
				disclosedOffset = writePos
			}
			writePos += n
			pos = segEnd
		}
	}

	metadata, err := parseCore(buff[:writePos])
	if err != nil {
		return nil, err
	}

	if disclosedOffset > 0 {
		metadata.disclosedVendors, err = parseVendorSection(buff[:writePos], uint(disclosedOffset)*8+3)
		if err != nil {
			return nil, err
		}
	}

	// TCF 2.3+ (policy version >= 5) requires Disclosed Vendors segment
	if metadata.TCFPolicyVersion() >= 5 && metadata.disclosedVendors == nil {
		return nil, errors.New("TCF 2.3+ requires Disclosed Vendors segment")
	}

	return metadata, nil
}

// Parse parses the TCF 2.0 vendor consent data from the bytes. This data should *not* be encoded (by base64 or any other encoding).
// If the data is malformed and cannot be interpreted as a vendor consent string, this will return an error.
// Note: This parses the core string only. Use ParseString to parse consent strings with additional segments (e.g., Disclosed Vendors).
func Parse(data []byte) (api.VendorConsents, error) {
	return parseCore(data)
}

// parseCore parses the core TCF 2.0 consent data and returns a pointer to allow modification.
func parseCore(data []byte) (*ConsentMetadata, error) {
	metadata, err := parseMetadata(data)
	if err != nil {
		return nil, err
	}

	var vendorConsents vendorConsentsResolver
	var vendorLegitInts vendorConsentsResolver

	var legitIntStart uint
	var pubRestrictsStart uint
	// Bit 229 determines whether or not the consent string encodes Vendor data in a RangeSection or BitField.
	// We know from parseMetadata that we have at least 29*8=232 bits available
	if isSet(data, 229) {
		vendorConsents, legitIntStart, err = parseRangeSection(data, metadata.MaxVendorID(), 230)
	} else {
		vendorConsents, legitIntStart, err = parseBitField(data, metadata.MaxVendorID(), 230)
	}
	if err != nil {
		return nil, err
	}

	metadata.vendorConsents = vendorConsents
	metadata.vendorLegitimateInterestStart = legitIntStart + 17
	legIntMaxVend, err := bitutils.ParseUInt16(data, legitIntStart)
	if err != nil {
		return nil, err
	}

	if legitIntStart+16 >= uint(len(data))*8 {
		return nil, fmt.Errorf("invalid consent data: no legitimate interest start position")
	}
	if isSet(data, legitIntStart+16) {
		vendorLegitInts, pubRestrictsStart, err = parseRangeSection(data, legIntMaxVend, metadata.vendorLegitimateInterestStart)
	} else {
		vendorLegitInts, pubRestrictsStart, err = parseBitField(data, legIntMaxVend, metadata.vendorLegitimateInterestStart)
	}
	if err != nil {
		return nil, err
	}

	metadata.vendorLegitimateInterests = vendorLegitInts
	metadata.pubRestrictionsStart = pubRestrictsStart

	pubRestrictions, _, err := parsePubRestriction(data, pubRestrictsStart)
	if err != nil {
		return nil, err
	}

	metadata.publisherRestrictions = pubRestrictions

	return &metadata, nil
}

// parseVendorSection parses a vendor section (used for Disclosed Vendors segment).
// startBit should point to the MaxVendorId field (after any segment header).
func parseVendorSection(data []byte, startBit uint) (vendorConsentsResolver, error) {
	maxVendorID, err := bitutils.ParseUInt16(data, startBit)
	if err != nil {
		return nil, err
	}

	// EncodingType is 1 bit after MaxVendorID
	isRange := isSet(data, startBit+16)

	// Payload starts after MaxVendorID (16 bits) + EncodingType (1 bit)
	payloadStart := startBit + 17

	if isRange {
		rs, _, err := parseRangeSection(data, maxVendorID, payloadStart)
		return rs, err
	}

	bf, _, err := parseBitField(data, maxVendorID, payloadStart)
	return bf, err
}

// IsConsentV2 return true if the consent strings looks like a tcf v2 consent string
func IsConsentV2(consent string) bool {
	return len(consent) > 0 && consent[0] == consentStringTCF2Prefix
}
