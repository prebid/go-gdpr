package consent2

import (
	"testing"

	"github.com/prebid/go-gdpr/consentconstants"
)

func TestBitField(t *testing.T) {
	// String built using http://gdpr-demo.labs.quantcast.com/user-examples/cookie-workshop.html
	// This sample includes a BitField.
	consent, err := Parse(decode(t, "COwGVJOOwGVJOADACHENAOCAAO6as_-AAAhoAFNLAAoAAAA"))
	assertNilError(t, err)
	assertUInt8sEqual(t, 2, consent.Version())
	assertUInt16sEqual(t, 3, consent.CmpID())
	assertUInt16sEqual(t, 2, consent.CmpVersion())
	assertUInt8sEqual(t, 7, consent.ConsentScreen())
	assertStringsEqual(t, "EN", consent.ConsentLanguage())
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
