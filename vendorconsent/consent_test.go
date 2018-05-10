package vendorconsent

import (
	"encoding/base64"
	"testing"
)

func TestIsSet(t *testing.T) {
	data := []byte{0xf6, 0x37}
	assertBoolsEqual(t, true, isSet(data, 0))
	assertBoolsEqual(t, true, isSet(data, 1))
	assertBoolsEqual(t, true, isSet(data, 2))
	assertBoolsEqual(t, true, isSet(data, 3))
	assertBoolsEqual(t, false, isSet(data, 4))
	assertBoolsEqual(t, true, isSet(data, 5))
	assertBoolsEqual(t, true, isSet(data, 6))
	assertBoolsEqual(t, false, isSet(data, 7))

	assertBoolsEqual(t, false, isSet(data, 8))
	assertBoolsEqual(t, false, isSet(data, 9))
	assertBoolsEqual(t, true, isSet(data, 10))
	assertBoolsEqual(t, true, isSet(data, 11))
	assertBoolsEqual(t, false, isSet(data, 12))
	assertBoolsEqual(t, true, isSet(data, 13))
	assertBoolsEqual(t, true, isSet(data, 14))
	assertBoolsEqual(t, true, isSet(data, 15))
}

// This checks error conditions to verify that we get errors back on malformed strings
func TestInvalidConsentStrings(t *testing.T) {
	// All strings here were encoded using https://cryptii.com/binary-to-base64 from binary to URL-encoded base64 string.
	// Beware: this tool only makes sense if your binary strings use full bytes (multiples of 8 digits).
	//
	// For future tests, a "basline" of valid binary using a BitField, segmented by different vendor consent string semantics, is:
	//
	// 000001                                 => Version
	// 001110001101011100100010100000101110   => Created date
	// 001110001101011100100011000110001010   => LastUpdated date
	// 000000000011                           => CmpId
	// 000000000010                           => CmpVersion
	// 000111                                 => ConsentScreen
	// 000100001101                           => ConsentLangugae
	// 000000001110                           => VendorListVersion
	// 001011010010110101101011               => PurposesAllowed
	// 0000000000000011                       => MaxVendorID
	// 0                                      => EncodingType
	// 000                                    => BitFieldSection
	//
	// These "bad requests" can be made by tweaking those values to get various errors.
	// Bad metadata
	assertInvalid(t, "BONciguONcjGKADACHENAOLS1rAA", "vendor consent strings are at least 22 bytes long. This one was 21")
	assertInvalid(t, "AONciguONcjGKADACHENAOLS1rAAMA", "the consent string encoded a Version of 0, but this value must be greater than or equal to 1")
	assertInvalid(t, "BONciguONcjGKADACHENAOLS1rAAAA", "the consent string encoded a MaxVendorID of 0, but this value must be greater than or equal to 1")
	assertInvalid(t, "BONciguONcjGKADACHENAALS1rAAMA", "the consent string encoded a VendorListVersion of 0, but this value must be greater than or equal to 1")

	// Bad BitFields
	assertInvalid(t, "BONciguONcjGKADACHENAOLS1rAAQA", "a BitField for 4 vendors requires a consent string of 23 bytes. This consent string had 22")
	assertInvalid(t, "BONciguONcjGKADACHENAOLS1rAAwAA", "a BitField for 12 vendors requires a consent string of 24 bytes. This consent string had 23")

	// Bad RangeSections
	assertInvalid(t, "BONciguONcjGKADACHENAOLS1rAAOAA", "vendor consent strings using RangeSections require at least 24 bytes. Got 23")                                   // This encodes 184 bits
	assertInvalid(t, "BONciguONcjGKADACHENAOLS1rAAPABAAA", "rangeSection expected a 16-bit vendorID to start at bit 187, but the consent string was only 25 bytes long")  // 1 single vendor, too few bits
	assertInvalid(t, "BONciguONcjGKADACHENAOLS1rAAPABgACA", "rangeSection expected a 16-bit vendorID to start at bit 203, but the consent string was only 26 bytes long") // 1 vendor range, too few bits
	assertInvalid(t, "BONciguONcjGKADACHENAOLS1rAAPACAACA", "rangeSection expected a 16-bit vendorID to start at bit 204, but the consent string was only 26 bytes long") // 2 single vendors, too few bits
	assertInvalid(t, "BONciguONcjGKADACHENAOLS1rAAPABAAAA", "bit 186 range entry excludes vendor 0, but only vendors [1, 3] are valid")
	assertInvalid(t, "BONciguONcjGKADACHENAOLS1rAAPABAAIA", "bit 186 range entry excludes vendor 4, but only vendors [1, 3] are valid")
	assertInvalid(t, "BONciguONcjGKADACHENAOLS1rAAPABgAAAAIA", "bit 186 range entry exclusion starts at 0, but the min vendor ID is 1")
	assertInvalid(t, "BONciguONcjGKADACHENAOLS1rAAPABgACAAgA", "bit 186 range entry exclusion ends at 4, but the max vendor ID is 3")
	assertInvalid(t, "BONciguONcjGKADACHENAOLS1rAAPABgAEAAIA", "bit 186 range entry excludes vendors [2, 1]. The start should be less than the end")
}

func assertInvalid(t *testing.T, urlEncodedString string, expectError string) {
	t.Helper()
	data, err := base64.RawURLEncoding.DecodeString(urlEncodedString)
	if err != nil {
		t.Fatalf("Failed to base64-decode string %s. Error: %v", urlEncodedString, err)
	}
	if consent, err := Parse(data); err == nil {
		t.Errorf("base64 URL-encoded string %s was considered valid, but shouldn't be. MaxVendorID: %d. len(data): %d", urlEncodedString, consent.MaxVendorID(), len(data))
	} else if err.Error() != expectError {
		t.Errorf(`error messages did not match. Expected "%s", got "%s": %v`, expectError, err.Error(), err)
	}
}

func decode(t *testing.T, encodedString string) []byte {
	data, err := base64.RawURLEncoding.DecodeString(encodedString)
	if err != nil {
		t.Fatalf("Failed to base64-decode string %s. Error: %v", encodedString, err)
	}
	return data
}

func assertUInt8sEqual(t *testing.T, expected uint8, actual uint8) {
	t.Helper()
	if actual != expected {
		t.Errorf("Ints were not equal. Expected %d, actual %d", expected, actual)
	}
}

func assertUInt16sEqual(t *testing.T, expected uint16, actual uint16) {
	t.Helper()
	if actual != expected {
		t.Errorf("Ints were not equal. Expected %d, actual %d", expected, actual)
	}
}

func assertIntsEqual(t *testing.T, expected int, actual int) {
	t.Helper()
	if actual != expected {
		t.Errorf("Ints were not equal. Expected %d, actual %d", expected, actual)
	}
}

func assertBoolsEqual(t *testing.T, expected bool, actual bool) {
	t.Helper()
	if actual != expected {
		t.Errorf("Bools were not equal. Expected %t, actual %t", expected, actual)
	}
}
