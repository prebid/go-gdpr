package vendorconsent

import (
	"encoding/base64"
	"testing"

	tcf2 "github.com/prebid/go-gdpr/vendorconsent/tcf2"
)

// This checks error conditions to verify that we get errors back on malformed strings
func TestInvalidConsentStrings20(t *testing.T) {
	// All strings here were encoded using https://cryptii.com/binary-to-base64 from binary to URL-encoded base64 string.
	// Beware: this tool only makes sense if your binary strings use full bytes (multiples of 8 digits).
	//
	// For future tests, a "baseline" of valid binary using a BitField, segmented by different vendor consent string semantics, is:
	//
	// 000010                                 => Version
	// 001110001101011100100010100000101110   => Created date
	// 001110001101011100100011000110001010   => LastUpdated date
	// 000000000011                           => CmpId
	// 000000000010                           => CmpVersion
	// 000111                                 => ConsentScreen
	// 000100001101                           => ConsentLanguage
	// 000000001110                           => VendorListVersion
	// 000010								  => TcfPolicyVersion
	// 0 									  => IsServiceSpecific
	// 0									  => UseNonStandardStacks
	// 100000000000							  => SpecialFeatureOptIns
	// 001011010010110101101011               => PurposesConsent
	// 111111111100000000000000				  => PurposesLITransparency
	// 0									  => PurposeOneTreatment
	// 010100010010							  => PublisherCC (US if I did tge math right)
	// 0000000000000011                       => MaxVendorID <= Vendor Consent
	// 0                                      => EncodingType
	// 000                                    => BitFieldSection
	// 0000000000000011                       => MaxVendorID <= Legitimate Interest
	// 0                                      => EncodingType
	// 000                                    => BitFieldSection
	//
	// 0000100011100011010111001000101000001011100011100011010111001000110001100010100000000000110000000000100001110001000011010000000011100000100010000000000000101101001011010110101111111111110000000000000000101000100100000000000000011000000000000000000110000
	// CONciguONcjGKADACHENAOCIAC0ta__AACiQABgAAYA
	//
	// These "bad requests" can be made by tweaking those values to get various errors.
	// Bad metadata
	assertInvalid20(t, "CONciguONcjGKADACHENAOCIAC0ta__AACiQAA", "vendor consent strings are at least 29 bytes long. This one was 28")
	assertInvalid20(t, "BONciguONcjGKADACHENAOCIAC0ta__AACiQABgAAYA", "the consent string encoded a Version of 1, but this value must be greater than or equal to 2")
	assertInvalid20(t, "CONciguONcjGKADACHENAACIAC0ta__AACiQABgAAYA", "the consent string encoded a VendorListVersion of 0, but this value must be greater than or equal to 1")

	// Bad BitFields
	assertInvalid20(t, "CONciguONcjGKADACHENAOCIAC0ta__AACiQAeAA", "a BitField for 60 vendors requires a consent string of 37 bytes. This consent string had 30")

	// Bad RangeSections
	assertInvalid20(t, "CONciguONcjGKADACHENAOCIAC0ta__AACiQABwA", "vendor consent strings using RangeSections require at least 31 bytes to read NumEntries. Got 30")          // This encodes 184 bits
	assertInvalid20(t, "CONciguONcjGKADACHENAOCIAC0ta__AACiQABwAQQ", "ParseUInt16 expected a 16-bit int to start at bit 243, but the consent string was only 31 bytes long")   // 1 single vendor, too few bits
	assertInvalid20(t, "CONciguONcjGKADACHENAOCIAC0ta__AACiQABwAYQAC", "ParseUInt16 expected a 16-bit int to start at bit 259, but the consent string was only 33 bytes long") // 1 vendor range, too few bits
	assertInvalid20(t, "CONciguONcjGKADACHENAOCIAC0ta__AACiQABwAgABA", "ParseUInt16 expected a 16-bit int to start at bit 260, but the consent string was only 33 bytes long") // 2 single vendors, too few bits
	assertInvalid20(t, "CONciguONcjGKADACHENAOCIAC0ta__AACiQABwAgAAAAA", "bit 242 range entry excludes vendor 0, but only vendors [1, 3] are valid")
	assertInvalid20(t, "CONciguONcjGKADACHENAOCIAC0ta__AACiQABwAgACAAA", "bit 242 range entry excludes vendor 4, but only vendors [1, 3] are valid")
	assertInvalid20(t, "CONciguONcjGKADACHENAOCIAC0ta__AACiQABwAgABAAAA", "bit 259 range entry excludes vendor 0, but only vendors [1, 3] are valid")
	assertInvalid20(t, "CONciguONcjGKADACHENAOCIAC0ta__AACiQABwAgABAAEA", "bit 259 range entry excludes vendor 4, but only vendors [1, 3] are valid")
	assertInvalid20(t, "CONciguONcjGKADACHENAOCIAC0ta__AACiQABwAoABAACA", "bit 242 range entry excludes vendors [2, 1]. The start should be less than the end")
}

func TestParseValidString20(t *testing.T) {
	parsed, err := ParseString("COyiILmOyiILmADACHENAPCAAAAAAAAAAAAAE5QBgALgAqgD8AQACSwEygJyAAAAAA")
	assertNilError(t, err)
	assertUInt16sEqual(t, 15, parsed.VendorListVersion())
}

func TestParseValidString20MaxVendorID0(t *testing.T) {
	parsed, err := ParseString("COwJz-rOwJz-rMLAEAFRAPCgAAAAAAAAAAqIAAAAAAAA")
	assertNilError(t, err)
	// if vendor consent is empty, max vendor ID may be 0. See https://github.com/InteractiveAdvertisingBureau/iabtcf-es/issues/121
	assertUInt16sEqual(t, 0, parsed.MaxVendorID())
	assertUInt16sEqual(t, 15, parsed.VendorListVersion())
}

func TestParseEmptyString(t *testing.T) {
	// checks there is no panic due to array out of bound
	var err error

	_, err = ParseString("")
	if err == nil {
		t.Errorf("expected an error from ParseString")
	}

	_, err = ParseVersion(make([]byte, 0))
	if err == nil {
		t.Errorf("expected an error from ParseVersion")
	}
}

func assertInvalid20(t *testing.T, urlEncodedString string, expectError string) {
	t.Helper()
	data, err := base64.RawURLEncoding.DecodeString(urlEncodedString)
	assertNilError(t, err)
	assertInvalidBytes20(t, data, expectError)
}

func assertInvalidBytes20(t *testing.T, data []byte, expectError string) {
	t.Helper()
	if consent, err := tcf2.Parse(data); err == nil {
		t.Errorf("base64 URL-encoded string %s was considered valid, but shouldn't be. MaxVendorID: %d. len(data): %d", base64.RawURLEncoding.EncodeToString(data), consent.MaxVendorID(), len(data))
	} else if err.Error() != expectError {
		t.Errorf(`error messages did not match. Expected "%s", got "%s": %v`, expectError, err.Error(), err)
	}
}
