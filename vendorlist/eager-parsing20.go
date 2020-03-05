package vendorlist

import (
	"encoding/json"
	"errors"

	"github.com/prebid/go-gdpr/consentconstants"
)

// ParseEagerly interprets and validates the Vendor List data up front, before returning it.
// The returned object can be shared safely between goroutines.
//
// This is ideal if:
//   1. You plan to call functions on the returned VendorList many times before discarding it.
//   2. You need strong input validation and good error messages.
//
// Otherwise, you may get better performance with ParseLazily.
func ParseEagerly20(data []byte) (VendorList, error) {
	var contract vendorListContract20
	if err := json.Unmarshal(data, &contract); err != nil {
		return nil, err
	}

	if contract.Version == 0 {
		return nil, errors.New("data.vendorListVersion was 0 or undefined. Versions should start at 1")
	}
	if len(contract.Vendors) == 0 {
		return nil, errors.New("data.vendors was undefined or had no elements")
	}

	parsedList := parsedVendorList20{
		version: contract.Version,
		vendors: make(map[uint16]parsedVendor20, len(contract.Vendors)),
	}

	for i := 0; i < len(contract.Vendors); i++ {
		thisVendor := contract.Vendors[i]
		parsedList.vendors[thisVendor.ID] = parseVendor20(thisVendor)
	}

	return parsedList, nil
}

func parseVendor20(contract vendorListVendorContract20) parsedVendor20 {
	parsed := parsedVendor20{
		purposes:            mapify(contract.Purposes),
		legitimateInterests: mapify(contract.LegitimateInterests),
		flexiblePurposes:    mapify(contract.FlexiblePurposes),
	}

	return parsed
}

type parsedVendorList20 struct {
	version uint16
	vendors map[uint16]parsedVendor20
}

func (l parsedVendorList20) Version() uint16 {
	return l.version
}

func (l parsedVendorList20) Vendor(vendorID uint16) Vendor {
	vendor, ok := l.vendors[vendorID]
	if ok {
		return vendor
	}
	return nil
}

type parsedVendor20 struct {
	purposes            map[consentconstants.Purpose]struct{}
	legitimateInterests map[consentconstants.Purpose]struct{}
	flexiblePurposes    map[consentconstants.Purpose]struct{}
}

func (l parsedVendor20) Purpose(purposeID consentconstants.Purpose) (hasPurpose bool) {
	_, hasPurpose = l.purposes[purposeID]
	if !hasPurpose {
		_, hasPurpose = l.flexiblePurposes[purposeID]
	}
	return
}

// LegitimateInterest retursn true if this vendor claims a "Legitimate Interest" to
// use data for the given purpose.
//
// For an explanation of legitimate interest, see https://www.gdpreu.org/the-regulation/key-concepts/legitimate-interest/
func (l parsedVendor20) LegitimateInterest(purposeID consentconstants.Purpose) (hasLegitimateInterest bool) {
	_, hasLegitimateInterest = l.legitimateInterests[purposeID]
	if !hasLegitimateInterest {
		_, hasLegitimateInterest = l.flexiblePurposes[purposeID]
	}
	return
}

type vendorListContract20 struct {
	Version uint16                       `json:"vendorListVersion"`
	Vendors []vendorListVendorContract20 `json:"vendors"`
}

type vendorListVendorContract20 struct {
	ID                  uint16  `json:"id"`
	Purposes            []uint8 `json:"purposes"`
	LegitimateInterests []uint8 `json:"legIntPurposes"`
	FlexiblePurposes    []uint8 `json:"flexiblePurposes"`
}
