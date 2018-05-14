package vendorlist

import (
	"testing"
)

func TestLazyParsedVendorList(t *testing.T) {
	AssertVendorlistCorrectness(t, ParseLazily)
}
