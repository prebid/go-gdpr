package vendorlist

import (
	"github.com/buger/jsonparser"
	"github.com/prebid/go-gdpr/consentconstants"
)

// ParseLazily20 returns a view of the data which re-calculates things on each function call.
// The returned object can be shared safely between goroutines.
//
// This is ideal if:
//   1. You only need to look up a few vendors or purpose IDs
//   2. You don't need good errors on malformed input
//
// Otherwise, you may get better performance with ParseEagerly20.
func ParseLazily20(data []byte) VendorList {
	return lazyVendorList20(data)
}

type lazyVendorList20 []byte

func (l lazyVendorList20) Version() uint16 {
	if val, ok := lazyParseInt(l, "vendorListVersion"); ok {
		return uint16(val)
	}
	return 0
}

func (l lazyVendorList20) Vendor(vendorID uint16) Vendor {
	var vendorBytes []byte
	jsonparser.ArrayEach(l, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if val, ok := lazyParseInt(value, "id"); ok {
			if uint16(val) == vendorID {
				vendorBytes = value
			}
		}
	}, "vendors")

	if len(vendorBytes) > 0 {
		return lazyVendor20(vendorBytes)
	}
	return nil
}

type lazyVendor20 []byte

func (l lazyVendor20) Purpose(purposeID consentconstants.Purpose) bool {
	exists := idExists(l, int(purposeID), "purposes")
	if exists {
		return true
	}
	return idExists(l, int(purposeID), "flexiblePurposes")
}

func (l lazyVendor20) LegitimateInterest(purposeID consentconstants.Purpose) bool {
	exists := idExists(l, int(purposeID), "legIntPurposes")
	if exists {
		return true
	}
	return idExists(l, int(purposeID), "flexiblePurposes")
}
