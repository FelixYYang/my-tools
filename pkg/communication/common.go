package communication

import "encoding/binary"

const (
	Connect    = 0x00000001
	Terminate  = 0x00000002 // 终止连接
	ActiveTest = 0x00000008 // 心跳
)

//字节转换成无符号整形
func bytesToUint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

//整形转换成字节
func uint32ToBytes(n uint32) []byte {
	x := n
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, x)
	return b
}

// IsRespCid 是否是回执命令
func IsRespCid(rid uint32) bool {
	return rid&0x80000000 == 0x80000000
}

// 转换回执命令为请求命令
func respCidToReqCid(rid uint32) uint32 {
	return rid - 0x80000000
}
