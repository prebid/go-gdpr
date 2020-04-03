package vendorconsent

import "github.com/prebid/go-gdpr/api"

// Parse parses the TCF 2.0 vendor consent data from the string. This string should *not* be encoded (by base64 or any other encoding).
// If the data is malformed and cannot be interpreted as a vendor consent string, this will return an error.
func Parse(data []byte) (api.VendorConsents, error) {
	metadata, err := parseMetadata(data)
	if err != nil {
		return nil, err
	}

	var vendorConsents vendorConsentsResolver
	var legIntStart uint
	// Bit 229 determines whether or not the consent string encodes Vendor data in a RangeSection or BitField.
	if isSet(data, 229) {
		vendorConsents, legIntStart, err = parseRangeSection(metadata)
	} else {
		vendorConsents, legIntStart, err = parseBitField(metadata)
	}

	metadata.vendorConsents = vendorConsents
	metadata.vendorLegitimateInterestStart = legIntStart

	if err != nil {
		return nil, err
	}
	return metadata, err

}
