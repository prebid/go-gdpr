package vendorlist

import (
	"testing"
)

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

func TestLazyParsedVendorlist(t *testing.T) {
	list := ParseLazily([]byte(testData))
	assertIntsEqual(t, 5, int(list.Version()))
	assertNil(t, list.Vendor(2), true)
	assertNil(t, list.Vendor(32), false)
}

func TestLazyParsedVendor(t *testing.T) {
	list := ParseLazily([]byte(testData))
	v := list.Vendor(32)
	assertBoolsEqual(t, true, v.Purpose(1))
	assertBoolsEqual(t, true, v.Purpose(2))
	assertBoolsEqual(t, false, v.Purpose(3))

	assertBoolsEqual(t, false, v.LegitimateInterest(1))
	assertBoolsEqual(t, false, v.LegitimateInterest(2))
	assertBoolsEqual(t, true, v.LegitimateInterest(3))
}

func assertIntsEqual(t *testing.T, expected int, actual int) {
	t.Helper()
	if actual != expected {
		t.Errorf("Ints were not equal. Expected %d, actual %d", expected, actual)
	}
}

func assertBoolsEqual(t *testing.T, expected bool, actual bool) {
	t.Helper()
	if actual != expected {
		t.Errorf("Bools were not equal. Expected %t, actual %t", expected, actual)
	}
}

func assertNil(t *testing.T, value Vendor, expectNil bool) {
	t.Helper()
	if expectNil && value != nil {
		t.Error("The vendor should be nil, but wasn't.")
	}
	if !expectNil && value == nil {
		t.Errorf("The vendor should not be nil, but was.")
	}
}
