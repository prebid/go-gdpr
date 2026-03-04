package vendorconsent

import (
	"fmt"

	"github.com/prebid/go-gdpr/bitutils"
)

// parseDisclosedVendorsSegment parses the Disclosed Vendors segment (SegmentType=1).
// This segment is mandatory in TCF 2.3.
func parseDisclosedVendorsSegment(data []byte) (vendorConsentsResolver, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data is empty")
	}

	// Need at least 3 bits for segment type + 16 bits for MaxVendorId + 1 bit for IsRangeEncoding
	if len(data) < 3 {
		return nil, fmt.Errorf("segment too short: %d bytes, need at least 3", len(data))
	}

	segmentType, err := bitutils.ParseByte8(data, 0)
	if err != nil {
		return nil, fmt.Errorf("parse segment type: %v", err)
	}
	segmentType = segmentType >> 5 // Get first 3 bits

	if segmentType != SegmentTypeDisclosedVendors {
		return nil, fmt.Errorf("expected segment type 1, got %d", segmentType)
	}

	maxVendorID, err := bitutils.ParseUInt16(data, 3)
	if err != nil {
		return nil, fmt.Errorf("parse MaxVendorId: %v", err)
	}

	// IsRangeEncoding is at bit 19 (0-based indexing)
	// see https://github.com/InteractiveAdvertisingBureau/GDPR-Transparency-and-Consent-Framework/blob/master/TCFv2/IAB%20Tech%20Lab%20-%20Consent%20string%20and%20vendor%20list%20formats%20v2.md#disclosed-vendors)
	isRangeEncoding := isSet(data, 19)

	// Create a temporary metadata just for parsing purposes
	tempMetadata := ConsentMetadata{data: data}

	if isRangeEncoding {
		rangeSection, _, err := parseRangeSection(tempMetadata, maxVendorID, 20)
		if err != nil {
			return nil, fmt.Errorf("parse range section: %v", err)
		}
		return rangeSection, nil
	}

	bitField, _, err := parseBitField(tempMetadata, maxVendorID, 20)
	if err != nil {
		return nil, fmt.Errorf("parse bit field: %v", err)
	}
	return bitField, nil
}
