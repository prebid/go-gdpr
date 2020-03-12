package vendorlist

import (
	"testing"
)

func AssertVendorlistCorrectness(t *testing.T, parser func(data []byte) VendorList, version int) {
	if version == 1 {
		t.Run("TestVendorList", vendorListTester(parser))
		t.Run("TestVendor", vendorTester(parser))
	} else {
		t.Run("TestVendorList", vendorListTester20(parser))
		t.Run("TestVendor", vendorTester20(parser))
	}
}

func vendorListTester(parser func(data []byte) VendorList) func(*testing.T) {
	return func(t *testing.T) {
		list := parser([]byte(testData))
		assertIntsEqual(t, 5, int(list.Version()))
		assertNil(t, list.Vendor(2), true)
		assertNil(t, list.Vendor(32), false)
	}
}

func vendorListTester20(parser func(data []byte) VendorList) func(*testing.T) {
	return func(t *testing.T) {
		list := parser([]byte(testData20))
		assertIntsEqual(t, 28, int(list.Version()))
		assertNil(t, list.Vendor(2), true)
		assertNil(t, list.Vendor(8), false)
	}
}

func vendorTester(parser func(data []byte) VendorList) func(*testing.T) {
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

func vendorTester20(parser func(data []byte) VendorList) func(*testing.T) {
	return func(t *testing.T) {
		list := parser([]byte(testData20))
		v := list.Vendor(8)
		assertBoolsEqual(t, true, v.Purpose(1))
		assertBoolsEqual(t, true, v.Purpose(2))
		assertBoolsEqual(t, true, v.Purpose(3))
		assertBoolsEqual(t, true, v.Purpose(4))
		assertBoolsEqual(t, false, v.Purpose(5))
		assertBoolsEqual(t, false, v.Purpose(6))

		assertBoolsEqual(t, false, v.LegitimateInterest(1))
		assertBoolsEqual(t, true, v.LegitimateInterest(2))
		assertBoolsEqual(t, false, v.LegitimateInterest(3))

		v = list.Vendor(80)
		assertBoolsEqual(t, true, v.Purpose(1))
		assertBoolsEqual(t, true, v.Purpose(2))
		assertBoolsEqual(t, false, v.Purpose(3))
		assertBoolsEqual(t, true, v.Purpose(4))
		assertBoolsEqual(t, false, v.Purpose(5))
		assertBoolsEqual(t, false, v.Purpose(6))

		assertBoolsEqual(t, false, v.LegitimateInterest(1))
		assertBoolsEqual(t, true, v.LegitimateInterest(2))
		assertBoolsEqual(t, false, v.LegitimateInterest(3))
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

const testData20 = `
{
	"gvlSpecificationVersion": 2,
	"vendorListVersion": 28,
	"tcfPolicyVersion": 2,
	"lastUpdated": "2020-03-05T16:05:29Z",
	"vendors": {
		"8": {
			"id": 8,
			"name": "Emerse Sverige AB",
			"purposes": [1, 3, 4],
			"legIntPurposes": [2, 7, 8, 9],
			"flexiblePurposes": [2, 9],
			"specialPurposes": [1, 2],
			"features": [1, 2],
			"specialFeatures": [],
			"policyUrl": "https://www.emerse.com/privacy-policy/"
		},
		"80": {
			"id": 80,
			"name": "Sharethrough, Inc",
			"purposes": [1, 2, 4, 7, 9, 10],
			"legIntPurposes": [],
			"flexiblePurposes": [2, 4, 7, 9, 10],
			"specialPurposes": [],
			"features": [],
			"specialFeatures": [],
			"policyUrl": "https://platform-cdn.sharethrough.com/privacy-policy"
		}
	}
}
`

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
