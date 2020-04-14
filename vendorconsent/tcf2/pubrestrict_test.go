package vendorconsent

import (
	"testing"
)

func TestPubRestrictions(t *testing.T) {
	baseConsent, err := Parse(decode(t, "COwAdDhOwAdDhN4ABAENAPCgAAQAAv___wAAAFP_AAp_4AI6ACACAA"))
	assertNilError(t, err)
	consent := baseConsent.(ConsentMetadata)
	assertUInt8sEqual(t, 2, consent.Version())
	assertUInt16sEqual(t, 888, consent.CmpID())
	assertUInt16sEqual(t, 1, consent.CmpVersion())
	assertUInt8sEqual(t, 0, consent.ConsentScreen())
	assertStringsEqual(t, "EN", consent.ConsentLanguage())
	assertUInt16sEqual(t, 15, consent.VendorListVersion())
	assertUInt16sEqual(t, 10, consent.MaxVendorID())

	assertBoolsEqual(t, true, consent.CheckPubRestriction(7, 1, 32))
	assertBoolsEqual(t, false, consent.CheckPubRestriction(7, 1, 7))
	assertBoolsEqual(t, false, consent.CheckPubRestriction(5, 1, 32))

	baseConsent, err = Parse(decode(t, "COxPe2TOxPe2TALABAENAPCgAAAAAAAAAAAAAFAAAAoAAA4IACACAIABgACAFA4ADACAAIygAGADwAQBIAIAIB0AEAEBSACACAA"))
	assertNilError(t, err)
	consent = baseConsent.(ConsentMetadata)

	assertBoolsEqual(t, false, consent.PurposeOneTreatment())
	assertBoolsEqual(t, false, consent.SpecialFeatureOptIn(3))
	assertBoolsEqual(t, true, consent.CheckPubRestriction(1, 0, 32))
	assertBoolsEqual(t, false, consent.CheckPubRestriction(1, 1, 32))
	assertBoolsEqual(t, true, consent.CheckPubRestriction(2, 0, 32))
	assertBoolsEqual(t, true, consent.CheckPubRestriction(2, 0, 5))
	assertBoolsEqual(t, true, consent.CheckPubRestriction(2, 0, 11))
	assertBoolsEqual(t, false, consent.CheckPubRestriction(2, 0, 44))
	assertBoolsEqual(t, false, consent.CheckPubRestriction(2, 0, 500))
	assertBoolsEqual(t, true, consent.CheckPubRestriction(2, 1, 32))
	assertBoolsEqual(t, false, consent.CheckPubRestriction(2, 1, 42))
}
