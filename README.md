# BitReader [![Go Report Card](https://goreportcard.com/badge/github.com/pektezol/bitreader)](https://goreportcard.com/report/github.com/pektezol/bitreader) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/pektezol/bitreader/blob/main/LICENSE) [![Go Reference](https://pkg.go.dev/badge/github.com/pektezol/bitreader.svg)](https://pkg.go.dev/github.com/pektezol/bitreader)
A simple bit reader with big/little-endian support for golang.\
Reads data from an existing byte array.\
Uses string manipulation (for now).\
Support reading up to 64 bits at one time.\
Checking for overflowing the data.

## Installation
```bash
$ go get github.com/pektezol/bitreader
```

## Usage

```go
import "github.com/pektezol/bitreader"

// data:  io.Reader  Data to read from a io stream
// le:    bool       Little-endian(true) or big-endian(false) state
reader := bitreader.Reader(data, le)

// Read First Bit
state, err := reader.ReadBool()

// Skip Bits/Bytes
err := reader.SkipBits(8)
err := reader.SkipBytes(4)

// Read Bits/Bytes
value, err := reader.ReadBytes(4)
value, err := reader.ReadBits(64) // up to 64 bits

// Wrapper functions
state := reader.TryReadBool()       // bool
value := reader.TryReadInt1()       // uint8
value := reader.TryReadInt8()       // uint8
value := reader.TryReadInt16()      // uint16
value := reader.TryReadInt32()      // uint32
value := reader.TryReadInt64()      // uint64
value := reader.TryReadFloat32()    // float32
value := reader.TryReadFloat64()    // float64
```

## Error Handling
ReadBits(x), ReadBytes(x), ReadBool(), SkipBits(x) and SkipBytes(x) functions returns an error message when they don't work as expected. It is advised to always handle errors. \
Wrapper functions, however, only returns the value and panics if an error is encountered.

## Bug Report / Feature Request
Using [Github Issues](https://github.com/pektezol/BitReader/issues/new/choose), you can report a bug that you encountered and/or request a feature that you would like to be added.

## Documentation

Full documentation can be found in https://pkg.go.dev/github.com/pektezol/bitreader

## License
This project is licensed under [MIT License](LICENSE).
