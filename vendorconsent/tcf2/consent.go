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
	consentStringTCF2Separator         = '.'
	consentStringTCF2Prefix            = 'C'
	segmentTypeDisclosedVendors uint8  = 1
)

// ParseString parses the TCF 2.0 vendor string base64 encoded
func ParseString(consent string) (api.VendorConsents, error) {
	if consent == "" {
		return nil, consentconstants.ErrEmptyDecodedConsent
	}

	segments := strings.Split(consent, string(consentStringTCF2Separator))

	// Decode and parse the Core String (first segment)
	buff := []byte(segments[0])
	decoded := buff
	n, err := base64.RawURLEncoding.Decode(decoded, buff)
	if err != nil {
		return nil, err
	}
	decoded = decoded[:n:n]

	metadata, err := parseCoreString(decoded)
	if err != nil {
		return nil, err
	}

	// Parse additional segments
	for _, seg := range segments[1:] {
		segBuff := []byte(seg)
		segDecoded := segBuff
		sn, err := base64.RawURLEncoding.Decode(segDecoded, segBuff)
		if err != nil {
			continue
		}
		segDecoded = segDecoded[:sn:sn]
		if len(segDecoded) == 0 {
			continue
		}

		// First 3 bits are the segment type
		segType := segDecoded[0] >> 5
		if segType == segmentTypeDisclosedVendors {
			dv, err := parseDisclosedVendors(segDecoded)
			if err != nil {
				continue
			}
			metadata.disclosedVendors = dv
		}
	}

	return metadata, nil
}

// Parse parses the TCF 2.0 vendor consent data from the string. This string should *not* be encoded (by base64 or any other encoding).
// If the data is malformed and cannot be interpreted as a vendor consent string, this will return an error.
func Parse(data []byte) (api.VendorConsents, error) {
	metadata, err := parseCoreString(data)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}

// parseCoreString parses the core string data and returns the populated ConsentMetadata.
func parseCoreString(data []byte) (ConsentMetadata, error) {
	metadata, err := parseMetadata(data)
	if err != nil {
		return ConsentMetadata{}, err
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
		return ConsentMetadata{}, err
	}

	metadata.vendorConsents = vendorConsents
	metadata.vendorLegitimateInterestStart = legitIntStart + 17
	legIntMaxVend, err := bitutils.ParseUInt16(data, legitIntStart)
	if err != nil {
		return ConsentMetadata{}, err
	}

	if legitIntStart+16 >= uint(len(data))*8 {
		return ConsentMetadata{}, fmt.Errorf("invalid consent data: no legitimate interest start position")
	}
	if isSet(data, legitIntStart+16) {
		vendorLegitInts, pubRestrictsStart, err = parseRangeSection(metadata, legIntMaxVend, metadata.vendorLegitimateInterestStart)
	} else {
		vendorLegitInts, pubRestrictsStart, err = parseBitField(metadata, legIntMaxVend, metadata.vendorLegitimateInterestStart)
	}
	if err != nil {
		return ConsentMetadata{}, err
	}

	metadata.vendorLegitimateInterests = vendorLegitInts
	metadata.pubRestrictionsStart = pubRestrictsStart

	pubRestrictions, _, err := parsePubRestriction(metadata, pubRestrictsStart)
	if err != nil {
		return ConsentMetadata{}, err
	}

	metadata.publisherRestrictions = pubRestrictions

	return metadata, nil
}

// parseDisclosedVendors parses a Disclosed Vendors segment.
func parseDisclosedVendors(data []byte) (vendorConsentsResolver, error) {
	// Minimum: 3 bits (SegmentType) + 16 bits (MaxVendorId) + 1 bit (IsRangeEncoding) = 20 bits = 3 bytes
	if len(data) < 3 {
		return nil, fmt.Errorf("disclosed vendors segment requires at least 3 bytes. Got %d", len(data))
	}

	maxVendorID, err := bitutils.ParseUInt16(data, 3)
	if err != nil {
		return nil, err
	}

	if isSet(data, 19) {
		return parseSegmentRangeSection(data, maxVendorID, 20)
	}
	segMetadata := ConsentMetadata{data: data}
	resolver, _, err := parseBitField(segMetadata, maxVendorID, 20)
	return resolver, err
}

// parseSegmentRangeSection parses a range-encoded vendor section from segment data.
// Unlike parseRangeSection, it does not require a minimum of 31 bytes for the full core string.
func parseSegmentRangeSection(data []byte, maxVendorID uint16, startbit uint) (*rangeSection, error) {
	numEntries, err := bitutils.ParseUInt12(data, startbit)
	if err != nil {
		return nil, err
	}

	currentOffset := startbit + 12
	consents := make([]rangeConsent, numEntries)
	for i := range consents {
		bitsConsumed, err := parseRangeConsent(&consents[i], data, currentOffset, maxVendorID)
		if err != nil {
			return nil, err
		}
		currentOffset = currentOffset + bitsConsumed
	}

	return &rangeSection{
		consents:    consents,
		maxVendorID: maxVendorID,
	}, nil
}

// IsConsentV2 return true if the consent strings looks like a tcf v2 consent string
func IsConsentV2(consent string) bool {
	return len(consent) > 0 && consent[0] == consentStringTCF2Prefix
}
