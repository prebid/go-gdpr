package vendorconsent

import (
	"encoding/base64"
	"testing"
)

// TestParseDisclosedVendors tests parsing of TCF 2.3 strings with disclosed vendors segment
func TestParseDisclosedVendors(t *testing.T) {
	// This is a TCF 2.3 string with disclosed vendors segment
	// Format: CoreString.DisclosedVendors

	// Core string (existing valid TCF 2.0 string)
	coreString := "COyiILmOyiILmADACHENAPCAAAAAAAAAAAAAE5QBgALgAqgD8AQACSwEygJyAAAAAA"

	// Create a disclosed vendors segment manually for testing
	// Binary structure:
	// - SegmentType=1: 001 (3 bits)
	// - MaxVendorId=10: 0000000000001010 (16 bits)
	// - IsRangeEncoding=0: 0 (1 bit - bitfield mode)
	// - Vendor bits (10 bits): 1010100000 (vendors 1, 3, 5 disclosed)
	//
	// Bit string: 001|0000000000001010|0|1010100000
	// Bytes: 00100000|00000001|01001010|10000000
	//        0x20     0x01     0x4a     0x80
	disclosedVendorsBytes := []byte{0x20, 0x01, 0x4a, 0x80}
	disclosedVendorsString := base64.RawURLEncoding.EncodeToString(disclosedVendorsBytes)

	consentString := coreString + "." + disclosedVendorsString

	consent, err := ParseString(consentString)
	assertNilError(t, err)

	// Test that core parsing still works
	assertUInt16sEqual(t, 15, consent.VendorListVersion())

	// Test disclosed vendors
	assertBoolsEqual(t, true, consent.VendorDisclosed(1))
	assertBoolsEqual(t, false, consent.VendorDisclosed(2))
	assertBoolsEqual(t, true, consent.VendorDisclosed(3))
	assertBoolsEqual(t, false, consent.VendorDisclosed(4))
	assertBoolsEqual(t, true, consent.VendorDisclosed(5))
	assertBoolsEqual(t, false, consent.VendorDisclosed(6))

	// Test HasDisclosedVendors
	assertBoolsEqual(t, true, consent.HasDisclosedVendors())
}

// TestBackwardCompatibilityNoDisclosedVendors tests that TCF 2.0/2.2 strings without
// disclosed vendors segment still work (backward compatibility)
func TestBackwardCompatibilityNoDisclosedVendors(t *testing.T) {
	// TCF 2.0 string without disclosed vendors segment
	consentString := "COyiILmOyiILmADACHENAPCAAAAAAAAAAAAAE5QBgALgAqgD8AQACSwEygJyAAAAAA"

	consent, err := ParseString(consentString)
	assertNilError(t, err)

	// Test that core parsing works
	assertUInt16sEqual(t, 15, consent.VendorListVersion())

	// VendorDisclosed should return false when no disclosed vendors segment exists
	assertBoolsEqual(t, false, consent.VendorDisclosed(1))
	assertBoolsEqual(t, false, consent.VendorDisclosed(100))

	// HasDisclosedVendors should return false
	assertBoolsEqual(t, false, consent.HasDisclosedVendors())
}

// TestEmptyDisclosedVendorsSegment tests handling of empty disclosed vendors segment
func TestEmptyDisclosedVendorsSegment(t *testing.T) {
	// TCF string with empty disclosed vendors segment (just a dot)
	consentString := "COyiILmOyiILmADACHENAPCAAAAAAAAAAAAAE5QBgALgAqgD8AQACSwEygJyAAAAAA."

	consent, err := ParseString(consentString)
	assertNilError(t, err)

	// Should still parse core string successfully
	assertUInt16sEqual(t, 15, consent.VendorListVersion())

	// VendorDisclosed should return false (no vendors disclosed)
	assertBoolsEqual(t, false, consent.VendorDisclosed(1))

	// HasDisclosedVendors should return false for empty segment
	assertBoolsEqual(t, false, consent.HasDisclosedVendors())
}

// TestMultipleSegments tests parsing string with multiple segments (core + disclosed + publisher)
func TestMultipleSegments(t *testing.T) {
	// Core string + disclosed vendors + publisher TC (third segment)
	coreString := "COyiILmOyiILmADACHENAPCAAAAAAAAAAAAAE5QBgALgAqgD8AQACSwEygJyAAAAAA"
	disclosedVendorsBytes := []byte{0x20, 0x01, 0x4a, 0x80}
	disclosedVendorsString := base64.RawURLEncoding.EncodeToString(disclosedVendorsBytes)
	publisherTCString := "YAAAAAAAAAAA" // placeholder

	consentString := coreString + "." + disclosedVendorsString + "." + publisherTCString

	consent, err := ParseString(consentString)
	assertNilError(t, err)

	// Core parsing should work
	assertUInt16sEqual(t, 15, consent.VendorListVersion())

	// HasDisclosedVendors should return true
	assertBoolsEqual(t, true, consent.HasDisclosedVendors())

	// Disclosed vendors should be parsed (ignoring publisher TC for now)
	assertBoolsEqual(t, true, consent.VendorDisclosed(1))
}

// TestSegmentsInAnyOrder tests that segments can appear in any order (TCF spec allows this)
func TestSegmentsInAnyOrder(t *testing.T) {
	coreString := "COwGVJOOwGVJOADACHENAOCAAO6as_-AAAhoAFNLAAoAAAA"

	// Disclosed vendors segment (type=1)
	// 001 0000000000011010 0 10101000000000100000100000000101000010011111
	// |    |               | |_ bitset
	// |    |               |__ IsRangeEncoding (0)
	// |    |_ maxVendorID (26)
	// |_ segment type (1)
	disclosedVendorsBytes := []byte{0x20, 0x03, 0x4a, 0x80, 0x20, 0x80, 0x50, 0x9f}

	disclosedVendorsString := base64.RawURLEncoding.EncodeToString(disclosedVendorsBytes)

	// Publisher TC segment (type=3) - minimal valid segment
	// Binary: 011|0000000000000000|... (type=3, no publisher restrictions)
	publisherTCBytes := []byte{0x60, 0x00, 0x00}
	publisherTCString := base64.RawURLEncoding.EncodeToString(publisherTCBytes)

	// Test order 1: Core.Disclosed.Publisher
	consent1, err := ParseString(coreString + "." + disclosedVendorsString + "." + publisherTCString)

	assertNilError(t, err)
	assertBoolsEqual(t, true, consent1.HasDisclosedVendors())

	assertUInt16sEqual(t, 26, consent1.VendorDisclosedMaxVendorID())

	assertBoolsEqual(t, true, consent1.VendorDisclosed(1))
	assertBoolsEqual(t, true, consent1.VendorDisclosed(3))
	assertBoolsEqual(t, true, consent1.VendorDisclosed(21))
	assertBoolsEqual(t, false, consent1.VendorDisclosed(27)) // greater than vendorDisclosedMaxVendorId

	// Test order 2: Core.Publisher.Disclosed (reversed order)
	consent2, err := ParseString(coreString + "." + publisherTCString + "." + disclosedVendorsString)

	// Both should have disclosed vendors
	assertBoolsEqual(t, true, consent2.HasDisclosedVendors())
	assertNilError(t, err)
	assertBoolsEqual(t, true, consent2.VendorDisclosed(1))
	assertBoolsEqual(t, true, consent2.VendorDisclosed(3))

	// Both should give same results
	assertBoolsEqual(t, consent1.VendorDisclosed(5), consent2.VendorDisclosed(5))
}
