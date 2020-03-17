package vendorlist2

import (
	"testing"
)

func TestLazyParsedVendorList(t *testing.T) {
	AssertVendorlistCorrectness(t, ParseLazily)
}
