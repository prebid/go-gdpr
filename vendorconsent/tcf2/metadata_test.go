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

func TestTCFPolicyVersion(t *testing.T) {
	baseConsent := "CPtGDMAPtGDMALMAAAENA_C_AAAAAAAAACiQAAAAAAAA"
	index := 22 // policy version is at the 23rd 6-bit base64 position
	tests := []struct {
		name       string
		base64Char string
		expected   uint8
	}{
		{
			name:       "char_A_bits_000000_is_version_0",
			base64Char: "A",
			expected:   0,
		},
		{
			name:       "char_B_bits_000001_is_version_1",
			base64Char: "B",
			expected:   1,
		},
		{
			name:       "char_C_bits_000010_is_version_2",
			base64Char: "C",
			expected:   2,
		},
		{
			name:       "char_E_bits_000100_is_version_4",
			base64Char: "E",
			expected:   4,
		},
		{
			name:       "char_I_bits_001000_is_version_8",
			base64Char: "I",
			expected:   8,
		},
		{
			name:       "char_Q_bits_010000_is_version_16",
			base64Char: "Q",
			expected:   16,
		},
		{
			name:       "char_g_bits_100000_is_version_32",
			base64Char: "g",
			expected:   32,
		},
		{
			name:       "char_underscore_bits_111111_is_version_63",
			base64Char: "_",
			expected:   63,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedConsent := baseConsent[:index] + tt.base64Char + baseConsent[index+1:]
			consent, err := Parse(decode(t, updatedConsent))
			assertNilError(t, err)
			assertUInt8sEqual(t, tt.expected, consent.TCFPolicyVersion())
		})
	}
}

func TestTCF2Fields(t *testing.T) {
	baseConsent, err := Parse(decode(t, "COx3XOeOx3XOeLkAAAENAfCIAAAAAHgAAIAAAAAAAAAA"))
	assertNilError(t, err)
	consent := baseConsent.(ConsentMetadata)

	assertBoolsEqual(t, true, consent.PurposeOneTreatment())
	assertBoolsEqual(t, true, consent.SpecialFeatureOptIn(1))
	assertBoolsEqual(t, false, consent.SpecialFeatureOptIn(2))
}

func TestLITransparency(t *testing.T) {
	baseConsent, err := Parse(decode(t, "COx3XOeOx3XOeLkAAAENAfCIAAAAAHgAAIAAAAAAAAAA"))
	assertNilError(t, err)
	consent := baseConsent.(ConsentMetadata)

	assertBoolsEqual(t, false, consent.PurposeLITransparency(1))
	assertBoolsEqual(t, true, consent.PurposeLITransparency(2))
	assertBoolsEqual(t, true, consent.PurposeLITransparency(3))
	assertBoolsEqual(t, true, consent.PurposeLITransparency(4))
	assertBoolsEqual(t, true, consent.PurposeLITransparency(5))
	assertBoolsEqual(t, false, consent.PurposeLITransparency(6))
	assertBoolsEqual(t, false, consent.PurposeLITransparency(7))
	assertBoolsEqual(t, false, consent.PurposeLITransparency(28))

}

func TestVendorDisclosedWithoutSegment(t *testing.T) {
	// TCF 2.0/2.2 consent string without disclosed vendors segment
	baseConsent, err := Parse(decode(t, "COx3XOeOx3XOeLkAAAENAfCIAAAAAHgAAIAAAAAAAAAA"))
	assertNilError(t, err)
	consent := baseConsent.(ConsentMetadata)

	// Should return false for all vendors when no disclosed vendors segment exists
	assertBoolsEqual(t, false, consent.VendorDisclosed(1))
	assertBoolsEqual(t, false, consent.VendorDisclosed(10))
	assertBoolsEqual(t, false, consent.VendorDisclosed(100))
	assertBoolsEqual(t, false, consent.VendorDisclosed(999))
}

func TestVendorDisclosedNilCheck(t *testing.T) {
	// Parse core string directly (no disclosed vendors)
	baseConsent, err := Parse(decode(t, "COx3XOeOx3XOeLkAAAENAfCIAAAAAHgAAIAAAAAAAAAA"))
	assertNilError(t, err)
	consent := baseConsent.(ConsentMetadata)

	// Verify disclosedVendors field is nil
	if consent.disclosedVendors != nil {
		t.Errorf("Expected disclosedVendors to be nil for core string without segment")
	}

	// VendorDisclosed should safely return false when disclosedVendors is nil
	assertBoolsEqual(t, false, consent.VendorDisclosed(1))
}

func TestHasDisclosedVendorsWithoutSegment(t *testing.T) {
	// TCF 2.0/2.2 consent string without disclosed vendors segment
	baseConsent, err := Parse(decode(t, "COx3XOeOx3XOeLkAAAENAfCIAAAAAHgAAIAAAAAAAAAA"))
	assertNilError(t, err)
	consent := baseConsent.(ConsentMetadata)

	// HasDisclosedVendors should return false when segment is not present
	assertBoolsEqual(t, false, consent.HasDisclosedVendors())
}
