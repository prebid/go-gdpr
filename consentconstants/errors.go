package consentconstants

import "errors"

// ErrEmptyDecodedConsent error raised when the consent string is empty
var ErrEmptyDecodedConsent = errors.New("decoded consent cannot be empty")
