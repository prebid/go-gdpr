package vendorconsent

import (
	"testing"
	"time"
)

func TestCreatedDate(t *testing.T) {
	consent, err := Parse(decode(t, "BIRAfK8OOHsDFABABBAAABAAAAAAEA"))
	assertNilError(t, err)

	created := consent.Created().UTC()
	year, month, day := created.Date()
	assertIntsEqual(t, 1998, year)
	assertIntsEqual(t, int(time.February), int(month))
	assertIntsEqual(t, 15, day)
	assertIntsEqual(t, 7, created.Hour())
	assertIntsEqual(t, 24, created.Minute())
	assertIntsEqual(t, 54, created.Second())
}

func TestLastUpdated(t *testing.T) {
	consent, err := Parse(decode(t, "BIRAfK8OOHsDFABABBAAABAAAAAAEA"))
	assertNilError(t, err)

	updated := consent.LastUpdated().UTC()
	year, month, day := updated.Date()
	assertIntsEqual(t, 2018, year)
	assertIntsEqual(t, int(time.May), int(month))
	assertIntsEqual(t, 21, day)
	assertIntsEqual(t, 18, updated.Hour())
	assertIntsEqual(t, 43, updated.Minute())
	assertIntsEqual(t, 18, updated.Second())
}

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
