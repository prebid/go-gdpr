package vendorconsent

import (
	"testing"
)

func TestRangeSectionConsent(t *testing.T) {
	// String built using http://acdn.adnxs.com/cmp/docs/#/tools/vendor-cookie-encoder
	// This sample encodes a mix of Single- and Range-typed consent exceptions.
	//
	// The values which don't have parsing functions implemented yet are listed below for future tests.
	//
	// cookie version = 1
	// created = Tue May 08 2018 12:31:13 GMT-0400 (EDT) (binary 001110001101011100100010100000101110)
	// last updated = Tue May 08 2018 12:35:13 GMT-0400 (EDT) (binary 001110001101011100100011000110001010)
	// cmpId = 3
	// cmpVersion = 2
	// consentScreen = 7
	// consentLanguage = "en" (binary 000100001101)
	// purposeIdBitString = 001011010010110101101011
	consent, err := Parse(decode(t, "BONciguONcjGKADACHENAOLS1rAHDAFAAEAASABQAMwAeACEAFw"))
	if err != nil {
		t.Fatalf("Failed to parse valid consent string: %v", err)
	}
	assertUInt8sEqual(t, 1, consent.Version())
	assertUInt16sEqual(t, 14, consent.VendorListVersion())
	assertUInt16sEqual(t, 112, consent.MaxVendorID())

	var s struct{}
	lackingConsent := map[uint16]struct{}{
		2:  s,
		4:  s,
		10: s,
		11: s,
		12: s,
		13: s,
		14: s,
		15: s,
		16: s,
		17: s,
		18: s,
		19: s,
		20: s,
		21: s,
		22: s,
		23: s,
		24: s,
		25: s,
		30: s,
		31: s,
		32: s,
		33: s,
		46: s,
	}

	for i := uint16(1); i <= consent.MaxVendorID(); i++ {
		_, ok := lackingConsent[i]
		assertBoolsEqual(t, !ok, consent.HasConsent(i))
	}
}

func TestParseUInt16(t *testing.T) {
	// Start with 01100000 00000000 00100000
	// Expect 00000000 00000001
	doParseIntTest(t, []byte{0x60, 0x00, 0x20}, 3, 0x1)

	// Start with 00100000 00001110 00000000
	// Expect 00000000 01110000
	doParseIntTest(t, []byte{0x20, 0x0e, 0x00}, 3, 0x70)

	// Start with 11110100 00010011
	// Expect 11110100 00010011
	doParseIntTest(t, []byte{0xf4, 0x13}, 0, 0xf413)
}

func doParseIntTest(t *testing.T, data []byte, offset int, expected int) {
	t.Helper()
	parsedVal, err := parseUInt16(data, uint(offset))
	if err != nil {
		t.Fatalf("Error parsing uint16: %v", err)
	}
	if parsedVal != uint16(expected) {
		t.Errorf("Failed to parse value. Got %d", parsedVal)
	}
}
