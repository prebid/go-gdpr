package vendorconsent

import "testing"

func TestLargeVendorListVersion(t *testing.T) {
	consent, err := Parse(decode(t, "BON96hFON96hFABABBAA4yAAAAAAEA"))
	assertNilError(t, err)
	assertUInt16sEqual(t, 3634, consent.VendorListVersion())
}

func TestLargeCmpID(t *testing.T) {
	consent, err := Parse(decode(t, "BOOG3gbOOG3gbFZABBAAABAAAAAAEA"))
	assertNilError(t, err)
	assertUInt16sEqual(t, 345, consent.CmpID())
}

func TestLargeCmpVersion(t *testing.T) {
	consent, err := Parse(decode(t, "BOOG4uyOOG4uyABFZBAAABAAAAAAEA"))
	assertNilError(t, err)
	assertUInt16sEqual(t, 345, consent.CmpVersion())
}

func TestLargeConsentScreen(t *testing.T) {
	consent, err := Parse(decode(t, "BOOG4uyOOG4uyABFZTAAABAAAAAAEA"))
	assertNilError(t, err)
	assertUInt8sEqual(t, 19, consent.ConsentScreen())
}

func TestLanguageExtremes(t *testing.T) {
	consent, err := Parse(decode(t, "BOOG9-6OOG9-6ABABBAZABAAAAAAEA"))
	assertNilError(t, err)
	assertStringsEqual(t, "AZ", consent.ConsentLanguage())

	consent, err = Parse(decode(t, "BOOG9-6OOG9-6ABABBZAABAAAAAAEA"))
	assertNilError(t, err)
	assertStringsEqual(t, "ZA", consent.ConsentLanguage())
}

func assertNilError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
