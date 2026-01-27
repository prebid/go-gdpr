package vendorconsent

import (
	"testing"

	"github.com/prebid/go-gdpr/consentconstants"
)

var consentStr = "CQeFq6FQeFq6FEZAAAENCZCAAL3AAEIAAAAAAHgACAB4AAgAAA.II7Nd_X__bX9n-_7_6ft0eY1f9_r37uQzDhfNs-8F3L_W_LwX32E7NF36tq4KmR4ku1bBIQNtHMnUDUmxaolVrzHsak2cpyNKJ_JkknsZe2dYGF9Pn9lD-YKZ7_5_9_f52T_9_9_-39z3_9f___dv_-__-vjf_599n_v9fV_78_Kf9______-____________8A"

func vendorTcfV2Allowed(consentStr string) bool {
	baseConsent, err := ParseString(consentStr)
	if err != nil {
		return false
	}
	consent := baseConsent.(*ConsentMetadata)
	return consent.VendorConsent(15) &&
		consent.VendorLegitInterest(15) &&
		consent.PurposeAllowed(consentconstants.Purpose(1)) &&
		consent.PurposeAllowed(consentconstants.Purpose(3)) &&
		consent.PurposeAllowed(consentconstants.Purpose(4)) &&
		consent.PurposeAllowed(consentconstants.Purpose(5)) &&
		consent.PurposeAllowed(consentconstants.Purpose(6)) &&
		consent.PurposeAllowed(consentconstants.Purpose(8)) &&
		consent.PurposeAllowed(consentconstants.Purpose(9)) &&
		consent.PurposeAllowed(consentconstants.Purpose(10)) &&
		consent.PurposeLITransparency(2) &&
		consent.PurposeLITransparency(7)
}

func BenchmarkVendorTcfV2Allowed(b *testing.B) {
	for b.Loop() {
		_ = vendorTcfV2Allowed(consentStr)
	}
}
