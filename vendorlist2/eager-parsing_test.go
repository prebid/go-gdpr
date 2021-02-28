package vendorlist2

import (
	"testing"

	"github.com/prebid/go-gdpr/api"
	"github.com/stretchr/testify/assert"
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

func TestParseEagerlyVendorsEmpty(t *testing.T) {
	vendorListJSON := `
{
	"gvlSpecificationVersion": 2,
	"vendorListVersion": 28,
	"tcfPolicyVersion": 2,
	"lastUpdated": "2020-03-05T16:05:29Z",
	"vendors": { }
}
`
	vendorList, err := ParseEagerly([]byte(vendorListJSON))

	assert.NoError(t, err)
	assert.Nil(t, vendorList.Vendor(0))
}
