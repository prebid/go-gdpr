package vendorconsent

import (
	"fmt"
	"testing"
)

func TestRangeSectionConsent(t *testing.T) {
	// String built using http://iabtcf.com/#/encode
	// This sample encodes a mix of Single- and Range-typed consent exceptions.
	consent, err := Parse(decode(t, "COyfVVoOyfVVoADACHENAwCAAAAAAAAAAAAAE5QBgALgAqgD8AQACSwEygJyAnSAMABgAFkAgQCDASeAmYBOgAA"))
	assertNilError(t, err)
	assertUInt8sEqual(t, 2, consent.Version())
	assertUInt16sEqual(t, 3, consent.CmpID())
	assertUInt16sEqual(t, 2, consent.CmpVersion())
	assertUInt8sEqual(t, 7, consent.ConsentScreen())
	assertStringsEqual(t, "EN", consent.ConsentLanguage())
	assertUInt16sEqual(t, 48, consent.VendorListVersion())
	assertUInt16sEqual(t, 626, consent.MaxVendorID())

	//  The above encoder doesn't support setting purposes.
	//	purposesWithConsent := buildMap(1, 3, 5, 6, 7, 10)
	//	for i := uint8(1); i <= 24; i++ {
	//		_, ok := purposesWithConsent[uint(i)]
	// 		assertBoolsEqual(t, ok, consent.PurposeAllowed(consentconstants.Purpose(i)))
	//	}

	vendorsWithConsent := buildMap(23, 42, 126, 127, 128, 587, 613, 626)
	for i := uint16(1); i <= consent.MaxVendorID(); i++ {
		_, expected := vendorsWithConsent[uint(i)]
		actual := consent.VendorConsent(i)

		if expected != actual {
			fmt.Printf("Vendor: %d failed\n", i)
		}

		assertBoolsEqual(t, expected, actual)
	}

	// TODO func VendorLegitInterest() should be added to api.VendorConsents
	consentMetadata := consent.(ConsentMetadata)
	vendorsLegitimateInterestWithConsent := buildMap(24, 44, 129, 130, 131, 591, 614, 628)

	for i := uint16(1); i <= consentMetadata.VendorLegitInterestMaxID(); i++ {
		_, expected := vendorsLegitimateInterestWithConsent[uint(i)]
		actual := consentMetadata.VendorLegitInterest(i)

		if expected != actual {
			fmt.Printf("VendorLegitInterest: %d failed\n", i)
		}

		assertBoolsEqual(t, expected, actual)
	}
}

// Prevents #10
func TestInvalidRangeEdgeCase(t *testing.T) {
	data := decode(t, "COwDzqZOwDzqZN4ABMENAPCAAP4AAP-AAAhoAFQAYABgAOABQAAAAA")
	data = data[:31]
	assertInvalidBytes(t, data[:31], "ParseUInt16 expected a 16-bit int to start at bit 243, but the consent string was only 31 bytes long")
}
