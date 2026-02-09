package vendorconsent

import (
	"testing"
)

func TestParseLegitIntSetWithBitField(t *testing.T) {
	// this test uses a crafted consent uses bit field, declares 10 vendors and legitimate interest without required content
	_, err := Parse(decode(t, "COvcSpYOvcSpYC9AAAENAPCAAAAAAAAAAAAAAFAAAAA"))
	assertError(t, err)
}

func TestParseLegitIntSetWithRangeSection(t *testing.T) {
	// this test uses a crafted consent uses range section, declares 10 vendors, 6 exceptions and legitimate interest without required content
	_, err := Parse(decode(t, "COvcSpYOvcSpYC9AAAENAPCAAAAAAAAAAAAAAFQBgAAgABAACAAEAAQAAgAA"))
	assertError(t, err)
}

func TestParseStringWithDisclosedVendors(t *testing.T) {
	// TCF 2.3 string with Disclosed Vendors segment (range-encoded):
	// Disclosed vendors: 1-5, 100, 404 (MaxVendorId=404)
	consent, err := ParseString("CQSbk4AQSbk4ANwAAAENAwCgAAAAAAAAAAYgACPAAAAA.IDKQA4AAgAKAGQAygAAA.YAAAAAAAAAAA")
	assertNilError(t, err)

	metadata := consent.(ConsentMetadata)
	assertUInt16sEqual(t, 404, metadata.MaxDisclosedVendorID())

	// Vendors in range 1-5 should be disclosed
	for id := uint16(1); id <= 5; id++ {
		assertBoolsEqual(t, true, metadata.DisclosedVendor(id))
	}
	// Vendor 6 should not be disclosed
	assertBoolsEqual(t, false, metadata.DisclosedVendor(6))
	// Vendor 100 should be disclosed
	assertBoolsEqual(t, true, metadata.DisclosedVendor(100))
	// Vendor 101 should not be disclosed
	assertBoolsEqual(t, false, metadata.DisclosedVendor(101))
	// Vendor 404 should be disclosed
	assertBoolsEqual(t, true, metadata.DisclosedVendor(404))

	// Core string fields should still work
	assertUInt8sEqual(t, 2, metadata.Version())
}

func TestParseStringWithoutDisclosedVendors(t *testing.T) {
	// TCF 2.2 string without Disclosed Vendors segment
	consent, err := ParseString("CQc78IAQc78IAAHABAFRCPFsAP_gAAAAAAAAKlwJ4AFgAYABUAC4AGQAQAAnABaADIAGgAWwAwgBzAD8AIQATgAuABlADjAIQARAAicBHAEdAJKAYoA0ACIgETAKWAXUAvMBgIDFgGMgMsAf2BAECMwEdgKlgAAAGKQAYAAgvwOgAwABBfghABgACC_BKADAAEF-AkAGAAIL8FoAMAAQX4A")
	assertNilError(t, err)

	metadata := consent.(ConsentMetadata)
	assertUInt16sEqual(t, 0, metadata.MaxDisclosedVendorID())
	assertBoolsEqual(t, false, metadata.DisclosedVendor(1))
}

func TestParseStringDisclosedVendorEdgeCases(t *testing.T) {
	// TCF 2.3 string - test boundary vendor IDs
	consent, err := ParseString("CQSbk4AQSbk4ANwAAAENAwCgAAAAAAAAAAYgACPAAAAA.IDKQA4AAgAKAGQAygAAA.YAAAAAAAAAAA")
	assertNilError(t, err)

	metadata := consent.(ConsentMetadata)

	// Vendor ID 0 should return false
	assertBoolsEqual(t, false, metadata.DisclosedVendor(0))
	// Vendor ID beyond MaxDisclosedVendorID should return false
	assertBoolsEqual(t, false, metadata.DisclosedVendor(405))
}
