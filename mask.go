package gosocket

import "encoding/binary"

//	对输入的字节数组payload进行掩码处理
//
// @0xAAC 亮点
func unMaskPayload(payload []byte, maskingKey []byte) {
	// 将4字节的掩码键转换为32位无符号整数mk
	var mk = binary.LittleEndian.Uint32(maskingKey)
	// 将mk左移32位，并加上原值，得到64位的ksy64，其前后32位都是相同的掩码键
	// 方便后续对8字节（64位）数据块的处理
	var key64 = uint64(mk)<<32 + uint64(mk)

	// 批量处理64字节块
	for len(payload) >= 64 {
		v := binary.LittleEndian.Uint64(payload) // payload[0:8]
		binary.LittleEndian.PutUint64(payload, v^key64)

		v = binary.LittleEndian.Uint64(payload[8:16])
		binary.LittleEndian.PutUint64(payload[8:16], v^key64)

		v = binary.LittleEndian.Uint64(payload[16:24])
		binary.LittleEndian.PutUint64(payload[16:24], v^key64)

		v = binary.LittleEndian.Uint64(payload[24:32])
		binary.LittleEndian.PutUint64(payload[24:32], v^key64)

		v = binary.LittleEndian.Uint64(payload[32:40])
		binary.LittleEndian.PutUint64(payload[32:40], v^key64)

		v = binary.LittleEndian.Uint64(payload[40:48])
		binary.LittleEndian.PutUint64(payload[40:48], v^key64)

		v = binary.LittleEndian.Uint64(payload[48:56])
		binary.LittleEndian.PutUint64(payload[48:56], v^key64)

		v = binary.LittleEndian.Uint64(payload[56:64])
		binary.LittleEndian.PutUint64(payload[56:64], v^key64)

		// 处理完64字节后，继续处理下一个64字节块
		payload = payload[64:]
	}

	// 剩余字节长度小于64，但大于等于8，批量处理剩余的8字节块
	for len(payload) >= 8 {
		v := binary.LittleEndian.Uint64(payload[:8])
		binary.LittleEndian.PutUint64(payload[:8], v^key64)
		payload = payload[8:]
	}

	var n = len(payload)
	for i := 0; i < n; i++ {
		// 等价于 i % 4， 0000 & 0011，0001 & 0011
		idx := i & 3
		payload[i] ^= maskingKey[idx]
	}
}
