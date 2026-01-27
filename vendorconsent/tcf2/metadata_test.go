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
	consent := baseConsent.(*ConsentMetadata)

	assertBoolsEqual(t, true, consent.PurposeOneTreatment())
	assertBoolsEqual(t, true, consent.SpecialFeatureOptIn(1))
	assertBoolsEqual(t, false, consent.SpecialFeatureOptIn(2))
}

func TestLITransparency(t *testing.T) {
	baseConsent, err := Parse(decode(t, "COx3XOeOx3XOeLkAAAENAfCIAAAAAHgAAIAAAAAAAAAA"))
	assertNilError(t, err)
	consent := baseConsent.(*ConsentMetadata)

	assertBoolsEqual(t, false, consent.PurposeLITransparency(1))
	assertBoolsEqual(t, true, consent.PurposeLITransparency(2))
	assertBoolsEqual(t, true, consent.PurposeLITransparency(3))
	assertBoolsEqual(t, true, consent.PurposeLITransparency(4))
	assertBoolsEqual(t, true, consent.PurposeLITransparency(5))
	assertBoolsEqual(t, false, consent.PurposeLITransparency(6))
	assertBoolsEqual(t, false, consent.PurposeLITransparency(7))
	assertBoolsEqual(t, false, consent.PurposeLITransparency(28))
}

func TestVendorDisclosed(t *testing.T) {
	baseConsent, err := ParseString("CQeakJVQeakJVHMAAAENCZCAAAAAAAAAAAAAAAAAAAAA.II7Nd_X__bX9n-_7_6ft0eY1f9_r37uQzDhfNs-8F3L_W_LwX32E7NF36tq4KmR4ku1bBIQNtHMnUDUmxaolVrzHsak2cpyNKJ_JkknsZe2dYGF9Pn9lD-YKZ7_5_9_f52T_9_9_-39z3_9f___dv_-__-vjf_599n_v9fV_78_Kf9______-____________8A")
	assertNilError(t, err)
	consent := baseConsent.(*ConsentMetadata)

	assertBoolsEqual(t, true, consent.VendorDisclosed(15))
}

func TestVendorDisclosed_NoSegment(t *testing.T) {
	// Consent string without disclosed vendors segment (no dot separator)
	baseConsent, err := ParseString("COvcSpYOvcSpYC9AAAENAPCAAAAAAAAAAAAACvwDQABAAIAAYABIAC4AJQAagA9ACEAPgAjIBJoCvAK-AAAAAA")
	assertNilError(t, err)
	consent := baseConsent.(*ConsentMetadata)

	// All vendors should return false when no disclosed segment exists
	assertBoolsEqual(t, false, consent.VendorDisclosed(1))
	assertBoolsEqual(t, false, consent.VendorDisclosed(15))
	assertBoolsEqual(t, false, consent.VendorDisclosed(100))
}

func TestVendorDisclosed_NotInList(t *testing.T) {
	// Consent string with disclosed vendors segment (bitfield encoding, MaxVendorID=1142)
	// Not all vendors 1-1142 are disclosed - only specific ones have their bits set
	baseConsent, err := ParseString("CQeakJVQeakJVHMAAAENCZCAAAAAAAAAAAAAAAAAAAAA.II7Nd_X__bX9n-_7_6ft0eY1f9_r37uQzDhfNs-8F3L_W_LwX32E7NF36tq4KmR4ku1bBIQNtHMnUDUmxaolVrzHsak2cpyNKJ_JkknsZe2dYGF9Pn9lD-YKZ7_5_9_f52T_9_9_-39z3_9f___dv_-__-vjf_599n_v9fV_78_Kf9______-____________8A")
	assertNilError(t, err)
	consent := baseConsent.(*ConsentMetadata)

	// Vendors that ARE disclosed
	assertBoolsEqual(t, true, consent.VendorDisclosed(1))
	assertBoolsEqual(t, true, consent.VendorDisclosed(15))
	assertBoolsEqual(t, true, consent.VendorDisclosed(1142)) // MaxVendorID boundary

	// Vendors that are NOT disclosed (bits are 0 in bitfield)
	assertBoolsEqual(t, false, consent.VendorDisclosed(51))
	assertBoolsEqual(t, false, consent.VendorDisclosed(999))
	assertBoolsEqual(t, false, consent.VendorDisclosed(1143))
}

func TestVendorDisclosed_RangeEncoding(t *testing.T) {
	// Range-encoded disclosed vendors segment with:
	// - MaxVendorId = 100
	// - 2 entries:
	//   - Single vendor 10 (IsARange=0)
	//   - Range 50-75 (IsARange=1)
	// Segment "IAyQAgAFQAyAEsAA" encodes this structure
	baseConsent, err := ParseString("COvcSpYOvcSpYC9AAAENAPCAAAAAAAAAAAAACvwDQABAAIAAYABIAC4AJQAagA9ACEAPgAjIBJoCvAK-AAAAAA.IAyQAgAFQAyAEsAA")
	assertNilError(t, err)
	consent := baseConsent.(*ConsentMetadata)

	// Single vendor entry (vendor 10)
	assertBoolsEqual(t, true, consent.VendorDisclosed(10))

	// Range entry (vendors 50-75)
	assertBoolsEqual(t, true, consent.VendorDisclosed(50)) // range start
	assertBoolsEqual(t, true, consent.VendorDisclosed(60)) // in range
	assertBoolsEqual(t, true, consent.VendorDisclosed(75)) // range end

	// Vendors NOT in disclosed list
	assertBoolsEqual(t, false, consent.VendorDisclosed(49))  // just before range
	assertBoolsEqual(t, false, consent.VendorDisclosed(76))  // just after range
	assertBoolsEqual(t, false, consent.VendorDisclosed(100)) // at MaxVendorId but not disclosed
}

func TestVendorDisclosed_MalformedSegment(t *testing.T) {
	// Consent string with invalid base64 in disclosed segment
	// The parser should return an error and not panic
	_, err := ParseString("COvcSpYOvcSpYC9AAAENAPCAAAAAAAAAAAAACvwDQABAAIAAYABIAC4AJQAagA9ACEAPgAjIBJoCvAK-AAAAAA.!!INVALID!!")
	assertError(t, err)
}

func TestTCF_2_3_RequiresDisclosedVendors(t *testing.T) {
	// TCF 2.2 (policy version 4) without disclosed vendors - should succeed
	// Policy version is stored at bits 133-138
	tcf22NoDisclosed := "CPtGDMAPtGDMALMAAAENA_EIAAAAAAAAACiQAAAAAAAA" // policy version 4
	_, err := ParseString(tcf22NoDisclosed)
	assertNilError(t, err)

	// TCF 2.3 (policy version 5) without disclosed vendors - should fail
	tcf23NoDisclosed := "CPtGDMAPtGDMALMAAAENA_FIAAAAAAAAACiQAAAAAAAA" // policy version 5
	_, err = ParseString(tcf23NoDisclosed)
	assertError(t, err)

	// TCF 2.3 (policy version 5) with disclosed vendors - should succeed
	tcf23WithDisclosed := "CPtGDMAPtGDMALMAAAENA_FIAAAAAAAAACiQAAAAAAAA.IAAAAAAAAAA"
	_, err = ParseString(tcf23WithDisclosed)
	assertNilError(t, err)
}
