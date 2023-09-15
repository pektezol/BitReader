// BitReader is a simple bit reader with big/little-endian support for golang.
package bitreader

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
)

// Reader is the main structure of our Reader.
// Whenever index == 0, we need to read a new byte from stream into currentByte
//
// stream io.Reader 	The underlying stream we're reading bytes from
// index uint8			The current index into the byte [0-7]
// currentByte byte		The byte we're currently reading from
// le bool 				Whether to read in little-endian order or not
type Reader struct {
	stream       io.Reader
	index        uint8
	currentByte  byte
	littleEndian bool
}

// NewReader is the main constructor that creates the Reader object
// with stream reader data and little-endian state.
func NewReader(stream io.Reader, littleEndian bool) *Reader {
	return &Reader{
		stream:       stream,
		index:        0,
		currentByte:  0,
		littleEndian: littleEndian,
	}
}

// NewReaderFromBytes is the main constructor that creates the Reader object
// with stream byte data and little-endian state.
func NewReaderFromBytes(stream []byte, littleEndian bool) *Reader {
	return &Reader{
		stream:       bytes.NewReader(stream),
		index:        0,
		currentByte:  0,
		littleEndian: littleEndian,
	}
}

// Fork is a function that copies the original reader into a new reader
// with all of its current values.
func (reader *Reader) Fork() (*Reader, error) {
	originalIndex := reader.index
	originalCurrentByte := reader.currentByte
	byteStream, err := io.ReadAll(reader.stream)
	if err != nil {
		return nil, err // Will only happen when there's no memory, lol
	}
	reader.stream = bytes.NewReader(byteStream)
	return &Reader{
		stream:       bytes.NewReader(byteStream),
		index:        uint8(originalIndex),
		currentByte:  originalCurrentByte,
		littleEndian: reader.littleEndian,
	}, nil
}

// TryReadBool is a wrapper function that gets the state of 1-bit.
//
// Returns true if 1, false if 0. Panics on overflow.
func (reader *Reader) TryReadBool() bool {
	flag, err := reader.ReadBool()
	if err != nil {
		panic(err)
	}
	return flag
}

// TryReadInt1 is a wrapper function that returns the value of 1-bit.
//
// Returns type uint8. Panics on overflow.
func (reader *Reader) TryReadInt1() uint8 {
	value, err := reader.ReadBits(1)
	if err != nil {
		panic(err)
	}
	return uint8(value)
}

// TryReadUInt8 is a wrapper function that returns the value of 8-bits.
//
// Returns uint8. Panics on overflow.
func (reader *Reader) TryReadUInt8() uint8 {
	value, err := reader.ReadBits(8)
	if err != nil {
		panic(err)
	}
	return uint8(value)
}

// TryReadSInt8 is a wrapper function that returns the value of 8-bits.
//
// Returns int8. Panics on overflow.
func (reader *Reader) TryReadSInt8() int8 {
	value, err := reader.ReadBits(8)
	if err != nil {
		panic(err)
	}
	return int8(value)
}

// TryReadUInt16 is a wrapper function that returns the value of 16-bits.
//
// Returns uint16. Panics on overflow.
func (reader *Reader) TryReadUInt16() uint16 {
	value, err := reader.ReadBits(16)
	if err != nil {
		panic(err)
	}
	return uint16(value)
}

// TryReadSInt16 is a wrapper function that returns the value of 16-bits.
//
// Returns uint16. Panics on overflow.
func (reader *Reader) TryReadSInt16() int16 {
	value, err := reader.ReadBits(16)
	if err != nil {
		panic(err)
	}
	return int16(value)
}

// TryReadUInt32 is a wrapper function that returns the value of 32-bits.
//
// Returns uint32. Panics on overflow.
func (reader *Reader) TryReadUInt32() uint32 {
	value, err := reader.ReadBits(32)
	if err != nil {
		panic(err)
	}
	return uint32(value)
}

// TryReadSInt32 is a wrapper function that returns the value of 32-bits.
//
// Returns int32. Panics on overflow.
func (reader *Reader) TryReadSInt32() int32 {
	value, err := reader.ReadBits(32)
	if err != nil {
		panic(err)
	}
	return int32(value)
}

// TryReadUInt64 is a wrapper function that returns the value of 64-bits.
//
// Returns uint64. Panics on overflow.
func (reader *Reader) TryReadUInt64() uint64 {
	value, err := reader.ReadBits(64)
	if err != nil {
		panic(err)
	}
	return value
}

// TryReadSInt64 is a wrapper function that returns the value of 64-bits.
//
// Returns int64. Panics on overflow.
func (reader *Reader) TryReadSInt64() int64 {
	value, err := reader.ReadBits(64)
	if err != nil {
		panic(err)
	}
	return int64(value)
}

// TryReadFloat32 is a wrapper function that returns the value of 32-bits.
//
// Returns float32. Panics on overflow.
func (reader *Reader) TryReadFloat32() float32 {
	value, err := reader.ReadBits(32)
	if err != nil {
		panic(err)
	}
	return math.Float32frombits(uint32(value))
}

// TryReadFloat64 is a wrapper function that returns the value of 64-bits.
//
// Returns float64. Panics on overflow.
func (reader *Reader) TryReadFloat64() float64 {
	value, err := reader.ReadBits(64)
	if err != nil {
		panic(err)
	}
	return math.Float64frombits(value)
}

// TryReadBits is a wrapper function that returns the value of bits specified in the parameter.
//
// Returns uint64. Panics on overflow.
func (reader *Reader) TryReadBits(bits int) uint64 {
	value, err := reader.ReadBits(bits)
	if err != nil {
		panic(err)
	}
	return value
}

// TryReadBytes is a wrapper function that returns the value of bits specified in the parameter.
//
// Returns uint64. Panics on overflow.
func (reader *Reader) TryReadBytes(bytes int) uint64 {
	value, err := reader.ReadBytes(bytes)
	if err != nil {
		panic(err)
	}
	return value
}

// TryReadString is a wrapper function that returns the string
// that is read until it is null-terminated.
//
// Returns string. Panics on overflow.
func (reader *Reader) TryReadString() string {
	text, err := reader.ReadString()
	if err != nil {
		panic(err)
	}
	return text
}

// TryReadStringLength is a wrapper function that returns the string
// that is read until the given length is reached or it is null-terminated.
//
// Returns string. Panics on overflow.
func (reader *Reader) TryReadStringLength(length int) string {
	text, err := reader.ReadStringLength(length)
	if err != nil {
		panic(err)
	}
	return text
}

// TryReadBytesToSlice is a wrapper function that reads the specified amount of bits
// from the parameter and puts each bit into a slice and returns this slice.
//
// Returns []byte. Panics on overflow.
func (reader *Reader) TryReadBitsToSlice(bits int) []byte {
	bytes := (bits / 8)
	if bits%8 != 0 {
		bytes++
	}
	out := make([]byte, bytes)
	for i := 0; i < bytes; i++ {
		if i == bytes-1 { // Not enough to fill a whole byte
			if bits%8 != 0 {
				val, err := reader.ReadBits(bits % 8)
				if err != nil {
					panic(err)
				}
				out[i] = byte(val)
			} else {
				val, err := reader.ReadBytes(1)
				if err != nil {
					panic(err)
				}
				out[i] = byte(val)
			}
			break
		} else {
			val, err := reader.ReadBytes(1)
			if err != nil {
				panic(err)
			}
			out[i] = byte(val)
		}
	}
	return out
}

// TryReadBytesToSlice is a wrapper function that reads the specified amount of bytes
// from the parameter and puts each byte into a slice and returns this slice.
//
// Returns []byte. Panics on overflow.
func (reader *Reader) TryReadBytesToSlice(bytes int) []byte {
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

// TryReadBytesToSlice is a wrapper function that reads the remaining bits
// left in the stream and returns the count of bits.
//
// Returns uint64. Panics on overflow.
func (reader *Reader) TryReadRemainingBits() uint64 {
	bits, err := reader.ReadRemainingBits()
	if err != nil {
		panic(err)
	}
	return bits
}

// ReadBool is a function that reads one bit and returns the state, error
// based on the output. Returns the read value in a bool format.
//
// Returns an error if there are no remaining bits.
func (reader *Reader) ReadBool() (bool, error) {
	val, err := reader.readBit()
	if err != nil {
		return false, err
	}
	return val == 1, nil
}

// ReadBits is a function that reads the specified amount of bits
// from the parameter and returns the value, error
// based on the output. It can read up to 64 bits. Returns the read
// value in type uint64.
//
// Returns an error if there are no remaining bits.
func (reader *Reader) ReadBits(bits int) (uint64, error) {
	if bits < 1 || bits > 64 {
		return 0, errors.New("ReadBits(bits) ERROR: Bits number should be between 1 and 64")
	}
	var val uint64
	for i := 0; i < bits; i++ {
		bit, err := reader.readBit()
		if err != nil {
			return 0, err
		}
		if reader.littleEndian {
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
func (reader *Reader) ReadBytes(bytes int) (uint64, error) {
	if bytes < 1 || bytes > 8 {
		return 0, errors.New("ReadBytes(bytes) ERROR: Bytes number should be between 1 and 8")
	}
	value, err := reader.ReadBits(bytes * 8)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// ReadString is a function that reads every byte
// until it is null-terminated (the byte is 0). Returns the
// string that is read until the null-termination.
//
// Returns an error if there are no remaining bits.
func (reader *Reader) ReadString() (string, error) {
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

// ReadStringLength is a function that reads every byte
// until the given length, or it is null-terminated (the byte is 0).
// Returns the string that is read until the lenth or null-termination.
// It will skip the remaining bytes if it is null-terminated.
//
// Returns an error if there are no remaining bits.
func (reader *Reader) ReadStringLength(length int) (string, error) {
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

// ReadBitsToSlice is a function that reads the specified amount of bits
// from the parameter and puts each bit into a slice and returns this slice.
//
// Returns an error if there are no remaining bits.
func (reader *Reader) ReadBitsToSlice(bits int) ([]byte, error) {
	bytes := (bits / 8)
	if bits%8 != 0 {
		bytes++
	}
	out := make([]byte, bytes)
	for i := 0; i < bytes; i++ {
		if i == bytes-1 { // Not enough to fill a whole byte
			if bits%8 != 0 {
				val, err := reader.ReadBits(bits % 8)
				if err != nil {
					return out, err
				}
				out[i] = byte(val)
			} else {
				val, err := reader.ReadBytes(1)
				if err != nil {
					return out, err
				}
				out[i] = byte(val)
			}
			break
		} else {
			val, err := reader.ReadBytes(1)
			if err != nil {
				return out, err
			}
			out[i] = byte(val)
		}
	}
	return out, nil
}

// ReadBytesToSlice is a function that reads the specified amount of bytes
// from the parameter and puts each byte into a slice and returns this slice.
//
// Returns an error if there are no remaining bytes.
func (reader *Reader) ReadBytesToSlice(bytes int) ([]byte, error) {
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

// SkipBits is a function that increases Reader index
// based on given input bits number.
//
// Returns an error if there are no remaining bits.
func (reader *Reader) SkipBits(bits int) error {
	// Read as many raw bytes as we can
	bytes := bits / 8
	if bytes > 0 {
		buf := make([]byte, bytes)
		_, err := reader.stream.Read(buf)
		if err != nil {
			return err
		}
		// The final read byte should be the new current byte
		reader.currentByte = buf[bytes-1]
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
func (reader *Reader) SkipBytes(bytes int) error {
	err := reader.SkipBits(bytes * 8)
	if err != nil {
		return err
	}
	return nil
}

// ReadRemainingBits is a function that reads the total amount of remaining bits in the stream.
// It first forks the original reader to check this count, so that it does not interfere with the original stream.
//
// Returns an error if there are no remaining bits.
func (reader *Reader) ReadRemainingBits() (uint64, error) {
	newReader, err := reader.Fork()
	if err != nil {
		return 0, err
	}
	var bits uint64 = 0
	for {
		err := newReader.SkipBits(1)
		if err != nil {
			break // EOF
		}
		fmt.Printf("%+v\n", newReader)
		bits++
	}
	return bits, nil
}

// readBit is a private function that reads a single bit from the stream.
// This is the main function that makes us read stream data.
func (reader *Reader) readBit() (uint8, error) {
	if reader.index == 0 {
		// Read a byte from stream into currentByte
		buffer := make([]byte, 1)
		// We are not checking for the n return value from stream.Read, because we are only reading 1 byte at a time.
		// Meaning if an EOF happens with a 1 byte read, we dont have any extra byte reading anyways.
		_, err := reader.stream.Read(buffer)
		if err != nil {
			return 0, err
		}
		reader.currentByte = buffer[0]
	}
	var val bool
	if reader.littleEndian {
		val = (reader.currentByte & (1 << reader.index)) != 0
	} else {
		val = (reader.currentByte & (1 << (7 - reader.index))) != 0
	}
	reader.index = (reader.index + 1) % 8
	if val {
		return 1, nil
	} else {
		return 0, nil
	}
}
