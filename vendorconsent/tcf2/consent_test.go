package vendorconsent

import (
	"testing"
)

func TestParseLegitIntSetWithBitField(t *testing.T) {
	// this test uses a crafted consent uses bit field, declares 10 vendors and legitimate interest without required content
	_, err := Parse(decode(t, "COvcSpYOvcSpYC9AAAENAPCAAAAAAAAAAAAAAFAAAAA"))
	assertError(t, err)
}

func TestParseLegitIntSetWithRangeSection(t *testing.T) {
	// this test uses a crafted consent uses range section, declares 10 vendors, 6 exceptions and legitimate interest without required content
	_, err := Parse(decode(t, "COvcSpYOvcSpYC9AAAENAPCAAAAAAAAAAAAAAFQBgAAgABAACAAEAAQAAgAA"))
	assertError(t, err)
}
