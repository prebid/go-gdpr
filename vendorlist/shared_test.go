package vendorlist

import (
	"testing"

	"github.com/prebid/go-gdpr/api"
)

func AssertVendorlistCorrectness(t *testing.T, parser func(data []byte) api.VendorList) {
	t.Run("TestVendorList", vendorListTester(parser))
	t.Run("TestVendor", vendorTester(parser))
}

func vendorListTester(parser func(data []byte) api.VendorList) func(*testing.T) {
	return func(t *testing.T) {
		list := parser([]byte(testData))
		assertIntsEqual(t, 5, int(list.Version()))
		assertIntsEqual(t, 0, int(list.SpecVersion()))
		assertNil(t, list.Vendor(2), true)
		assertNil(t, list.Vendor(32), false)
	}
}

func vendorTester(parser func(data []byte) api.VendorList) func(*testing.T) {
	return func(t *testing.T) {
		list := parser([]byte(testData))
		v := list.Vendor(32)
		assertBoolsEqual(t, true, v.Purpose(1))
		assertBoolsEqual(t, true, v.Purpose(2))
		assertBoolsEqual(t, false, v.Purpose(3))

		assertBoolsEqual(t, false, v.LegitimateInterest(1))
		assertBoolsEqual(t, false, v.LegitimateInterest(2))
		assertBoolsEqual(t, true, v.LegitimateInterest(3))
	}
}

const testData = `
{
  "vendorListVersion": 5,
  "vendors": [
    {
      "id": 32,
      "purposeIds": [
        1,
        2
      ],
      "legIntPurposeIds": [
        3
      ],
      "featureIds": [
        2,
        3
      ]
    }
  ]
}
`

func assertIntsEqual(t *testing.T, expected, actual int) {
	t.Helper()

	if actual != expected {
		t.Errorf("Ints were not equal. Expected %d, actual %d", expected, actual)
	}
}

func assertBoolsEqual(t *testing.T, expected, actual bool) {
	t.Helper()

	if actual != expected {
		t.Errorf("Bools were not equal. Expected %t, actual %t", expected, actual)
	}
}

func assertNil(t *testing.T, value api.Vendor, expectNil bool) {
	t.Helper()

	if expectNil && value != nil {
		t.Error("The vendor should be nil, but wasn't.")
	}

	if !expectNil && value == nil {
		t.Errorf("The vendor should not be nil, but was.")
	}
}
