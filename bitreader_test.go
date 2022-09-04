package bitreader

import (
	"testing"
)

// 01110001, 00001101, 00000000, 00000000, 10100010, 00011011, 00000000, 00000000, 11001100
var TestArray = [...]byte{113, 13, 0, 0, 162, 27, 0, 0, 204}

func TestReadBit(t *testing.T) {
	bitreader := Reader(TestArray[:], false)
	expected := []bool{false, true, true, true}
	for i := range expected {
		value, err := bitreader.ReadBit()
		if err != nil {
			t.Fatal(err)
		}
		if value != expected[i] {
			t.Fatalf("ReadBit FAIL for index %d: Expected %t, Got %t", i, expected[i], value)
		}
	}
}

func TestReadBitLE(t *testing.T) {
	bitreader := Reader(TestArray[:], true)
	expected := []bool{true, false, false, false}
	for i := range expected {
		value, err := bitreader.ReadBit()
		if err != nil {
			t.Fatal(err)
		}
		if value != expected[i] {
			t.Fatalf("ReadBitLSB FAIL for index %d: Expected %t, Got %t", i, expected[i], value)
		}
	}
}

func TestReadBits32(t *testing.T) {
	bitreader := Reader(TestArray[:], false)
	expected := []int{3793354753, 2288779267} // 11100010000110100000000000000001, 10001000011011000000000000000011
	expectedBool := []bool{false, false}
	for i := range expected {
		bool, err := bitreader.ReadBit()
		if bool != expectedBool[i] {
			t.Fatalf("ReadBits32 ReadBit FAIL for index %d: Expected %t, Got %t", i, expectedBool[i], bool)
		}
		if err != nil {
			t.Fatal(err)
		}
		value, err := bitreader.ReadBits(32)
		if err != nil {
			t.Fatal(err)
		}
		if value != expected[i] {
			t.Fatalf("ReadBits32 FAIL for index %d: Expected %d, Got %d", i, expected[i], value)
		}
	}
}

func TestReadBitsLE(t *testing.T) {
	bitreader := Reader(TestArray[:], true)
	expected := []int{1720, 1768} // 11010111000, 11011101000
	for i := range expected {
		bitreader.ReadBit()
		value, err := bitreader.ReadBits(32)
		if err != nil {
			t.Fatal(err)
		}
		if value != expected[i] {
			t.Fatalf("ReadBits32LSB FAIL for index %d: Expected %d, Got %d", i, expected[i], value)
		}
	}
}

func TestSkipBits(t *testing.T) {
	bitreader := Reader(TestArray[:], false)
	expected := []bool{true, true, false, true} //00001101
	err := bitreader.SkipBits(12)
	if err != nil {
		t.Fatal(err)
	}
	for i := range expected {
		value, err := bitreader.ReadBit()
		if err != nil {
			t.Fatal(err)
		}
		if value != expected[i] {
			t.Fatalf("SkipBits ReadBit FAIL for index %d: Expected %t, Got %t", i, expected[i], value)
		}
	}
}

func TestSkipBitsLE(t *testing.T) {
	bitreader := Reader(TestArray[:], true)
	expected := []bool{false, false, false, false} //10110000
	bitreader.SkipBits(12)
	for i := range expected {
		value, err := bitreader.ReadBit()
		if err != nil {
			t.Fatal(err)
		}
		if value != expected[i] {
			t.Fatalf("SkipBits ReadBit FAIL for index %d: Expected %t, Got %t", i, expected[i], value)
		}
	}
}
