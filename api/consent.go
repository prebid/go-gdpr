package api

import (
	"time"

	"github.com/prebid/go-gdpr/consentconstants"
)

// VendorConsents is a GDPR Vendor Consent string, as defined by IAB Europe. For technical details,
// see https://github.com/InteractiveAdvertisingBureau/GDPR-Transparency-and-Consent-Framework/blob/master/Consent%20string%20and%20vendor%20list%20formats%20v1.1%20Final.md#vendor-consent-string-format-
type VendorConsents interface {
	// The version of the Consent string.
	Version() uint8

	// The time that the consent string was first created
	Created() time.Time

	// The time that the consent string was last updated
	LastUpdated() time.Time

	// The ID of the CMP used to update the consent string.
	CmpID() uint16

	// The version of the CMP used to update the consent string
	CmpVersion() uint16

	// The number of the CMP screen where consent was given
	ConsentScreen() uint8

	// The two-letter ISO639-1 language code used by the CMP to ask for consent, in uppercase.
	ConsentLanguage() string

	// The VendorListVersion which is needed to interpret this consent string.
	//
	// The IAB is hosting these on their webpage. For example, version 2 of the
	// Vendor List can be found at https://vendorlist.consensu.org/v-2/vendorlist.json
	//
	// For other versions, just replace the "v-*" path with the value returned here.
	// The latest version can always be found at https://vendorlist.consensu.org/vendorlist.json
	VendorListVersion() uint16

	// TCFPolicyVersion indicates the TCF policy version needed to interpret this consent string.
	TCFPolicyVersion() uint8

	// MaxVendorID describes how many vendors are encoded into the string.
	// This is the upper bound (inclusive) on valid inputs for HasConsent(id).
	MaxVendorID() uint16

	// Determine if the user has consented to use data for the given Purpose.
	//
	// If the purpose is converted from an int > 24, the return value is undefined because
	// the consent string doesn't have room for more purposes than that.
	PurposeAllowed(id consentconstants.Purpose) bool

	// Determine if a given vendor has consent to collect or receive user info.
	//
	// This function's behavior is undefined for "invalid" IDs.
	// IDs with value < 1 or value > MaxVendorID() are definitely invalid, but IDs within that range
	// may still be invalid, depending on the Vendor List.
	//
	// It is the caller's responsibility to get the right Vendor List version for the semantics of the ID.
	// For more information, see VendorListVersion().
	VendorConsent(id uint16) bool

	// VendorDisclosed determines if a given vendor was disclosed to the user.
	// This is mandatory in TCF 2.3 and is used to verify that vendors (especially those
	// declaring only Special Purposes) were actually presented to the user.
	// Mandatory Inclusion: All TC Strings created on or after February 28, 2026, must include the disclosedVendors segment.
	// Source https://iabeurope.eu/all-you-need-to-know-about-the-transition-to-tcf-v2-3/#:~:text=Any%20TC%20String%20created%20without,TC%20String%20created%20under%202.3.
	//
	// For TCF 2.2 strings that don't include a DisclosedVendors segment,
	// this returns false to maintain backward compatibility.
	VendorDisclosed(id uint16) bool

	VendorDisclosedMaxVendorID() uint16

	// HasDisclosedVendors returns true if the consent string includes a disclosedVendors segment.
	// This method is particularly useful during the TCF 2.3 transition phase (mandatory from March 1st, 2025)
	// to distinguish between:
	//   - Legacy TCF 2.0/2.2 consent strings that legitimately don't have this segment (returns false)
	//   - TCF 2.3+ consent strings that should have this segment (returns true if present, false if malformed)
	//
	// Note: VendorDisclosed(id) returns false both when the segment is missing AND when
	// a vendor is not disclosed, so use HasDisclosedVendors() to disambiguate these cases.
	HasDisclosedVendors() bool
}
