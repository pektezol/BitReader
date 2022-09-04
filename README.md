# BitReader
A simple bit reader with big/little-endian support for golang.\
Reads data from an existing byte array.\
Uses string manipulation (for now).\
Support reading up to 64 bits at one time.\
Checking for overflowing the data.

## Installation
```bash
$ go get github.com/bisaxa/bitreader
```

## Usage

```go
import "github.com/bisaxa/bitreader"

// data: []byte  Data to read from byte array
// le:     bool  Little-endian(true) or big-endian(false) state
reader := bitreader.Reader(data, le)

// read first bit
state, err := reader.ReadBit()

// skip bits/bytes
err := reader.SkipBits(8)
err := reader.SkipBytes(4)

// read bits
value, err := reader.ReadBits(11)
value, err := reader.ReadBits(64) // up to 64 bits
```

## Error Handling
ReadBits(x), ReadBit(), SkipBits(x) and SkipBytes(x) functions returns an error message when they don't work as expected. It is advised to always handle errors.

## License
This project is licensed under [MIT License](LICENSE).