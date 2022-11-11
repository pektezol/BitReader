// BitReader is a simple bit reader with big/little-endian support for golang.
// It can read stream data from an io.Reader; can read from os.File and a byte array with bytes.NewReader(array).
// Uses bitwise operations for v2.
// Supports reading up to 64 bits at one time.
// Includes wrapper functions for most used data types.
// Error checking on all but wrapper functions.
// Thanks to github.com/mlugg for the big help!
package bitreader

import (
	"fmt"
	"io"
	"math"
)

// ReaderType is the main structure of our Reader.
// Whenever index == 0, we need to read a new byte from stream into curByte
//
// stream io.Reader The underlying stream we're reading bytes from
// index uint18		The current index into the byte [0-7]
// curByte byte		The byte we're currently reading from
// le bool 			Whether to read in little-endian order
type ReaderType struct {
	stream  io.Reader
	index   uint8
	curByte byte
	le      bool
}

// Reader is the main constructor that creates the ReaderType object
// with stream data and little-endian state.
func Reader(stream io.Reader, le bool) *ReaderType {
	return &ReaderType{
		stream:  stream,
		index:   0,
		curByte: 0, // Initial value doesn't matter, it'll be read as soon as we try to read any bits
		le:      le,
	}
}

// TryReadBool is a wrapper function that gets the state of 1-bit,
// returns true if 1, false if 0. Panics on error.
func (reader *ReaderType) TryReadBool() bool {
	flag, err := reader.ReadBool()
	if err != nil {
		panic(err)
	}
	return flag
}

// TryReadInt1 is a wrapper function that returns the value of 1-bit.
// Returns type uint8. Panics on error.
func (reader *ReaderType) TryReadInt1() uint8 {
	value, err := reader.ReadBits(1)
	if err != nil {
		panic(err)
	}
	return uint8(value)
}

// TryReadInt8 is a wrapper function that returns the value of 8-bits.
// Returns uint8. Panics on error.
func (reader *ReaderType) TryReadInt8() uint8 {
	value, err := reader.ReadBits(8)
	if err != nil {
		panic(err)
	}
	return uint8(value)
}

// TryReadInt16 is a wrapper function that returns the value of 16-bits.
// Returns uint16. Panics on error.
func (reader *ReaderType) TryReadInt16() uint16 {
	value, err := reader.ReadBits(16)
	if err != nil {
		panic(err)
	}
	return uint16(value)
}

// TryReadInt32 is a wrapper function that returns the value of 32-bits.
// Returns uint32. Panics on error.
func (reader *ReaderType) TryReadInt32() uint32 {
	value, err := reader.ReadBits(32)
	if err != nil {
		panic(err)
	}
	return uint32(value)
}

// TryReadInt64 is a wrapper function that returns the value of 64-bits.
// Returns uint64. Panics on error.
func (reader *ReaderType) TryReadInt64() uint64 {
	value, err := reader.ReadBits(64)
	if err != nil {
		panic(err)
	}
	return value
}

// TryReadFloat32 is a wrapper function that returns the value of 32-bits.
// Returns float32. Panics on error.
func (reader *ReaderType) TryReadFloat32() float32 {
	value, err := reader.ReadBits(32)
	if err != nil {
		panic(err)
	}
	return math.Float32frombits(uint32(value))
}

// TryReadFloat64 is a wrapper function that returns the value of 64-bits.
// Returns float64. Panics on error.
func (reader *ReaderType) TryReadFloat64() float64 {
	value, err := reader.ReadBits(64)
	if err != nil {
		panic(err)
	}
	return math.Float64frombits(value)
}

// TryReadBits is a wrapper function that returns the value of bits specified in the parameter.
// Returns uint64. Panics on error.
func (reader *ReaderType) TryReadBits(bits int) uint64 {
	value, err := reader.ReadBits(bits)
	if err != nil {
		panic(err)
	}
	return value
}

// TryReadBytes is a wrapper function that returns the value of bits specified in the parameter.
// Returns uint64. Panics on error.
func (reader *ReaderType) TryReadBytes(bytes int) uint64 {
	value, err := reader.ReadBytes(bytes)
	if err != nil {
		panic(err)
	}
	return value
}

// TryReadString is a wrapper function that returns the string
// that is read until it is null-terminated.
func (reader *ReaderType) TryReadString() string {
	text, _ := reader.ReadString()
	return text
}

// TryReadStringLen is a wrapper function that returns the string
// that is read until the given length is reached or it is null-terminated.
func (reader *ReaderType) TryReadStringLen(length int) string {
	text, _ := reader.ReadStringLen(length)
	return text
}

// TryReadBytesToSlice is a wrapper function that reads the specified amount of bits
// from the parameter and puts each bit into a slice and returns this slice.
func (reader *ReaderType) TryReadBitsToSlice(bits int) []byte {
	bytes := (bits / 8)
	if bits%8 != 0 {
		bytes++
	}
	out := make([]byte, bytes)
	for i := 0; i < bytes; i++ {
		if i == bytes-1 { // Not enough to fill a whole byte
			val, err := reader.ReadBits(bits % 8)
			if err != nil {
				panic(err)
			}
			out[i] = byte(val)
			break
		}
		val, err := reader.ReadBytes(1)
		if err != nil {
			panic(err)
		}
		out[i] = byte(val)
	}
	return out
}

// TryReadBytesToSlice is a wrapper function that reads the specified amount of bytes
// from the parameter and puts each byte into a slice and returns this slice.
func (reader *ReaderType) TryReadBytesToSlice(bytes int) []byte {
	var out []byte
	for i := 0; i < bytes; i++ {
		val, err := reader.ReadBytes(1)
		if err != nil {
			panic(err)
		}
		out = append(out, byte(val))
	}
	return out
}

// SkipBits is a function that increases Reader index
// based on given input bits number.
//
// Returns an error if there are no remaining bits.
func (reader *ReaderType) SkipBits(bits int) error {
	// Read as many raw bytes as we can
	bytes := bits / 8
	buf := make([]byte, bytes)
	_, err := reader.stream.Read(buf)
	if err != nil {
		return err
	}
	// The final read byte should be the new current byte
	if bytes > 0 {
		reader.curByte = buf[bytes-1]
	}
	// Read the extra bits
	for i := bytes * 8; i < bits; i++ {
		_, err := reader.readBit()
		if err != nil {
			return err
		}
	}
	return nil
}

// SkipBytes is a function that increases Reader index
// based on given input bytes number.
//
// Returns an error if there are no remaining bits.
func (reader *ReaderType) SkipBytes(bytes int) error {
	err := reader.SkipBits(bytes * 8)
	if err != nil {
		return err
	}
	return nil
}

// ReadString is a function that reads every byte
// until it is null-terminated (the byte is 0). Returns the
// string that is read until the null-termination.
//
// Returns an error if there are no remaining bits.
func (reader *ReaderType) ReadString() (string, error) {
	var out string
	for {
		value, err := reader.ReadBytes(1)
		if err != nil {
			return out, err
		}
		if value == 0 {
			break
		}
		out += string(rune(value))
	}
	return out, nil
}

// ReadStringLen is a function that reads every byte
// until the given length, or it is null-terminated (the byte is 0).
// Returns the string that is read until the lenth or null-termination.
// It will skip the remaining bytes if it is null-terminated.
//
// Returns an error if there are no remaining bits.
func (reader *ReaderType) ReadStringLen(length int) (string, error) {
	var out string
	for i := 0; i < length; i++ {
		value, err := reader.ReadBytes(1)
		if err != nil {
			return out, err
		}
		if value == 0 {
			reader.SkipBytes(length - 1 - i)
			break
		}
		out += string(rune(value))
	}
	return out, nil
}

// ReadBits is a function that reads the specified amount of bits
// from the parameter and returns the value, error
// based on the output. It can read up to 64 bits. Returns the read
// value in type uint64.
//
// Returns an error if there are no remaining bits.
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

// ReadBytes is a function that reads the specified amount of bytes
// from the parameter and returns the value, error
// based on the output. It can read up to 8 bytes. Returns the read
// value in type uint64.
//
// Returns an error if there are no remaining bits.
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

// ReadBitsToSlice is a function that reads the specified amount of bits
// from the parameter and puts each bit into a slice and returns this slice.
//
// Returns an error if there are no remaining bits.
func (reader *ReaderType) ReadBitsToSlice(bits int) ([]byte, error) {
	bytes := (bits / 8)
	if bits%8 != 0 {
		bytes++
	}
	out := make([]byte, bytes)
	for i := 0; i < bytes; i++ {
		if i == bytes-1 { // Not enough to fill a whole byte
			val, err := reader.ReadBits(bits % 8)
			if err != nil {
				return out, err
			}
			out[i] = byte(val)
			break
		}
		val, err := reader.ReadBytes(1)
		if err != nil {
			return out, err
		}
		out[i] = byte(val)
	}
	return out, nil
}

// ReadBytesToSlice is a function that reads the specified amount of bytes
// from the parameter and puts each byte into a slice and returns this slice.
//
// Returns an error if there are no remaining bytes.
func (reader *ReaderType) ReadBytesToSlice(bytes int) ([]byte, error) {
	var out []byte
	for i := 0; i < bytes; i++ {
		val, err := reader.ReadBytes(1)
		if err != nil {
			return out, err
		}
		out = append(out, byte(val))
	}
	return out, nil
}

// ReadBool is a function that reads one bit and returns the state, error
// based on the output. Returns the read value in a bool format.
//
// Returns an error if there are no remaining bits.
func (reader *ReaderType) ReadBool() (bool, error) {
	val, err := reader.readBit()
	if err != nil {
		return false, err
	}
	return val == 1, nil
}

// readBit is a private function that reads a single bit from the stream.
// This is the main function that makes us read stream data.
func (reader *ReaderType) readBit() (uint8, error) {
	if reader.index == 0 {
		// Read a byte from stream into curByte
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
