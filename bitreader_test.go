// BitReader is a simple bit reader with big/little-endian support for golang.
// %83.2 coerage
package bitreader

import (
	"bytes"
	"io"
	"math"
	"reflect"
	"testing"
)

func TestNewReader(t *testing.T) {
	stream := bytes.NewReader([]byte{0x01, 0x02, 0x03})
	type args struct {
		stream       io.Reader
		littleEndian bool
	}
	tests := []struct {
		name string
		args args
		want *Reader
	}{
		{
			name: "ReaderLE",
			args: args{
				stream:       stream,
				littleEndian: true,
			},
			want: &Reader{
				stream:       stream,
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
		},
		{
			name: "ReaderBE",
			args: args{
				stream:       stream,
				littleEndian: false,
			},
			want: &Reader{
				stream:       stream,
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReader(tt.args.stream, tt.args.littleEndian); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReader() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestNewReaderFromBytes(t *testing.T) {
	type args struct {
		stream       []byte
		littleEndian bool
	}
	tests := []struct {
		name string
		args args
		want *Reader
	}{
		{
			name: "ReaderLE",
			args: args{
				stream:       []byte{0x01, 0x02, 0x03},
				littleEndian: true,
			},
			want: &Reader{
				stream:       bytes.NewReader([]byte{0x01, 0x02, 0x03}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
		},
		{
			name: "ReaderBE",
			args: args{
				stream:       []byte{0x01, 0x02, 0x03},
				littleEndian: false,
			},
			want: &Reader{
				stream:       bytes.NewReader([]byte{0x01, 0x02, 0x03}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReaderFromBytes(tt.args.stream, tt.args.littleEndian); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReaderFromBytes() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_Fork(t *testing.T) {
	stream := bytes.NewReader([]byte{53})
	tests := []struct {
		name    string
		reader  *Reader
		want    *Reader
		wantErr bool
	}{
		{
			name: "Fork",
			reader: &Reader{
				stream:       stream,
				index:        4,
				currentByte:  53,
				littleEndian: false,
			},
			want: &Reader{
				index:        4,
				currentByte:  53,
				littleEndian: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.reader.Fork()
			tt.want.stream = got.stream
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.Fork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader.Fork() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadBool(t *testing.T) {
	tests := []struct {
		name   string
		reader *Reader
		want   bool
	}{
		{
			name: "ReadBoolTrueLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b00000001}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			want: true,
		},
		{
			name: "ReadBoolTrueBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10000000}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want: true,
		},
		{
			name: "ReadBoolFalseLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b00000010}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			want: false,
		},
		{
			name: "ReadBoolFalseBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b01000000}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.TryReadBool(); got != tt.want {
				t.Errorf("Reader.TryReadBool() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadInt1(t *testing.T) {
	tests := []struct {
		name   string
		reader *Reader
		want   uint8
	}{
		{
			name: "ReadInt1TrueLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b00000001}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			want: 0b1,
		},
		{
			name: "ReadInt1TrueBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10000000}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want: 0b1,
		},
		{
			name: "ReadInt1FalseLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b00000010}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			want: 0b0,
		},
		{
			name: "ReadInt1FalseBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b01000000}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want: 0b0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.TryReadInt1(); got != tt.want {
				t.Errorf("Reader.TryReadInt1() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadUInt8(t *testing.T) {
	tests := []struct {
		name   string
		reader *Reader
		want   uint8
	}{
		{
			name: "ReadUInt8LE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{202}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			want: 202,
		},
		{
			name: "ReadUInt8BE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{202}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want: 202,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.TryReadUInt8(); got != tt.want {
				t.Errorf("Reader.TryReadUInt8() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadSInt8(t *testing.T) {
	tests := []struct {
		name   string
		reader *Reader
		want   int8
	}{
		{
			name: "ReadSInt8LE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{202}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			want: -54,
		},
		{
			name: "ReadSInt8BE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{202}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want: -54,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.TryReadSInt8(); got != tt.want {
				t.Errorf("Reader.TryReadSInt8() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadUInt16(t *testing.T) {
	tests := []struct {
		name   string
		reader *Reader
		want   uint16
	}{
		{
			name: "ReadUInt16LE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10101010, 0b01010101}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			want: 0b0101010110101010,
		},
		{
			name: "ReadUInt16BE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10101010, 0b01010101}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want: 0b1010101001010101,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.TryReadUInt16(); got != tt.want {
				t.Errorf("Reader.TryReadUInt16() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadSInt16(t *testing.T) {
	tests := []struct {
		name   string
		reader *Reader
		want   int16
	}{
		{
			name: "ReadSInt16LE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10101010, 0b01010101}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			want: 21930,
		},
		{
			name: "ReadSInt16BE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10101010, 0b01010101}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want: -21931,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.TryReadSInt16(); got != tt.want {
				t.Errorf("Reader.TryReadSInt16() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadUInt32(t *testing.T) {
	tests := []struct {
		name   string
		reader *Reader
		want   uint32
	}{
		{
			name: "ReadUInt32LE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10101010, 0b01010101, 0b11110000, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			want: 0b00001111111100000101010110101010,
		},
		{
			name: "ReadUInt32BE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10101010, 0b01010101, 0b11110000, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want: 0b10101010010101011111000000001111,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.TryReadUInt32(); got != tt.want {
				t.Errorf("Reader.TryReadUInt32() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadSInt32(t *testing.T) {
	tests := []struct {
		name   string
		reader *Reader
		want   int32
	}{
		{
			name: "ReadSInt32LE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10101010, 0b01010101, 0b11110000, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			want: 267408810,
		},
		{
			name: "ReadSInt32BE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10101010, 0b01010101, 0b11110000, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want: -1437208561,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.TryReadSInt32(); got != tt.want {
				t.Errorf("Reader.TryReadSInt32() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadUInt64(t *testing.T) {
	tests := []struct {
		name   string
		reader *Reader
		want   uint64
	}{
		{
			name: "ReadUInt64LE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10101010, 0b01010101, 0b11110000, 0b00001111, 0b10101010, 0b01010101, 0b11110000, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			want: 0b0000111111110000010101011010101000001111111100000101010110101010,
		},
		{
			name: "ReadUInt64BE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10101010, 0b01010101, 0b11110000, 0b00001111, 0b10101010, 0b01010101, 0b11110000, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want: 0b1010101001010101111100000000111110101010010101011111000000001111,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.TryReadUInt64(); got != tt.want {
				t.Errorf("Reader.TryReadUInt64() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadSInt64(t *testing.T) {
	tests := []struct {
		name   string
		reader *Reader
		want   int64
	}{
		{
			name: "ReadSInt64LE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10101010, 0b01010101, 0b11110000, 0b00001111, 0b10101010, 0b01010101, 0b11110000, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			want: 1148512093879686570,
		},
		{
			name: "ReadSInt64BE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10101010, 0b01010101, 0b11110000, 0b00001111, 0b10101010, 0b01010101, 0b11110000, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want: -6172763764168462321,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.TryReadSInt64(); got != tt.want {
				t.Errorf("Reader.TryReadSInt64() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadFloat32(t *testing.T) {
	tests := []struct {
		name   string
		reader *Reader
		want   float32
	}{
		{
			name: "ReadFloat32LE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10101010, 0b01010101, 0b11110000, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			want: math.Float32frombits(0b00001111111100000101010110101010),
		},
		{
			name: "ReadFloat32BE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10101010, 0b01010101, 0b11110000, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want: math.Float32frombits(0b10101010010101011111000000001111),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.TryReadFloat32(); got != tt.want {
				t.Errorf("Reader.TryReadFloat32() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadFloat64(t *testing.T) {
	tests := []struct {
		name   string
		reader *Reader
		want   float64
	}{
		{
			name: "ReadFloat64LE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10101010, 0b01010101, 0b11110000, 0b00001111, 0b10101010, 0b01010101, 0b11110000, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			want: math.Float64frombits(0b0000111111110000010101011010101000001111111100000101010110101010),
		},
		{
			name: "ReadFloat64BE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10101010, 0b01010101, 0b11110000, 0b00001111, 0b10101010, 0b01010101, 0b11110000, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want: math.Float64frombits(0b1010101001010101111100000000111110101010010101011111000000001111),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.TryReadFloat64(); got != tt.want {
				t.Errorf("Reader.TryReadFloat64() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadBits(t *testing.T) {
	type args struct {
		bits int
	}
	tests := []struct {
		name   string
		reader *Reader
		args   args
		want   uint64
	}{
		{
			name: "ReadBitsLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110000, 0b01010101}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			args: args{
				bits: 12,
			},
			want: 0b010111110000,
		},
		{
			name: "ReadBitsBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110000, 0b01010101}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			args: args{
				bits: 12,
			},
			want: 0b111100000101,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.TryReadBits(tt.args.bits); got != tt.want {
				t.Errorf("Reader.TryReadBits() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadBytes(t *testing.T) {
	type args struct {
		bytes int
	}
	tests := []struct {
		name   string
		reader *Reader
		args   args
		want   uint64
	}{
		{
			name: "ReadBytesLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110000, 0b01010101}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			args: args{
				bytes: 2,
			},
			want: 0b0101010111110000,
		},
		{
			name: "ReadBytesBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110000, 0b01010101}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			args: args{
				bytes: 2,
			},
			want: 0b1111000001010101,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.TryReadBytes(tt.args.bytes); got != tt.want {
				t.Errorf("Reader.TryReadBytes() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadString(t *testing.T) {
	tests := []struct {
		name   string
		reader *Reader
		want   string
	}{
		{
			name: "ReadStringLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{'H', 'e', 'l', 'l', 'o', 0, '!'}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			want: "Hello",
		},
		{
			name: "ReadStringBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{'W', 'o', 'r', 'l', 'd', 0, '!'}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want: "World",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.TryReadString(); got != tt.want {
				t.Errorf("Reader.TryReadString() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadStringLength(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name   string
		reader *Reader
		args   args
		want   string
	}{
		{
			name: "ReadStringLengthLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{'H', 'e', 'l', 'l', 'o', 0, '!'}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			args: args{
				length: 4,
			},
			want: "Hell",
		},
		{
			name: "ReadStringLengthBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{'W', 'o', 'r', 'l', 'd', '!', '?'}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			args: args{
				length: 6,
			},
			want: "World!",
		},
		{
			name: "ReadStringLengthNullHitBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{'W', 'o', 'r', 'l', 'd', 0, '!', '?'}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			args: args{
				length: 7,
			},
			want: "World",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.TryReadStringLength(tt.args.length); got != tt.want {
				t.Errorf("Reader.TryReadStringLength() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadBitsToSlice(t *testing.T) {
	type args struct {
		bits int
	}
	tests := []struct {
		name   string
		reader *Reader
		args   args
		want   []byte
	}{
		{
			name: "ReadBitsToSliceBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110010, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			args: args{
				bits: 12,
			},
			want: []byte{0b11110010, 0b0},
		},
		{
			name: "ReadBitsToSliceLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110010, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			args: args{
				bits: 12,
			},
			want: []byte{0b11110010, 0b00001111},
		},
		{
			name: "ReadBitsToSliceBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110010, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			args: args{
				bits: 16,
			},
			want: []byte{0b11110010, 0b00001111},
		},
		{
			name: "ReadBitsToSliceLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110010, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			args: args{
				bits: 16,
			},
			want: []byte{0b11110010, 0b00001111},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.TryReadBitsToSlice(tt.args.bits); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader.TryReadBitsToSlice() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadBytesToSlice(t *testing.T) {
	type args struct {
		bytes int
	}
	tests := []struct {
		name   string
		reader *Reader
		args   args
		want   []byte
	}{
		{
			name: "ReadBytesToSliceBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110010, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			args: args{
				bytes: 2,
			},
			want: []byte{0b11110010, 0b00001111},
		},
		{
			name: "ReadBytesToSliceLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110010, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			args: args{
				bytes: 2,
			},
			want: []byte{0b11110010, 0b00001111},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reader.TryReadBytesToSlice(tt.args.bytes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader.TryReadBytesToSlice() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_TryReadRemainingBits(t *testing.T) {
	tests := []struct {
		name    string
		reader  *Reader
		want    uint64
		wantErr bool
	}{
		{
			name: "ReadRemainingBits",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0x11, 0x22}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want:    16,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.reader.TryReadRemainingBits()
			if got != tt.want {
				t.Errorf("Reader.TryReadRemainingBits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_ReadBool(t *testing.T) {
	tests := []struct {
		name    string
		reader  *Reader
		want    bool
		wantErr bool
	}{
		{
			name: "ReadBoolTrueLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b00000001}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			want: true,
		},
		{
			name: "ReadBoolTrueBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b10000000}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want: true,
		},
		{
			name: "ReadBoolFalseLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b00000010}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			want: false,
		},
		{
			name: "ReadBoolFalseBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b01000000}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.reader.ReadBool()
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.ReadBool() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Reader.ReadBool() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_ReadBits(t *testing.T) {
	type args struct {
		bits int
	}
	tests := []struct {
		name    string
		reader  *Reader
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name: "ReadBitsLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110000, 0b01010101}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			args: args{
				bits: 12,
			},
			want: 0b010111110000,
		},
		{
			name: "ReadBitsBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110000, 0b01010101}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			args: args{
				bits: 12,
			},
			want: 0b111100000101,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.reader.ReadBits(tt.args.bits)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.ReadBits() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Reader.ReadBits() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_ReadBytes(t *testing.T) {
	type args struct {
		bytes int
	}
	tests := []struct {
		name    string
		reader  *Reader
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name: "ReadBytesLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110000, 0b01010101}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			args: args{
				bytes: 2,
			},
			want: 0b0101010111110000,
		},
		{
			name: "ReadBytesBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110000, 0b01010101}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			args: args{
				bytes: 2,
			},
			want: 0b1111000001010101,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.reader.ReadBytes(tt.args.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.ReadBytes() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Reader.ReadBytes() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_ReadString(t *testing.T) {
	tests := []struct {
		name    string
		reader  *Reader
		want    string
		wantErr bool
	}{
		{
			name: "ReadStringLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{'H', 'e', 'l', 'l', 'o', 0, '!'}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			want: "Hello",
		},
		{
			name: "ReadStringBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{'W', 'o', 'r', 'l', 'd', 0, '!'}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want: "World",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.reader.ReadString()
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.ReadString() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Reader.ReadString() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_ReadStringLength(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name    string
		reader  *Reader
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ReadStringLengthLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{'H', 'e', 'l', 'l', 'o', 0, '!'}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			args: args{
				length: 4,
			},
			want: "Hell",
		},
		{
			name: "ReadStringLengthBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{'W', 'o', 'r', 'l', 'd', '!', '?'}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			args: args{
				length: 6,
			},
			want: "World!",
		},
		{
			name: "ReadStringLengthNullHitBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{'W', 'o', 'r', 'l', 'd', 0, '!', '?'}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			args: args{
				length: 7,
			},
			want: "World",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.reader.ReadStringLength(tt.args.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.ReadStringLength() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Reader.ReadStringLength() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_ReadBitsToSlice(t *testing.T) {
	type args struct {
		bits int
	}
	tests := []struct {
		name    string
		reader  *Reader
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "ReadBitsToSliceBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110010, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			args: args{
				bits: 12,
			},
			want: []byte{0b11110010, 0b0},
		},
		{
			name: "ReadBitsToSliceLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110010, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			args: args{
				bits: 12,
			},
			want: []byte{0b11110010, 0b00001111},
		},
		{
			name: "ReadBitsToSliceBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110010, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			args: args{
				bits: 16,
			},
			want: []byte{0b11110010, 0b00001111},
		},
		{
			name: "ReadBitsToSliceLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110010, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			args: args{
				bits: 16,
			},
			want: []byte{0b11110010, 0b00001111},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.reader.ReadBitsToSlice(tt.args.bits)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.ReadBitsToSlice() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader.ReadBitsToSlice() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_ReadBytesToSlice(t *testing.T) {
	type args struct {
		bytes int
	}
	tests := []struct {
		name    string
		reader  *Reader
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "ReadBytesToSliceBE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110010, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			args: args{
				bytes: 2,
			},
			want: []byte{0b11110010, 0b00001111},
		},
		{
			name: "ReadBytesToSliceLE",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0b11110010, 0b00001111}),
				index:        0,
				currentByte:  0,
				littleEndian: true,
			},
			args: args{
				bytes: 2,
			},
			want: []byte{0b11110010, 0b00001111},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.reader.ReadBytesToSlice(tt.args.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.ReadBytesToSlice() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader.ReadBytesToSlice() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestReader_ReadRemainingBits(t *testing.T) {
	tests := []struct {
		name    string
		reader  *Reader
		want    uint64
		wantErr bool
	}{
		{
			name: "ReadRemainingBits",
			reader: &Reader{
				stream:       bytes.NewReader([]byte{0x11, 0x22}),
				index:        0,
				currentByte:  0,
				littleEndian: false,
			},
			want:    16,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.reader.ReadRemainingBits()
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.ReadRemainingBits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Reader.ReadRemainingBits() = %v, want %v", got, tt.want)
			}
		})
	}
}
