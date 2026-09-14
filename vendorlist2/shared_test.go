package vendorlist2

import (
	"testing"

	"github.com/prebid/go-gdpr/api"
)

func AssertVendorListCorrectness(t *testing.T, gvl api.VendorList) {
	v := gvl.Vendor(8)
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

	v = gvl.Vendor(80)
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

const testDataSpecVersion2 = `
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

const testDataSpecVersion2Empty = `
{
	"gvlSpecificationVersion": 2,
	"vendorListVersion": 28,
	"tcfPolicyVersion": 2,
	"lastUpdated": "2020-03-05T16:05:29Z",
	"vendors": { }
}
`

const testDataSpecVersion3 = `
{
	"gvlSpecificationVersion": 3,
	"vendorListVersion": 1,
	"tcfPolicyVersion": 4,
	"lastUpdated": "2023-05-18T16:07:14Z",
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
			"dataRetention": {
				"stdRetention": 30,
				"purposes": { "9": 180 },
				"specialPurposes": {}
			},
		   "dataDeclaration" : [ 1, 2, 4, 6 ],
		   "urls": [
			  {
				"langId": "en",
				"privacy": "https://vendorname.com/gdpr.html",
				"legIntClaim": "https://vendorname.com/gdpr.html#li"
			  },
			  {
				"langId": "fr",
				"privacy": "https://vendorname.com/fr/gdpr.html",
				"legIntClaim": "https://vendorname.com/fr/gdpr.html#li"
			  }
		   ]
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
			"dataRetention": {
				"stdRetention": 30,
				"purposes": { "9": 180 },
				"specialPurposes": {}
			},
		   "dataDeclaration" : [ 1, 2, 4, 6 ],
		   "urls": [
			  {
				"langId": "en",
				"privacy": "https://vendorname.com/gdpr.html",
				"legIntClaim": "https://vendorname.com/gdpr.html#li"
			  },
			  {
				"langId": "fr",
				"privacy": "https://vendorname.com/fr/gdpr.html",
				"legIntClaim": "https://vendorname.com/fr/gdpr.html#li"
			  }
		   ]
		}
	}
}
`

const testDataSpecVersion3Empty = `
{
	"gvlSpecificationVersion": 3,
	"vendorListVersion": 1,
	"tcfPolicyVersion": 4,
	"lastUpdated": "2023-05-18T16:07:14Z",
	"vendors": { }
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
