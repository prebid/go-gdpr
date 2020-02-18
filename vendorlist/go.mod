module github.com/prebid/go-gdpr/vendorlist

go 1.13

replace github.com/prebid/go-gdpr => ../.

replace github.com/prebid/go-gdpr/consentconstants => ../consentconstants

require (
	github.com/buger/jsonparser v0.0.0-20191204142016-1a29609e0929
	github.com/prebid/go-gdpr v0.6.0 // indirect
	github.com/prebid/go-gdpr/consentconstants v0.0.0-00010101000000-000000000000
)
