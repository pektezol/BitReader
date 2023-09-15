# BitReader [![Go Reference](https://pkg.go.dev/badge/github.com/pektezol/bitreader.svg)](https://pkg.go.dev/github.com/pektezol/bitreader) [![Go Report Card](https://goreportcard.com/badge/github.com/pektezol/bitreader)](https://goreportcard.com/report/github.com/pektezol/bitreader) [![License: LGPL 2.1](https://img.shields.io/badge/License-LGPL_v2.1-blue.svg)](https://github.com/pektezol/bitreader/blob/main/LICENSE) 
The simplest bit reader with big/little-endian support for Golang.

## Installation
```bash
$ go get github.com/pektezol/bitreader
```

## Usage Examples

```go
import "github.com/pektezol/bitreader"

// ioStream:        io.Reader  Data to read from an io stream
// byteStream:      []byte     Data to read from a byte slice
// littleEndian:    bool       Little-endian(true) or big-endian(false) state
reader := bitreader.NewReader(ioStream, le)
reader := bitreader.NewReaderFromBytes(byteStream, le)

// Fork Reader, Copies Current Reader
newReader, err := reader.Fork()

// Read Total Number of Bits Left
bits, err := reader.ReadRemainingBits()

// Read First Bit
state, err := reader.ReadBool()

// Read Bits/Bytes
value, err := reader.ReadBits(64)       // up to 64 bits
value, err := reader.ReadBytes(8)       // up to 8 bytes

// Read String
text, err := reader.ReadString()            // null-terminated
text, err := reader.ReadStringLength(256)   // length-specified

// Read Bits/Bytes into Slice
arr, err := reader.ReadBitsToSlice(128)
arr, err := reader.ReadBytesToSlice(64)

// Skip Bits/Bytes
err := reader.SkipBits(8)
err := reader.SkipBytes(4)

// Wrapper functions
state := reader.TryReadBool()           // bool
value := reader.TryReadInt1()           // uint8
value := reader.TryReadUInt8()          // uint8
value := reader.TryReadSInt8()          // int8
value := reader.TryReadUInt16()         // uint16
value := reader.TryReadSInt16()         // int16
value := reader.TryReadUInt32()         // uint32
value := reader.TryReadSInt32()         // int32
value := reader.TryReadUInt64()         // uint64
value := reader.TryReadSInt64()         // int64
value := reader.TryReadFloat32()        // float32
value := reader.TryReadFloat64()        // float64
value := reader.TryReadBits(64)         // uint64
value := reader.TryReadBytes(8)         // uint64
text := reader.TryReadString()          // string
text := reader.TryReadStringLength(64)  // string
arr := reader.TryReadBitsToSlice(1024)  // []byte
arr := reader.TryReadBytesToSlice(128)  // []byte
bits := reader.TryReadRemainingBits()   // uint64
```

## Error Handling
All ReadXXX(), SkipXXX() and Fork() functions returns an error message when they don't work as expected. It is advised to always handle errors. \
Wrapper functions, however, only returns the value and panics if an error is encountered, for the sake of ease of use.

## Bug Report / Feature Request
Using [Github Issues](https://github.com/pektezol/BitReader/issues/new/choose), you can report a bug that you encountered and/or request a feature that you would like to be added.

## Documentation

Full documentation can be found in https://pkg.go.dev/github.com/pektezol/bitreader

## License
This project is licensed under [GNU Lesser General Public License version 2.1](https://www.gnu.org/licenses/old-licenses/lgpl-2.1).
