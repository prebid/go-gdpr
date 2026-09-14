package vendorconsent

import (
	"testing"

	"github.com/prebid/go-gdpr/consentconstants"
)

func TestRangeSectionConsent(t *testing.T) {
	// String built using http://acdn.adnxs.com/cmp/docs/#/tools/vendor-cookie-encoder
	// This sample encodes a mix of Single- and Range-typed consent exceptions.
	consent, err := Parse(decode(t, "BONciguONcjGKADACHENAOLS1rAHDAFAAEAASABQAMwAeACEAFw"))
	assertNilError(t, err)
	assertUInt8sEqual(t, 1, consent.Version())
	assertUInt16sEqual(t, 3, consent.CmpID())
	assertUInt16sEqual(t, 2, consent.CmpVersion())
	assertUInt8sEqual(t, 7, consent.ConsentScreen())
	assertStringsEqual(t, "EN", consent.ConsentLanguage())
	assertUInt16sEqual(t, 14, consent.VendorListVersion())
	assertUInt16sEqual(t, 112, consent.MaxVendorID())

	purposesWithConsent := buildMap(3, 5, 6, 8, 11, 13, 14, 16, 18, 19, 21, 23, 24)
	for i := uint8(1); i <= 24; i++ {
		_, ok := purposesWithConsent[uint(i)]
		assertBoolsEqual(t, ok, consent.PurposeAllowed(consentconstants.Purpose(i)))
	}

	vendorsLackingConsent := buildMap(2, 4, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 30, 31, 32, 33, 46)
	for i := uint16(1); i <= consent.MaxVendorID(); i++ {
		_, ok := vendorsLackingConsent[uint(i)]
		assertBoolsEqual(t, !ok, consent.VendorConsent(i))
	}
}

// Prevents #10
func TestInvalidRangeEdgeCase(t *testing.T) {
	data := decode(t, "BOQA9AtOQA9AtABABBAAABAAAAAGSAHAACAAMAAoABwAEgALAAaA")
	data = data[:36]
	assertInvalidBytes(t, data[:36], "bit 288 was supposed to start a new RangeEntry, but the consent string was only 36 bytes long")
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

func doParseIntTest(t *testing.T, data []byte, offset, expected int) {
	t.Helper()

	parsedVal, err := parseUInt16(data, uint(offset))
	assertNilError(t, err)

	if parsedVal != uint16(expected) {
		t.Errorf("Failed to parse value. Got %d", parsedVal)
	}
}
