package bitreader

import (
	"bytes"
	"testing"
)

// TODO: Write better unit tests

// 01110001, 00001101, 00000000, 00000000, 10100010, 00011011, 00000000, 00000000, 11001100
var TestArray = [...]byte{113, 13, 0, 0, 162, 27, 0, 0, 204}

func TestTryReadFloat32(t *testing.T) {
	bitreader := Reader(bytes.NewReader(TestArray[:]), false)
	expected := []float32{6.98198182157e+29, -2.10064170919e-18}
	for i := range expected {
		value := bitreader.TryReadFloat32()
		if value != expected[i] {
			t.Fatalf("TryReadFloat32 FAIL for index %d: Expected %f, Got %f", i, expected[i], value)
		}
	}
}

func TestTryReadFloat64(t *testing.T) {
	bitreader := Reader(bytes.NewReader(TestArray[:]), false)
	expected := []float64{3.68828741038253948851462939603e+236}
	for i := range expected {
		value := bitreader.TryReadFloat64()
		if value != expected[i] {
			t.Fatalf("TryReadFloat64 FAIL for index %d: Expected %f, Got %f", i, expected[i], value)
		}
	}
}

func TestTryReadInt8(t *testing.T) {
	bitreader := Reader(bytes.NewReader(TestArray[:]), false)
	expected := []int{113, 13, 0}
	for i := range expected {
		value := bitreader.TryReadInt8()
		if int(value) != expected[i] {
			t.Fatalf("TryReadInt8 FAIL for index %d: Expected %d, Got %d", i, expected[i], value)
		}
	}
}

func TestTryReadInt16(t *testing.T) {
	bitreader := Reader(bytes.NewReader(TestArray[:]), false)
	expected := []int{28941, 0, 41499, 0}
	for i := range expected {
		value := bitreader.TryReadInt16()
		if int(value) != expected[i] {
			t.Fatalf("TryReadInt16 FAIL for index %d: Expected %d, Got %d", i, expected[i], value)
		}
	}
}

func TestTryReadInt32(t *testing.T) {
	bitreader := Reader(bytes.NewReader(TestArray[:]), false)
	expected := []int{1896677376, 2719678464}
	for i := range expected {
		value := bitreader.TryReadInt32()
		if int(value) != expected[i] {
			t.Fatalf("TryReadInt32 FAIL for index %d: Expected %d, Got %d", i, expected[i], value)
		}
	}
}

func TestTryReadInt64(t *testing.T) {
	bitreader := Reader(bytes.NewReader(TestArray[:]), false)
	expected := []int{8146167303702773760}
	for i := range expected {
		value := bitreader.TryReadInt64()
		if int(value) != expected[i] {
			t.Fatalf("TryReadInt64 FAIL for index %d: Expected %d, Got %d", i, expected[i], value)
		}
	}
}

func TestReadBit(t *testing.T) {
	bitreader := Reader(bytes.NewReader(TestArray[:]), false)
	expected := []bool{false, true, true, true}
	for i := range expected {
		value, err := bitreader.ReadBool()
		if err != nil {
			t.Fatal(err)
		}
		if value != expected[i] {
			t.Fatalf("ReadBit FAIL for index %d: Expected %t, Got %t", i, expected[i], value)
		}
	}
}

func TestReadBitLE(t *testing.T) {
	bitreader := Reader(bytes.NewReader(TestArray[:]), true)
	expected := []bool{true, false, false, false}
	for i := range expected {
		value, err := bitreader.ReadBool()
		if err != nil {
			t.Fatal(err)
		}
		if value != expected[i] {
			t.Fatalf("ReadBitLE FAIL for index %d: Expected %t, Got %t", i, expected[i], value)
		}
	}
}

func TestReadBits(t *testing.T) {
	bitreader := Reader(bytes.NewReader(TestArray[:]), false)
	expected := []int{3793354753, 2288779267} // 11100010000110100000000000000001, 10001000011011000000000000000011
	expectedBool := []bool{false, false}
	for i := range expected {
		bool, err := bitreader.ReadBool()
		if bool != expectedBool[i] {
			t.Fatalf("ReadBits ReadBit FAIL for index %d: Expected %t, Got %t", i, expectedBool[i], bool)
		}
		if err != nil {
			t.Fatal(err)
		}
		value, err := bitreader.ReadBits(32)
		if err != nil {
			t.Fatal(err)
		}
		if int(value) != expected[i] {
			t.Fatalf("ReadBits FAIL for index %d: Expected %d, Got %d", i, expected[i], value)
		}
	}
}

func TestReadBitsLE(t *testing.T) {
	bitreader := Reader(bytes.NewReader(TestArray[:]), true)
	expected := []int{1720, 1768} // 11010111000, 11011101000
	for i := range expected {
		bitreader.ReadBool()
		value, err := bitreader.ReadBits(32)
		if err != nil {
			t.Fatal(err)
		}
		if int(value) != expected[i] {
			t.Fatalf("ReadBitsLE FAIL for index %d: Expected %d, Got %d", i, expected[i], value)
		}
	}
}

func TestReadBytes(t *testing.T) {
	bitreader := Reader(bytes.NewReader(TestArray[:]), false)
	expected := []int{3793354753, 2288779267} // 11100010000110100000000000000001, 10001000011011000000000000000011
	expectedBool := []bool{false, false}
	for i := range expected {
		bool, err := bitreader.ReadBool()
		if bool != expectedBool[i] {
			t.Fatalf("ReadBytes ReadBit FAIL for index %d: Expected %t, Got %t", i, expectedBool[i], bool)
		}
		if err != nil {
			t.Fatal(err)
		}
		value, err := bitreader.ReadBytes(4)
		if err != nil {
			t.Fatal(err)
		}
		if int(value) != expected[i] {
			t.Fatalf("ReadBytes FAIL for index %d: Expected %d, Got %d", i, expected[i], value)
		}
	}
}

func TestReadBytesLE(t *testing.T) {
	bitreader := Reader(bytes.NewReader(TestArray[:]), true)
	expected := []int{1720, 1768} // 11010111000, 11011101000
	for i := range expected {
		bitreader.ReadBool()
		value, err := bitreader.ReadBytes(4)
		if err != nil {
			t.Fatal(err)
		}
		if int(value) != expected[i] {
			t.Fatalf("ReadBytesLE FAIL for index %d: Expected %d, Got %d", i, expected[i], value)
		}
	}
}

func TestSkipBits(t *testing.T) {
	bitreader := Reader(bytes.NewReader(TestArray[:]), false)
	expected := []bool{true, true, false, true} //00001101
	err := bitreader.SkipBits(12)
	if err != nil {
		t.Fatal(err)
	}
	for i := range expected {
		value, err := bitreader.ReadBool()
		if err != nil {
			t.Fatal(err)
		}
		if value != expected[i] {
			t.Fatalf("SkipBits ReadBit FAIL for index %d: Expected %t, Got %t", i, expected[i], value)
		}
	}
}

func TestSkipBitsLE(t *testing.T) {
	bitreader := Reader(bytes.NewReader(TestArray[:]), true)
	expected := []bool{false, false, false, false} //10110000
	bitreader.SkipBits(12)
	for i := range expected {
		value, err := bitreader.ReadBool()
		if err != nil {
			t.Fatal(err)
		}
		if value != expected[i] {
			t.Fatalf("SkipBitsLE ReadBit FAIL for index %d: Expected %t, Got %t", i, expected[i], value)
		}
	}
}
