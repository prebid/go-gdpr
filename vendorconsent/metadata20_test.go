package vendorconsent

import (
	"testing"
	"time"
)

func TestCreatedDate20(t *testing.T) {
	consent, err := Parse20(decode(t, "COvcSpYOvcSpYC9AAAENAPCAAAAAAAAAAAAACvwDQABAAIAAYABIAC4AJQAagA9ACEAPgAjIBJoCvAK-AAAAAA"))
	assertNilError(t, err)
	created := consent.Created().UTC()
	year, month, day := created.Date()
	assertIntsEqual(t, 2020, year)
	assertIntsEqual(t, int(time.February), int(month))
	assertIntsEqual(t, 27, day)
	assertIntsEqual(t, 19, created.Hour())
	assertIntsEqual(t, 51, created.Minute())
	assertIntsEqual(t, 49, created.Second())
}

func TestLastUpdated20(t *testing.T) {
	consent, err := Parse20(decode(t, "COvcSpYOvcSpYC9AAAENAPCAAAAAAAAAAAAACvwDQABAAIAAYABIAC4AJQAagA9ACEAPgAjIBJoCvAK-AAAAAA"))
	assertNilError(t, err)
	updated := consent.LastUpdated().UTC()
	year, month, day := updated.Date()
	assertIntsEqual(t, 2020, year)
	assertIntsEqual(t, int(time.February), int(month))
	assertIntsEqual(t, 27, day)
	assertIntsEqual(t, 19, updated.Hour())
	assertIntsEqual(t, 51, updated.Minute())
	assertIntsEqual(t, 49, updated.Second())
}

func TestLargeVendorListVersion20(t *testing.T) {
	consent, err := Parse(decode(t, "BON96hFON96hFABABBAA4yAAAAAAEA"))
	assertNilError(t, err)
	assertUInt16sEqual(t, 3634, consent.VendorListVersion())
}

func TestLargeCmpID20(t *testing.T) {
	consent, err := Parse(decode(t, "BOOG3gbOOG3gbFZABBAAABAAAAAAEA"))
	assertNilError(t, err)
	assertUInt16sEqual(t, 345, consent.CmpID())
}

func TestLargeCmpVersion20(t *testing.T) {
	consent, err := Parse(decode(t, "BOOG4uyOOG4uyABFZBAAABAAAAAAEA"))
	assertNilError(t, err)
	assertUInt16sEqual(t, 345, consent.CmpVersion())
}

func TestLargeConsentScreen20(t *testing.T) {
	consent, err := Parse(decode(t, "BOOG4uyOOG4uyABFZTAAABAAAAAAEA"))
	assertNilError(t, err)
	assertUInt8sEqual(t, 19, consent.ConsentScreen())
}

func TestLanguageExtremes20(t *testing.T) {
	consent, err := Parse(decode(t, "BOOG9-6OOG9-6ABABBAZABAAAAAAEA"))
	assertNilError(t, err)
	assertStringsEqual(t, "AZ", consent.ConsentLanguage())

	consent, err = Parse(decode(t, "BOOG9-6OOG9-6ABABBZAABAAAAAAEA"))
	assertNilError(t, err)
	assertStringsEqual(t, "ZA", consent.ConsentLanguage())
}
