package vendorconsent

import "testing"

func TestBitField(t *testing.T) {
	// String built using http://acdn.adnxs.com/cmp/docs/#/tools/vendor-cookie-encoder
	// This sample includes a BitField.
	//
	// The values which don't have parsing functions implemented yet are listed below for future tests.
	//
	// cookie version = 1
	// created = Sun May 06 2018 12:31:13 GMT-0400 (EDT) (binary 001110001101010101111100101000101010)
	// last updated = Mon May 07 2018 01:42:15 GMT-0400 (EDT) (binary 001110001101010111110000100000100110)
	// cmpId = 3
	// cmpVersion = 2
	// consentScreen = 7
	// consentLanguage = "en" (binary 000100001101)
	// purposeIdBitString = 111011101001101010110011
	consent, err := Parse(decode(t, "BONV8oqONXwgmADACHENAO7pqzAAppY"))
	if err != nil {
		t.Fatalf("Failed to parse valid consent string: %v", err)
	}
	assertUInt8sEqual(t, 1, consent.Version())
	assertUInt16sEqual(t, 14, consent.VendorListVersion())
	assertUInt16sEqual(t, 10, consent.MaxVendorID())

	var s struct{}
	hasConsent := map[uint16]struct{}{
		1:  s,
		2:  s,
		4:  s,
		7:  s,
		9:  s,
		10: s,
	}

	for i := uint16(1); i <= consent.MaxVendorID(); i++ {
		_, ok := hasConsent[i]
		assertBoolsEqual(t, ok, consent.HasConsent(i))
	}
}
