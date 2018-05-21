package vendorconsent

import (
	"testing"

	"github.com/prebid/go-gdpr/consentconstants"
)

func TestBitField(t *testing.T) {
	// String built using http://acdn.adnxs.com/cmp/docs/#/tools/vendor-cookie-encoder
	// This sample includes a BitField.
	//
	// The values which don't have parsing functions implemented yet are listed below for future tests.
	//
	// cookie version = 1
	// created = Sun May 06 2018 12:31:13 GMT-0400 (EDT) (binary 001110001101010101111100101000101010)
	// last updated = Mon May 07 2018 01:42:15 GMT-0400 (EDT) (binary 001110001101010111110000100000100110)
	// consentScreen = 7
	// consentLanguage = "en" (binary 000100001101)
	consent, err := Parse(decode(t, "BONV8oqONXwgmADACHENAO7pqzAAppY"))
	assertNilError(t, err)
	assertUInt8sEqual(t, 1, consent.Version())
	assertUInt16sEqual(t, 3, consent.CmpID())
	assertUInt16sEqual(t, 2, consent.CmpVersion())
	assertUInt16sEqual(t, 14, consent.VendorListVersion())
	assertUInt16sEqual(t, 10, consent.MaxVendorID())

	purposesAllowed := buildMap(1, 2, 3, 5, 6, 7, 9, 12, 13, 15, 17, 19, 20, 23, 24)
	for i := uint8(1); i <= 24; i++ {
		_, ok := purposesAllowed[uint(i)]
		assertBoolsEqual(t, ok, consent.PurposeAllowed(consentconstants.Purpose(i)))
	}

	vendorsWithConsent := buildMap(1, 2, 4, 7, 9, 10)
	for i := uint16(1); i <= consent.MaxVendorID(); i++ {
		_, ok := vendorsWithConsent[uint(i)]
		assertBoolsEqual(t, ok, consent.VendorConsent(i))
	}
}
