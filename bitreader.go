// paksd
package bitreader

import (
	"fmt"
	"io"
	"math"
)

type ReaderType struct {
	stream  io.Reader // the underlying stream we're reading bytes from
	index   uint8     // 0-7, the current index into the byte
	curByte byte      // the byte we're currently reading from
	le      bool      // whether to read in little-endian order
	// Whenever index == 0, we need to read a new byte from stream into curByte
}

func Reader(stream io.Reader, le bool) *ReaderType {
	return &ReaderType{
		stream:  stream,
		index:   0,
		curByte: 0, // initial value doesn't matter, it'll be read as soon as we try to read any bits
		le:      le,
	}
}

func (reader *ReaderType) TryReadBool() bool {
	flag, err := reader.ReadBool()
	if err != nil {
		panic(err)
	}
	return flag
}

func (reader *ReaderType) TryReadInt1() uint8 {
	value, err := reader.ReadBits(1)
	if err != nil {
		panic(err)
	}
	return uint8(value)
}

func (reader *ReaderType) TryReadInt8() uint8 {
	value, err := reader.ReadBits(8)
	if err != nil {
		panic(err)
	}
	return uint8(value)
}

func (reader *ReaderType) TryReadInt16() uint16 {
	value, err := reader.ReadBits(16)
	if err != nil {
		panic(err)
	}
	return uint16(value)
}

func (reader *ReaderType) TryReadInt32() uint32 {
	value, err := reader.ReadBits(32)
	if err != nil {
		panic(err)
	}
	return uint32(value)
}

func (reader *ReaderType) TryReadInt64() uint64 {
	value, err := reader.ReadBits(64)
	if err != nil {
		panic(err)
	}
	return value
}

func (reader *ReaderType) TryReadFloat32() float32 {
	value, err := reader.ReadBits(32)
	if err != nil {
		panic(err)
	}
	return math.Float32frombits(uint32(value))
}

func (reader *ReaderType) TryReadFloat64() float64 {
	value, err := reader.ReadBits(64)
	if err != nil {
		panic(err)
	}
	return math.Float64frombits(value)
}

func (reader *ReaderType) SkipBits(bits int) error {
	// read as many raw bytes as we can
	bytes := bits / 8
	buf := make([]byte, bytes)
	_, err := reader.stream.Read(buf)
	if err != nil {
		return err
	}
	// the final read byte should be the new current byte
	if bytes > 0 {
		reader.curByte = buf[bytes-1]
	}
	// read the extra bits
	for i := bytes * 8; i < bits; i++ {
		_, err := reader.readBit()
		if err != nil {
			return err
		}
	}
	return nil
}

func (reader *ReaderType) SkipBytes(bytes int) error {
	err := reader.SkipBits(bytes * 8)
	if err != nil {
		return err
	}
	return nil
}

// Read up to 64 bits from the stream
func (reader *ReaderType) ReadBits(bits int) (uint64, error) {
	if bits < 1 || bits > 64 {
		return 0, fmt.Errorf("ReadBits(bits) ERROR: Bits number should be between 1 and 64.")
	}
	var val uint64
	for i := 0; i < bits; i++ {
		bit, err := reader.readBit()
		if err != nil {
			return 0, err
		}

		if reader.le {
			val |= uint64(bit) << i
		} else {
			val |= uint64(bit) << (bits - 1 - i)
		}
	}
	return val, nil
}

func (reader *ReaderType) ReadBytes(bytes int) (uint64, error) {
	if bytes < 1 || bytes > 8 {
		return 0, fmt.Errorf("ReadBytes(bytes) ERROR: Bytes number should be between 1 and 8.")
	}
	value, err := reader.ReadBits(bytes * 8)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// Read a single bool from the stream
func (reader *ReaderType) ReadBool() (bool, error) {
	val, err := reader.readBit()
	if err != nil {
		return false, err
	}
	return val == 1, nil
}

// Read a single bit from the stream
func (reader *ReaderType) readBit() (uint8, error) {
	if reader.index == 0 {
		// read a byte from stream into curByte
		buf := make([]byte, 1)
		_, err := reader.stream.Read(buf)
		if err != nil {
			return 0, err
		}
		reader.curByte = buf[0]
	}
	var val bool
	if reader.le {
		val = (reader.curByte & (1 << reader.index)) != 0
	} else {
		val = (reader.curByte & (1 << (7 - reader.index))) != 0
	}
	reader.index = (reader.index + 1) % 8
	if val {
		return 1, nil
	} else {
		return 0, nil
	}
}
