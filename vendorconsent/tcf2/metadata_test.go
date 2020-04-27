package vendorconsent

import (
	"testing"
	"time"
)

func TestCreatedDate(t *testing.T) {
	consent, err := Parse(decode(t, "COvcSpYOvcSpYC9AAAENAPCAAAAAAAAAAAAACvwDQABAAIAAYABIAC4AJQAagA9ACEAPgAjIBJoCvAK-AAAAAA"))
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

func TestLastUpdate(t *testing.T) {
	consent, err := Parse(decode(t, "COvcSpYOvcSpYC9AAAENAPCAAAAAAAAAAAAACvwDQABAAIAAYABIAC4AJQAagA9ACEAPgAjIBJoCvAK-AAAAAA"))
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

func TestLargeCmpID(t *testing.T) {
	consent, err := Parse(decode(t, "COyiBqdOyiBqdObAAAENAfCIAP8AAH-AAAAAB4AXQQgEAAAgoAAAAABAIYQUAAAAAAAAAAAAAAAIQIQCxIvkgQMAAAABgAIAAAAAAAAAAABAZAkAAA"))
	assertNilError(t, err)
	assertUInt16sEqual(t, 923, consent.CmpID())
}

func TestLargeCmpVersion(t *testing.T) {
	consent, err := Parse(decode(t, "COyiCPlOyiCPlKxMIAENAfCAAAAAAAAAAAAAAAAAAAAA"))
	assertNilError(t, err)
	assertUInt16sEqual(t, 776, consent.CmpVersion())
}

func TestLargeConsentScreen(t *testing.T) {
	consent, err := Parse(decode(t, "COyiFYuOyiFYuDKAA4ENAfCAAAAAAAAAAAAAAAAAAAAA"))
	assertNilError(t, err)
	assertUInt8sEqual(t, 56, consent.ConsentScreen())
}

func TestLanguageExtremes(t *testing.T) {
	consent, err := Parse(decode(t, "COyiHgFOyiHgFN4ABABGAPCAAAAAAAAAAAAAAFAAAAoAAAA"))
	assertNilError(t, err)
	assertStringsEqual(t, "BG", consent.ConsentLanguage())

	consent, err = Parse(decode(t, "COyiHgFOyiHgFN4ABASVAPCAAAAAAAAAAAAAAFAAAAoAAAA"))
	assertNilError(t, err)
	assertStringsEqual(t, "SV", consent.ConsentLanguage())
}

func TestTCF2Fields(t *testing.T) {
	baseConsent, err := Parse(decode(t, "COx3XOeOx3XOeLkAAAENAfCIAAAAAHgAAIAAAAAAAAAA"))
	assertNilError(t, err)
	consent := baseConsent.(ConsentMetadata)

	assertBoolsEqual(t, true, consent.PurposeOneTreatment())
	assertBoolsEqual(t, true, consent.SpecialFeatureOptIn(1))
	assertBoolsEqual(t, false, consent.SpecialFeatureOptIn(2))
}
