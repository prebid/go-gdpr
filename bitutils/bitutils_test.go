package bitutils

import (
	"testing"
)

// Define some test data

// 0000 0100 1010 0010 0000 0011 1011 0001 0000 0000 0010 1011

var testdata = []byte{0x04, 0xa2, 0x03, 0xb1, 0x00, 0x2b}

type testDefinition struct {
	data   []byte // The data to feed the function
	offset uint   // The bit offset in the byte slice to start
	value  uint64 // The value we expect the function to return (64 bit to allow for future functions that extract larger ints)
}

var test4Bits = []testDefinition{
	{testdata, 21, 7},
	{testdata, 12, 2},
	{testdata, 44, 11},
	{testdata, 6, 2},
	{[]byte{0x10}, 0, 1},
	{[]byte{0x92}, 4, 2},
	{[]byte{0x99}, 1, 3},
	{[]byte{0x01, 0xe0}, 7, 15},
}

func TestParseByte4(t *testing.T) {
	b, err := ParseByte4(testdata, 46)
	assertStringsEqual(t, "ParseByte4 expected 4 bits to start at bit 46, but the consent string was only 6 bytes long (needs second byte)", err.Error())

	b, err = ParseByte4(testdata, 80)
	assertStringsEqual(t, "ParseByte4 expected 4 bits to start at bit 80, but the consent string was only 6 bytes long", err.Error())

	for _, test := range test4Bits {
		b, err = ParseByte4(test.data, test.offset)
		assertNilError(t, err)
		assertBytesEqual(t, byte(test.value), b)
	}
}

func TestParseByte8(t *testing.T) {
	b, err := ParseByte8(testdata, 4)
	assertNilError(t, err)
	assertBytesEqual(t, 0x4a, b)
}

func TestParseUInt12(t *testing.T) {
	i, err := ParseUInt12(testdata, 10)
	assertNilError(t, err)
	assertUInt16sEqual(t, 2176, i)
	i, err = ParseUInt12(testdata, 16)
	assertNilError(t, err)
	assertUInt16sEqual(t, 59, i)
	i, err = ParseUInt12(testdata, 19)
	assertNilError(t, err)
	assertUInt16sEqual(t, 472, i)
	i, err = ParseUInt12(testdata, 1)
	assertNilError(t, err)
	assertUInt16sEqual(t, 148, i)
	i, err = ParseUInt12(testdata, 44)
	assertStringsEqual(t, "ParseUInt12 expected a 12-bit int to start at bit 44, but the consent string was only 6 bytes long", err.Error())
}

func assertNilError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func assertStringsEqual(t *testing.T, expected string, actual string) {
	t.Helper()
	if actual != expected {
		t.Errorf("Strings were not equal. Expected %s, actual %s", expected, actual)
	}
}

func assertBytesEqual(t *testing.T, expected byte, actual byte) {
	t.Helper()
	if actual != expected {
		t.Errorf("bytes were not equal. Expected %d, actual %d", expected, actual)
	}
}

func assertUInt8sEqual(t *testing.T, expected uint8, actual uint8) {
	t.Helper()
	if actual != expected {
		t.Errorf("Ints were not equal. Expected %d, actual %d", expected, actual)
	}
}

func assertUInt16sEqual(t *testing.T, expected uint16, actual uint16) {
	t.Helper()
	if actual != expected {
		t.Errorf("Ints were not equal. Expected %d, actual %d", expected, actual)
	}
}

func assertIntsEqual(t *testing.T, expected int, actual int) {
	t.Helper()
	if actual != expected {
		t.Errorf("Ints were not equal. Expected %d, actual %d", expected, actual)
	}
}

func assertBoolsEqual(t *testing.T, expected bool, actual bool) {
	t.Helper()
	if actual != expected {
		t.Errorf("Bools were not equal. Expected %t, actual %t", expected, actual)
	}
}
