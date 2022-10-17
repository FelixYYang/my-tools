package netx

import (
	"encoding/binary"
	"errors"
	"io"
)

// Packer 分包器
// 将流数据按指定规则分包
type Packer interface {
	Pack(r io.Reader) (Package, error)
}

type Package struct {
	Len  uint64
	Data []byte
}

func NewLenPacker(option LenOption) Packer {
	l := &lenPacker{MaxLen: option.MaxLen, LenType: option.LenType, Offset: option.Offset, ByteOrder: option.ByteOrder}
	var lenByteLen uint64
	switch option.LenType.(type) {
	case int8, uint8:
		lenByteLen = 1
	case int16, uint16:
		lenByteLen = 2
	case int32, uint32:
		lenByteLen = 4
	case int64, uint64:
		lenByteLen = 8
	default:
		panic("LenType is invalid")
	}
	l.OffsetR = l.Offset + lenByteLen
	if l.OffsetR > l.MaxLen {
		panic("Offset+lenByteLen > MaxLen")
	}
	l.leftData = make([]byte, l.MaxLen)
	return l
}

type LenOption struct {
	MaxLen    uint64
	LenType   any
	Offset    uint64
	ByteOrder binary.ByteOrder
}

type lenPacker struct {
	MaxLen    uint64
	LenType   any
	ByteOrder binary.ByteOrder
	Offset    uint64
	OffsetR   uint64
	leftData  []byte
}

func (p *lenPacker) Pack(r io.Reader) (Package, error) {
	_, err := io.ReadFull(r, p.leftData[:p.OffsetR])
	if err != nil {
		return Package{}, err
	}
	totalLen := p.parseLen(p.leftData[p.Offset:p.OffsetR])
	if totalLen == 0 {
		return Package{}, errors.New("total len parse err")
	}
	if totalLen > p.MaxLen {
		return Package{}, errors.New("totalLen > MaxLen")
	}
	_, err = io.ReadFull(r, p.leftData[p.OffsetR:totalLen])
	if err != nil {
		return Package{}, err
	}

	return Package{
		Len:  totalLen,
		Data: p.leftData[:totalLen],
	}, nil
}

func (p *lenPacker) parseLen(buf []byte) (lenNum uint64) {
	switch p.LenType.(type) {
	case int8, uint8:
		lenNum = uint64(buf[0])
	case int16, uint16:
		lenNum = uint64(p.ByteOrder.Uint16(buf))
	case int32, uint32:
		lenNum = uint64(p.ByteOrder.Uint32(buf))
	case int64, uint64:
		lenNum = p.ByteOrder.Uint64(buf)
	default:
		panic("LenType is invalid")
	}
	return
}
