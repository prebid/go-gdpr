package vendorconsent

import (
	"testing"
)

func TestPubRestrictions(t *testing.T) {
	baseConsent, err := Parse(decode(t, "COwAdDhOwAdDhN4ABAENAPCgAAQAAv___wAAAFP_AAp_4AI6ACACAA"))
	assertNilError(t, err)

	consent := baseConsent.(ConsentMetadata)
	// Extra random verification to ensure the basics are solid
	assertUInt8sEqual(t, 2, consent.Version())
	assertUInt16sEqual(t, 888, consent.CmpID())
	assertUInt16sEqual(t, 1, consent.CmpVersion())
	assertUInt8sEqual(t, 0, consent.ConsentScreen())
	assertStringsEqual(t, "EN", consent.ConsentLanguage())
	assertUInt16sEqual(t, 15, consent.VendorListVersion())
	assertUInt16sEqual(t, 10, consent.MaxVendorID())

	// A pub restriction was set on Puropse 7, type 1, for vendor 32
	assertBoolsEqual(t, true, consent.CheckPubRestriction(7, 1, 32))
	// Verify a that a different vendor is not flagged
	assertBoolsEqual(t, false, consent.CheckPubRestriction(7, 1, 7))
	// Verify that a different purpose is not flagged for the same vendor
	assertBoolsEqual(t, false, consent.CheckPubRestriction(5, 1, 32))
}

func TestPubRestrictions2(t *testing.T) {
	baseConsent, err := Parse(decode(t, "COxPe2TOxPe2TALABAENAPCgAAAAAAAAAAAAAFAAAAoAAA4IACACAIABgACAFA4ADACAAIygAGADwAQBIAIAIB0AEAEBSACACAA"))
	assertNilError(t, err)

	consent := baseConsent.(ConsentMetadata)

	assertBoolsEqual(t, false, consent.PurposeOneTreatment())
	assertBoolsEqual(t, false, consent.SpecialFeatureOptIn(3))

	// Pub restriction 0 set on purpose 1
	assertBoolsEqual(t, true, consent.CheckPubRestriction(1, 0, 32))
	// Verify that restriction type 1 isn't erroneously flagged
	assertBoolsEqual(t, false, consent.CheckPubRestriction(1, 1, 32))
	// Pub restriction 0 for purpose 2 is set on several vendors, but not all
	assertBoolsEqual(t, true, consent.CheckPubRestriction(2, 0, 32))
	assertBoolsEqual(t, true, consent.CheckPubRestriction(2, 0, 5))
	assertBoolsEqual(t, true, consent.CheckPubRestriction(2, 0, 11))
	assertBoolsEqual(t, false, consent.CheckPubRestriction(2, 0, 44))
	assertBoolsEqual(t, false, consent.CheckPubRestriction(2, 0, 500))
	// Pub restriction 1 on purpose 2 also set for vendor 32. This does not make sense, but not explicitly
	// disallowed by the spec. Verifying that the code handles it
	assertBoolsEqual(t, true, consent.CheckPubRestriction(2, 1, 32))
	// Verifying that this special case doesn't bleed over to other vendors.
	assertBoolsEqual(t, false, consent.CheckPubRestriction(2, 1, 42))
}

func TestPubRestrictions3(t *testing.T) {
	_, err := Parse(decode(t, "COzSDo9OzSDo9B9AAAENAiCAALAAAAAAAAAACOQAQCOAAAAA"))
	assertNilError(t, err)
}
