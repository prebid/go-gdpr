package vendorlist2

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
		assertIntsEqual(t, 28, int(list.Version()))
		assertIntsEqual(t, 2, int(list.SpecVersion()))
		assertNil(t, list.Vendor(2), true)
		assertNil(t, list.Vendor(8), false)
	}
}

func vendorTester(parser func(data []byte) api.VendorList) func(*testing.T) {
	return func(t *testing.T) {
		list := parser([]byte(testData))
		v := list.Vendor(8)
		assertBoolsEqual(t, true, v.Purpose(1))
		assertBoolsEqual(t, true, v.PurposeStrict(1))
		assertBoolsEqual(t, true, v.Purpose(2))
		assertBoolsEqual(t, false, v.PurposeStrict(2))
		assertBoolsEqual(t, true, v.Purpose(3))
		assertBoolsEqual(t, true, v.PurposeStrict(3))
		assertBoolsEqual(t, true, v.Purpose(4))
		assertBoolsEqual(t, true, v.PurposeStrict(4))
		assertBoolsEqual(t, false, v.Purpose(5))
		assertBoolsEqual(t, false, v.PurposeStrict(5))
		assertBoolsEqual(t, false, v.Purpose(6))
		assertBoolsEqual(t, false, v.PurposeStrict(6))

		assertBoolsEqual(t, false, v.LegitimateInterest(1))
		assertBoolsEqual(t, false, v.LegitimateInterestStrict(1))
		assertBoolsEqual(t, true, v.LegitimateInterest(2))
		assertBoolsEqual(t, true, v.LegitimateInterestStrict(2))
		assertBoolsEqual(t, false, v.LegitimateInterest(3))
		assertBoolsEqual(t, false, v.LegitimateInterestStrict(3))

		assertBoolsEqual(t, true, v.SpecialPurpose(1))
		assertBoolsEqual(t, true, v.SpecialPurpose(2))
		assertBoolsEqual(t, false, v.SpecialPurpose(3)) // Does not exist yet

		assertBoolsEqual(t, true, v.SpecialFeature(1))
		assertBoolsEqual(t, true, v.SpecialFeature(2))
		assertBoolsEqual(t, false, v.SpecialFeature(3)) // Does not exist yet

		v = list.Vendor(80)
		assertBoolsEqual(t, true, v.Purpose(1))
		assertBoolsEqual(t, true, v.PurposeStrict(1))
		assertBoolsEqual(t, true, v.Purpose(2))
		assertBoolsEqual(t, true, v.PurposeStrict(2))
		assertBoolsEqual(t, false, v.Purpose(3))
		assertBoolsEqual(t, false, v.PurposeStrict(3))
		assertBoolsEqual(t, true, v.Purpose(4))
		assertBoolsEqual(t, true, v.PurposeStrict(4))
		assertBoolsEqual(t, false, v.Purpose(5))
		assertBoolsEqual(t, false, v.PurposeStrict(5))
		assertBoolsEqual(t, false, v.Purpose(6))
		assertBoolsEqual(t, false, v.PurposeStrict(6))

		assertBoolsEqual(t, false, v.LegitimateInterest(1))
		assertBoolsEqual(t, false, v.LegitimateInterestStrict(1))
		assertBoolsEqual(t, true, v.LegitimateInterest(2))
		assertBoolsEqual(t, false, v.LegitimateInterestStrict(2))
		assertBoolsEqual(t, false, v.LegitimateInterest(3))
		assertBoolsEqual(t, false, v.LegitimateInterestStrict(3))

		assertBoolsEqual(t, false, v.SpecialPurpose(1))
		assertBoolsEqual(t, false, v.SpecialPurpose(2))
		assertBoolsEqual(t, false, v.SpecialPurpose(3)) // Does not exist yet
		
		assertBoolsEqual(t, false, v.SpecialFeature(1))
		assertBoolsEqual(t, false, v.SpecialFeature(2))
		assertBoolsEqual(t, false, v.SpecialFeature(3)) // Does not exist yet
	}

}

const testData = `
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
			"specialFeatures": [1, 2],
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

func assertNil(t *testing.T, value api.Vendor, expectNil bool) {
	t.Helper()
	if expectNil && value != nil {
		t.Error("The vendor should be nil, but wasn't.")
	}
	if !expectNil && value == nil {
		t.Errorf("The vendor should not be nil, but was.")
	}
}
