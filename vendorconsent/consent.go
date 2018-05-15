package vendorconsent

import "github.com/prebid/go-gdpr/consentconstants"

// VendorConsents is a GDPR Vendor Consent string, as defined by IAB Europe. For technical details,
// see https://github.com/InteractiveAdvertisingBureau/GDPR-Transparency-and-Consent-Framework/blob/master/Consent%20string%20and%20vendor%20list%20formats%20v1.1%20Final.md#vendor-consent-string-format-
type VendorConsents interface {
	// The version of the Consent string.
	Version() uint8

	// The VendorListVersion which is needed to interpret this consent string.
	//
	// The IAB is hosting these on their webpage. For example, version 2 of the
	// Vendor List can be found at https://vendorlist.consensu.org/v-2/vendorlist.json
	//
	// For other versions, just replace the "v-*" path with the value returned here.
	// The latest version can always be found at https://vendorlist.consensu.org/vendorlist.json
	VendorListVersion() uint16

	// MaxVendorID describes how many vendors are encoded into the string.
	// This is the upper bound (inclusive) on valid inputs for HasConsent(id).
	MaxVendorID() uint16

	// Determine if the user has consented to use data for the given Purpose.
	//
	// If the purpose is converted from an int > 24, the return value is undefined because
	// the consent string doesn't have room for more purposes than that.
	PurposeAllowed(id consentconstants.Purpose) bool

	// Determine if a given vendor has consent to collect or receive user info.
	//
	// This function's behavior is undefined for "invalid" IDs.
	// IDs with value < 1 or value > MaxVendorID() are definitely invalid, but IDs within that range
	// may still be invalid, depending on the Vendor List.
	//
	// It is the caller's responsibility to get the right Vendor List version for the semantics of the ID.
	// For more information, see VendorListVersion().
	VendorConsent(id uint16) bool
}

// Parse the vendor consent data from the string. This string should *not* be encoded (by base64 or any other encoding).
// If the data is malformed and cannot be interpreted as a vendor consent string, this will return an error.
func Parse(data []byte) (VendorConsents, error) {
	metadata, err := parseMetadata(data)
	if err != nil {
		return nil, err
	}

	// Bit 172 determines whether or not the consent string encodes Vendor data in a RangeSection or BitField.
	if isSet(data, 172) {
		return parseRangeSection(metadata)
	}

	return parseBitField(metadata)
}

// Returns true if the bitIndex'th bit in data is a 1, and false if it's a 0.
func isSet(data []byte, bitIndex uint) bool {
	byteIndex := bitIndex / 8
	bitOffset := bitIndex % 8
	return byteToBool(data[byteIndex] & (0x80 >> bitOffset))
}
