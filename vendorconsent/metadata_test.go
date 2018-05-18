package vendorconsent

import "testing"

func TestLargeVendorListVersion(t *testing.T) {
	consent, err := Parse(decode(t, "BON96hFON96hFABABBAA4yAAAAAAEA"))
	if err != nil {
		t.Fatalf("Failed to parse valid consent string: %v", err)
	}

	assertUInt16sEqual(t, 3634, consent.VendorListVersion())
}
