package vendorconsent

import "fmt"

func parseBitField20(data consentMetadata20) (*consentBitField20, error) {
	vendorBitsRequired := data.MaxVendorID()

	// BitFields start at bit 230. This means the last two bits of byte 28 are part of the bitfield.
	// In this case "others" will never be used, and we don't risk an index-out-of-bounds by using it.
	if vendorBitsRequired <= 3 {
		return &consentBitField20{
			consentMetadata20: data,
			firstTwo:          data[28],
			others:            nil,
		}, nil
	}

	otherBytesRequired := (vendorBitsRequired - 2) / 8
	if (vendorBitsRequired-2)%8 > 0 {
		otherBytesRequired = otherBytesRequired + 1
	}
	dataLengthRequired := 29 + otherBytesRequired
	if uint(len(data)) < uint(dataLengthRequired) {
		return nil, fmt.Errorf("a BitField for %d vendors requires a consent string of %d bytes. This consent string had %d", vendorBitsRequired, dataLengthRequired, len(data))
	}

	return &consentBitField20{
		consentMetadata20: data,
		firstTwo:          data[28],
		others:            data[29:],
	}, nil
}

// A BitField has len(MaxVendorID()) entries, with one bit for every vendor in the range.
type consentBitField20 struct {
	consentMetadata20
	firstTwo byte
	others   []byte
}

func (f *consentBitField20) VendorConsent(id uint16) bool {
	if id < 1 || id > f.MaxVendorID() {
		return false
	}
	// Careful here... vendor IDs start at index 1...
	if id <= 2 {
		return byteToBool(f.firstTwo & (0x04 >> id))
	}
	return isSet(f.others, uint(id-3))
}
