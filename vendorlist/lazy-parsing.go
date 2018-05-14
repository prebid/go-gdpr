package vendorlist

import (
	"strconv"

	"github.com/buger/jsonparser"
)

// ParseLazily returns a view of the data which re-calculates things on each function call.
//
// This is ideal if:
//   1. You only need to look up a few vendors or purpose IDs
//   2. You don't need good errors on malformed input
//
// Otherwise, you'd probably do better with a "ParseEagerly([]byte) (VendorList, error)"" function.
// PRs for this are welcome if it suits you!
func ParseLazily(data []byte) VendorList {
	return lazyVendorList(data)
}

type lazyVendorList []byte

func (l lazyVendorList) Version() uint16 {
	if val, ok := lazyParseInt(l, "vendorListVersion"); ok {
		return uint16(val)
	}
	return 0
}

func (l lazyVendorList) Vendor(vendorID uint16) Vendor {
	var vendorBytes []byte
	jsonparser.ArrayEach(l, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if val, ok := lazyParseInt(value, "id"); ok {
			if uint16(val) == vendorID {
				vendorBytes = value
			}
		}
	}, "vendors")

	if len(vendorBytes) > 0 {
		return lazyVendor(vendorBytes)
	}
	return nil
}

type lazyVendor []byte

func (l lazyVendor) Purpose(purposeID uint8) bool {
	return idExists(l, int(purposeID), "purposeIds")
}

func (l lazyVendor) LegitimateInterest(purposeID uint8) bool {
	return idExists(l, int(purposeID), "legIntPurposeIds")
}

// Returns false unless "id" exists in an array located at "data.key".
func idExists(data []byte, id int, key string) bool {
	hasID := false

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err == nil && dataType == jsonparser.Number {
			if intVal, err := strconv.ParseInt(string(value), 10, 0); err == nil {
				if int(intVal) == id {
					hasID = true
				}
			}
		}
	}, key)

	return hasID
}

func lazyParseInt(data []byte, key string) (int, bool) {
	if value, dataType, _, err := jsonparser.Get(data, key); err == nil && dataType == jsonparser.Number {
		intVal, err := strconv.Atoi(string(value))
		if err != nil {
			return 0, false
		}
		return intVal, true
	}
	return 0, false
}