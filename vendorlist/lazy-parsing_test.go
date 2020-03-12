package vendorlist

import (
	"testing"
)

func TestLazyParsedVendorList(t *testing.T) {
	AssertVendorlistCorrectness(t, ParseLazily, 1)
	AssertVendorlistCorrectness(t, ParseLazily20, 2)
}
