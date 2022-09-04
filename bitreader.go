package bitreader

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

type ReaderType struct {
	data  []byte // Reader data from a byte array
	base  int    // Current reader byte location
	index int    // Current reader index
	le    bool   // Little endian or big endian?
}

func Reader(data []byte, le bool) *ReaderType {
	dataClone := data
	if le {
		for index, byteValue := range data {
			dataClone[index] = bits.Reverse8(byteValue)
		}
	}
	return &ReaderType{
		data:  dataClone,
		base:  0,
		index: 0,
		le:    le,
	}
}

func (reader *ReaderType) SkipBits(bits int) error {
	if bits <= 0 {
		return fmt.Errorf("SkipBits Error: Bits value %d lower or equals than 0.", bits)
	}
	err := reader.checkAvailableBits(bits)
	if err != nil {
		return err
	}
	for reader.index+bits > 7 {
		reader.base++
		reader.index = 0
		bits -= 8
	}
	reader.index += bits
	return nil
}

func (reader *ReaderType) SkipBytes(bytes int) error {
	err := reader.SkipBits(bytes * 8)
	if err != nil {
		return err
	}
	return nil
}

func (reader *ReaderType) ReadBits(bits int) (int, error) {
	if bits <= 0 {
		return -1, fmt.Errorf("ReadBits Error: Bits value %d lower or equals than 0.", bits)
	}
	if bits > 64 {
		return -1, fmt.Errorf("ReadBits Error: Bits value %d higher than 64.", bits)
	}
	err := reader.checkAvailableBits(bits)
	if err != nil {
		return -1, err
	}
	if reader.le {
		var output string
		// Go to last bit and read backwards from there
		reader.base += bits / 8
		reader.index += bits % 8
		if reader.index > 7 {
			reader.index -= 8
			reader.base++
		}
		for i := 0; i < bits; i++ {
			reader.index--
			if reader.index < 0 {
				reader.base--
				reader.index = 7
			}
			binary := fmt.Sprintf("%08b", reader.data[reader.base])
			binaryArr := strings.Split(binary, "")
			output += binaryArr[reader.index]
		}
		// Return to last bit after reading
		reader.base += bits / 8
		reader.index += bits % 8
		if reader.index > 7 {
			reader.index -= 8
		}
		// Conversion of string binary to int
		value, err := strconv.ParseUint(output, 2, 64)
		if err != nil {
			return -1, fmt.Errorf("%s", err)
		}
		return int(value), nil
	} else {
		var output string
		for i := 0; i < bits; i++ {
			binary := fmt.Sprintf("%08b", reader.data[reader.base])
			binaryArr := strings.Split(binary, "")
			output += binaryArr[reader.index]
			reader.index++
			if reader.index > 7 {
				reader.base++
				reader.index = 0
			}
		}
		// Conversion of string binary to int
		value, err := strconv.ParseUint(output, 2, 64)
		if err != nil {
			return -1, fmt.Errorf("%s", err)
		}
		return int(value), nil
	}
}

func (reader *ReaderType) ReadBit() (bool, error) {
	value, err := reader.ReadBits(1)
	if err != nil {
		return false, fmt.Errorf("ReadBit Error: %s", err)
	}
	return value != 0, nil
}

func (reader *ReaderType) checkAvailableBits(bits int) error {
	availableBits := (len(reader.data)-reader.base)*8 - reader.index
	if availableBits < bits {
		return fmt.Errorf("BitReaderOutOfBounds: Wanted to read/skip %d bit(s) but only %d bit(s) is/are available.", bits, availableBits)
	}
	return nil
}
