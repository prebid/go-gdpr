package vendorlist2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEagerlyVendorList(t *testing.T) {
	tests := []struct {
		name                  string
		vendorList            string
		vendorListSpecVersion uint16
		vendorListVersion     uint16
	}{
		{
			name:                  "vendor_list_spec_2",
			vendorList:            testDataSpecVersion2,
			vendorListSpecVersion: 2,
			vendorListVersion:     28,
		},
		{
			name:                  "vendor_list_spec_3",
			vendorList:            testDataSpecVersion3,
			vendorListSpecVersion: 3,
			vendorListVersion:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedGVL, err := ParseEagerly([]byte(tt.vendorList))
			assert.NoError(t, err)
			assert.Equal(t, tt.vendorListSpecVersion, parsedGVL.SpecVersion())
			assert.Equal(t, tt.vendorListVersion, parsedGVL.Version())
			assert.NotNil(t, parsedGVL.Vendor(8))
			assert.NotNil(t, parsedGVL.Vendor(80))
			AssertVendorListCorrectness(t, parsedGVL)
		})
	}
}

func TestParseEagerlyEmptyVendorList(t *testing.T) {
	tests := []struct {
		name                  string
		vendorList            string
		vendorListSpecVersion uint16
		vendorListVersion     uint16
	}{
		{
			name:                  "vendor_list_spec_2",
			vendorList:            testDataSpecVersion2Empty,
			vendorListSpecVersion: 2,
			vendorListVersion:     28,
		},
		{
			name:                  "vendor_list_spec_3",
			vendorList:            testDataSpecVersion3Empty,
			vendorListSpecVersion: 3,
			vendorListVersion:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedGVL, err := ParseEagerly([]byte(tt.vendorList))
			assert.NoError(t, err)
			assert.Equal(t, tt.vendorListSpecVersion, parsedGVL.SpecVersion())
			assert.Equal(t, tt.vendorListVersion, parsedGVL.Version())
			assert.Nil(t, parsedGVL.Vendor(8))
			assert.Nil(t, parsedGVL.Vendor(80))
		})
	}
}
