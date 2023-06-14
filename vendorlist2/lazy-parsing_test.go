package vendorlist2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLazilyVendorList(t *testing.T) {
	tests := []struct{
		name              string
		vendorList        string
		vendorListVersion uint16
	}{
		{
			name:       "vendor list spec 2",
			vendorList: testDataSpecVersion2,
			vendorListVersion: 28,
		},
		{
			name:       "vendor list spec 3",
			vendorList: testDataSpecVersion3,
			vendorListVersion: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedGVL := ParseLazily([]byte(tt.vendorList))
			assert.Equal(t, tt.vendorListVersion, parsedGVL.Version())
			assert.NotNil(t, parsedGVL.Vendor(8))
			assert.NotNil(t, parsedGVL.Vendor(80))
			AssertVendorListCorrectness(t, parsedGVL)
		})
	}
}

func TestParseLazilyEmptyVendorList(t *testing.T) {
	tests := []struct{
		name              string
		vendorList        string
		vendorListVersion uint16
	}{
		{
			name:       "vendor list spec 2",
			vendorList: testDataSpecVersion2Empty,
			vendorListVersion: 28,
		},
		{
			name:       "vendor list spec 3",
			vendorList: testDataSpecVersion3Empty,
			vendorListVersion: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedGVL := ParseLazily([]byte(tt.vendorList))
			assert.Equal(t, tt.vendorListVersion, parsedGVL.Version())
			assert.Nil(t, parsedGVL.Vendor(8))
			assert.Nil(t, parsedGVL.Vendor(80))
		})
	}
}
