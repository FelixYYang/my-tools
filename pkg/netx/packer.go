package netx

import (
	"encoding/binary"
	"errors"
	"io"
)

// UnPacker 分包器
// 将流数据按指定规则分包
type UnPacker interface {
	// UnPack reads data from the given io.Reader and unpacks it into the provided byte slice.
	//
	// Parameters:
	// - r: the io.Reader from which to read the data.
	// - buf: the byte slice into which the data will be unpacked.
	//
	// Returns:
	// - int: the number of bytes read and unpacked.
	// - error: an error if the unpacking process encountered any issues.
	UnPack(io.Reader, []byte) (int, error)
}

type Package struct {
	Len  uint64
	Data []byte
}

// NewLenUnPacker creates a new instance of LenUnPacker based on the provided LenOption.
//
// Parameters:
// - option: the LenOption containing the configuration for the LenUnPacker.
//
// Returns:
// - UnPacker: the newly created LenUnPacker instance.
func NewLenUnPacker(option LenOption) UnPacker {
	byteOrder := option.ByteOrder
	var lenByteLen int
	var parseLenFunc func([]byte) int
	switch option.LenType.(type) {
	case int8, uint8:
		parseLenFunc = func(buf []byte) int {
			return int(buf[0])
		}
		lenByteLen = 1
	case int16, uint16:
		parseLenFunc = func(buf []byte) int {
			return int(byteOrder.Uint16(buf))
		}
		lenByteLen = 2
	case int32, uint32:
		parseLenFunc = func(buf []byte) int {
			return int(byteOrder.Uint32(buf))
		}
		lenByteLen = 4
	case int64, uint64:
		parseLenFunc = func(buf []byte) int {
			return int(byteOrder.Uint64(buf))
		}
		lenByteLen = 8
	default:
		panic("LenType is invalid")
	}
	l := &lenPacker{MaxLen: option.MaxLen, Offset: option.Offset, parseLen: parseLenFunc}
	l.OffsetR = l.Offset + lenByteLen
	if l.OffsetR > l.MaxLen {
		panic("Offset+lenByteLen > MaxLen")
	}
	return l
}

type LenOption struct {
	MaxLen    int
	LenType   any
	Offset    int
	ByteOrder binary.ByteOrder
}

type lenPacker struct {
	MaxLen   int
	Offset   int
	OffsetR  int
	parseLen func([]byte) int
}

func (p *lenPacker) UnPack(r io.Reader, buf []byte) (int, error) {
	_, err := io.ReadFull(r, buf[:p.OffsetR])
	if err != nil {
		return 0, err
	}
	totalLen := p.parseLen(buf[p.Offset:p.OffsetR])
	if totalLen == 0 {
		return 0, errors.New("total len parse err")
	}
	if totalLen > p.MaxLen {
		return 0, errors.New("totalLen > MaxLen")
	}
	_, err = io.ReadFull(r, buf[p.OffsetR:totalLen])
	if err != nil {
		return 0, err
	}
	return totalLen, nil
}
