package vendorconsent

import (
	"encoding/base64"
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

// Segment types defined in TCF 2.x specification.
// https://github.com/InteractiveAdvertisingBureau/GDPR-Transparency-and-Consent-Framework/blob/master/TCFv2/IAB%20Tech%20Lab%20-%20Consent%20string%20and%20vendor%20list%20formats%20v2.md#publisher-purposes-transparency-and-consent
// Only `SegmentTypeDisclosedVendors` is used in this file, but all types are included for specification completeness.
const (
	SegmentTypeCoreString       = 0
	SegmentTypeDisclosedVendors = 1
	SegmentTypePublisherTC      = 3
)

// ParseString parses the TCF 2.0 vendor string base64 encoded
func ParseString(consent string) (api.VendorConsents, error) {
	if consent == "" {
		return nil, consentconstants.ErrEmptyDecodedConsent
	}

	consentMeta, err := parseCoreAndDisclosedVendors(consent)
	if err != nil {
		return nil, err
	}

	return consentMeta, nil
}

// Parse parses the TCF 2.0 "Core string" segment. This string should *not* be encoded (by base64 or any other encoding).
// If the data is malformed and cannot be interpreted as a vendor consent string, this will return an error.
func Parse(data []byte) (api.VendorConsents, error) {
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
		vendorConsents, legitIntStart, err = parseRangeSection(metadata, metadata.MaxVendorID(), 230)
	} else {
		vendorConsents, legitIntStart, err = parseBitField(metadata, metadata.MaxVendorID(), 230)
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
		vendorLegitInts, pubRestrictsStart, err = parseRangeSection(metadata, legIntMaxVend, metadata.vendorLegitimateInterestStart)
	} else {
		vendorLegitInts, pubRestrictsStart, err = parseBitField(metadata, legIntMaxVend, metadata.vendorLegitimateInterestStart)
	}
	if err != nil {
		return nil, err
	}

	metadata.vendorLegitimateInterests = vendorLegitInts
	metadata.pubRestrictionsStart = pubRestrictsStart

	pubRestrictions, _, err := parsePubRestriction(metadata, pubRestrictsStart)
	if err != nil {
		return nil, err
	}

	metadata.publisherRestrictions = pubRestrictions

	return metadata, err
}

func parseCoreAndDisclosedVendors(consent string) (ConsentMetadata, error) {
	// Split TCF 2.0 segments by '.'
	// Format: [Core String].[Disclosed Vendors].[Publisher TC]
	segments := strings.Split(consent, string(consentStringTCF2Separator))

	// Parse the core string (always first segment)
	coreSegmentDecoded, err := decodeSegment(segments[0])
	if err != nil {
		return ConsentMetadata{}, err
	}

	// Parse the core string
	result, err := Parse(coreSegmentDecoded)
	if err != nil {
		return ConsentMetadata{}, err
	}

	metadata := result.(ConsentMetadata)

	// Parse disclosed vendors segment if present (TCF 2.3+)
	// Iterate through segments to find disclosed vendors by type (segments after Core String segment can be in any order)
	for _, segment := range segments[1:] {
		if segment == "" {
			continue
		}

		decoded, err := decodeSegment(segment)
		if err != nil {
			return ConsentMetadata{}, err
		}

		segmentType, err := getSegmentType(decoded)
		if err != nil {
			return ConsentMetadata{}, err
		}

		if segmentType == SegmentTypeDisclosedVendors { // Disclosed Vendors segment
			disclosedVendors, err := parseDisclosedVendorsSegment(decoded)
			if err != nil {
				return ConsentMetadata{}, fmt.Errorf("failed to parse disclosed vendors segment: %v", err)
			}
			metadata.disclosedVendors = disclosedVendors
			metadata.hasDisclosedVendors = true
			break
		}
	}

	return metadata, nil
}

// IsConsentV2 return true if the consent strings looks like a tcf v2 consent string
func IsConsentV2(consent string) bool {
	return len(consent) > 0 && consent[0] == consentStringTCF2Prefix
}

// decodeSegment decodes a base64 encoded segment string.
func decodeSegment(segmentString string) ([]byte, error) {
	if segmentString == "" {
		return nil, fmt.Errorf("empty segment string")
	}

	buff := []byte(segmentString)
	decoded := buff
	n, err := base64.RawURLEncoding.Decode(decoded, buff)
	if err != nil {
		return nil, fmt.Errorf("failed to decode segment: %v", err)
	}

	return decoded[:n:n], nil
}

// getSegmentType extracts the 3-bit segment type from the segment data
func getSegmentType(data []byte) (uint8, error) {
	if len(data) < 1 {
		return 0, fmt.Errorf("segment too short")
	}

	segmentType := data[0] >> 5
	return segmentType, nil
}
