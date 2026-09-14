package vendorlist

import (
	"testing"

	"github.com/prebid/go-gdpr/api"
)

func TestEagerlyParsedVendorList(t *testing.T) {
	AssertVendorlistCorrectness(t, func(data []byte) api.VendorList {
		vendorList, err := ParseEagerly(data)
		if err != nil {
			t.Errorf("ParseEagerly returned an unexpected error: %v", err)
		}

		return vendorList
	})
}
