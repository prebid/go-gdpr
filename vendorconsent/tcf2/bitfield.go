package vendorconsent

import (
	"fmt"
)

func parseBitField(metadata consentMetadata) (*consentBitField, uint, error) {
	data := metadata.data
	vendorBitsRequired := metadata.MaxVendorID()

	// BitFields start at bit 230. This means the last three bits of byte 28 are part of the bitfield.
	// In this case "others" will never be used, and we don't risk an index-out-of-bounds by using it.
	if vendorBitsRequired <= 3 {
		return &consentBitField{
			maxVendorID: vendorBitsRequired,
			firstTwo:    data[28],
			others:      nil,
		}, uint(230 + vendorBitsRequired), nil
	}

	otherBytesRequired := (vendorBitsRequired - 3) / 8
	if (vendorBitsRequired-3)%8 > 0 {
		otherBytesRequired = otherBytesRequired + 1
	}
	dataLengthRequired := 28 + otherBytesRequired
	if uint(len(data)) < uint(dataLengthRequired) {
		return nil, 0, fmt.Errorf("a BitField for %d vendors requires a consent string of %d bytes. This consent string had %d", vendorBitsRequired, dataLengthRequired, len(data))
	}

	return &consentBitField{
		maxVendorID: vendorBitsRequired,
		firstTwo:    data[28],
		others:      data[29:],
	}, uint(230 + vendorBitsRequired), nil
}

// A BitField has len(MaxVendorID()) entries, with one bit for every vendor in the range.
type consentBitField struct {
	maxVendorID uint16
	firstTwo    byte
	others      []byte
}

func (f *consentBitField) VendorConsent(id uint16) bool {
	if id < 1 || id > f.maxVendorID {
		return false
	}
	// Careful here... vendor IDs start at index 1...
	if id <= 3 {
		return byteToBool(f.firstTwo & (0x04 >> id))
	}
	return isSet(f.others, uint(id-3))
}

// byteToBool returns false if val is 0, and true otherwise
func byteToBool(val byte) bool {
	return val != 0
}
