package vendorconsent

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/prebid/go-gdpr/consentconstants"
)

// Parse the metadata from the consent string.
// This returns an error if the input is too short to answer questions about that data.
func parseMetadata(data []byte) (consentMetadata, error) {
	if len(data) < 22 {
		return nil, fmt.Errorf("vendor consent strings are at least 22 bytes long. This one was %d", len(data))
	}
	metadata := consentMetadata(data)
	if metadata.MaxVendorID() < 1 {
		return nil, fmt.Errorf("the consent string encoded a MaxVendorID of %d, but this value must be greater than or equal to 1", metadata.MaxVendorID())
	}
	if metadata.Version() < 1 {
		return nil, fmt.Errorf("the consent string encoded a Version of %d, but this value must be greater than or equal to 1", metadata.Version())
	}
	if metadata.VendorListVersion() == 0 {
		return nil, errors.New("the consent string encoded a VendorListVersion of 0, but this value must be greater than or equal to 1")

	}
	return consentMetadata(data), nil
}

// consemtMetadata implements the parts of the VendorConsents interface which are common
// to BitFields and RangeSections. This relies on Parse to have done some validation already,
// to make sure that functions on it don't overflow the bounds of the byte array.
type consentMetadata []byte

func (c consentMetadata) Version() uint8 {
	return uint8(c[0] >> 2)
}

func (c consentMetadata) VendorListVersion() uint16 {
	// The vendor list version is stored in bits 120 - 131
	rightByte := ((c[16] & 0xf0) >> 4) | ((c[15] & 0x0f) << 4)
	leftByte := c[15] >> 4
	return binary.BigEndian.Uint16([]byte{leftByte, rightByte})
}

func (c consentMetadata) MaxVendorID() uint16 {
	// The max vendor ID is stored in bits 156 - 171
	leftByte := byte((c[19]&0x0f)<<4 + (c[20]&0xf0)>>4)
	rightByte := byte((c[20]&0x0f)<<4 + (c[21]&0xf0)>>4)
	return binary.BigEndian.Uint16([]byte{leftByte, rightByte})
}

func (c consentMetadata) PurposeAllowed(id consentconstants.Purpose) bool {
	// Purposes are stored in bits 132 - 155. The interface contract only defines behavior for ints in the range [1, 24]...
	// so in the valid range, this won't even overflow a uint8.
	return isSet(c, uint(id)+131)
}
