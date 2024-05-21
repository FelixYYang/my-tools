package netx

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"testing"
)

func TestInt32(t *testing.T) {
	var lenSize int64 = 4
	for i := 0; i < 1000; i++ {
		maxLen := rand.Int63n(10000) + lenSize
		option := LenOption{
			MaxLen:    int(maxLen),
			ByteOrder: binary.BigEndian,
		}
		option.Offset = int(rand.Int63n(int64(option.MaxLen) - lenSize + 1))
		var rightLen = rand.Int31n(int32(option.MaxLen - (option.Offset + int(lenSize)) + 1))
		var totalLen = int32(option.Offset) + int32(lenSize) + rightLen
		option.LenType = totalLen
		buffer := bytes.NewBuffer(make([]byte, 0, option.MaxLen))
		if option.Offset > 0 {
			buffer.Write(make([]byte, option.Offset))
		}
		if err := binary.Write(buffer, option.ByteOrder, option.LenType); err != nil {
			t.Error(err)
			return
		}
		if rightLen > 0 {
			buffer.Write(make([]byte, rightLen))
		}
		p := NewLenUnPacker(option)
		packBuf := make([]byte, option.MaxLen)
		n, err := p.UnPack(buffer, packBuf)
		if err != nil {
			t.Error(err)
			return
		}
		pack := packBuf[:n]
		if len(pack) != int(totalLen) {
			t.Errorf("total len not equal %d != %d", len(pack), option.LenType)
			return
		}
	}
}

func TestInt64(t *testing.T) {
	var lenSize int64 = 8
	for i := 0; i < 10000; i++ {
		maxLen := rand.Int63n(10000) + lenSize
		option := LenOption{
			MaxLen:    int(maxLen),
			ByteOrder: binary.BigEndian,
		}
		option.Offset = int(rand.Int63n(int64(option.MaxLen) - lenSize + 1))
		var rightLen = rand.Int31n(int32(option.MaxLen - (option.Offset + int(lenSize)) + 1))
		var totalLen = int64(int32(option.Offset) + int32(lenSize) + rightLen)
		option.LenType = totalLen
		buffer := bytes.NewBuffer(make([]byte, 0, option.MaxLen))
		if option.Offset > 0 {
			buffer.Write(make([]byte, option.Offset))
		}
		if err := binary.Write(buffer, option.ByteOrder, option.LenType); err != nil {
			t.Error(err)
			return
		}
		if rightLen > 0 {
			buffer.Write(make([]byte, rightLen))
		}
		p := NewLenUnPacker(option)
		packBuf := make([]byte, option.MaxLen)
		n, err := p.UnPack(buffer, packBuf)
		if err != nil {
			t.Error(err)
			return
		}
		pack := packBuf[:n]
		if len(pack) != int(totalLen) {
			t.Errorf("total len not equal %d != %d", len(pack), option.LenType)
			return
		}
	}
}
