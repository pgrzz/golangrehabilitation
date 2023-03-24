package mq

import (
	"encoding/binary"
	"fmt"
	"io"
)

type writeBuffer struct {
	w io.Writer
	b [16]byte
}

func (wb *writeBuffer) writeInt8(i int8) {
	wb.b[0] = byte(i)
	wb.Write(wb.b[:1])
}

func (wb *writeBuffer) writeInt16(i int16) {
	binary.BigEndian.PutUint16(wb.b[:2], uint16(i))
	wb.Write(wb.b[:2])
}

func (wb *writeBuffer) writeInt32(i int32) {
	binary.BigEndian.PutUint32(wb.b[:4], uint32(i))
	wb.Write(wb.b[:4])
}

func (wb *writeBuffer) writeInt64(i int64) {
	binary.BigEndian.PutUint64(wb.b[:8], uint64(i))
	wb.Write(wb.b[:8])
}

func (wb *writeBuffer) writeVarInt(i int64) {
	u := uint64((i << 1) ^ (i >> 63))
	n := 0

	for u >= 0x80 && n < len(wb.b) {
		wb.b[n] = byte(u) | 0x80
		u >>= 7
		n++
	}

	if n < len(wb.b) {
		wb.b[n] = byte(u)
		n++
	}

	wb.Write(wb.b[:n])
}

func (wb *writeBuffer) writeString(s string) {
	wb.writeInt16(int16(len(s)))
	wb.WriteString(s)
}

func (wb *writeBuffer) writeVarString(s string) {
	wb.writeVarInt(int64(len(s)))
	wb.WriteString(s)
}

func (wb *writeBuffer) writeNullableString(s *string) {
	if s == nil {
		wb.writeInt16(-1)
	} else {
		wb.writeString(*s)
	}
}

func (wb *writeBuffer) writeBytes(b []byte) {
	n := len(b)
	if b == nil {
		n = -1
	}
	wb.writeInt32(int32(n))
	wb.Write(b)
}

func (wb *writeBuffer) writeVarBytes(b []byte) {
	if b != nil {
		wb.writeVarInt(int64(len(b)))
		wb.Write(b)
	} else {
		//-1 is used to indicate nil key
		wb.writeVarInt(-1)
	}
}

func (wb *writeBuffer) writeBool(b bool) {
	v := int8(0)
	if b {
		v = 1
	}
	wb.writeInt8(v)
}

func (wb *writeBuffer) writeArrayLen(n int) {
	wb.writeInt32(int32(n))
}

func (wb *writeBuffer) writeArray(n int, f func(int)) {
	wb.writeArrayLen(n)
	for i := 0; i < n; i++ {
		f(i)
	}
}

func (wb *writeBuffer) writeVarArray(n int, f func(int)) {
	wb.writeVarInt(int64(n))
	for i := 0; i < n; i++ {
		f(i)
	}
}

func (wb *writeBuffer) writeStringArray(a []string) {
	wb.writeArray(len(a), func(i int) { wb.writeString(a[i]) })
}

func (wb *writeBuffer) writeInt32Array(a []int32) {
	wb.writeArray(len(a), func(i int) { wb.writeInt32(a[i]) })
}

func (wb *writeBuffer) write(a interface{}) {
	switch v := a.(type) {
	case int8:
		wb.writeInt8(v)
	case int16:
		wb.writeInt16(v)
	case int32:
		wb.writeInt32(v)
	case int64:
		wb.writeInt64(v)
	case string:
		wb.writeString(v)
	case []byte:
		wb.writeBytes(v)
	case bool:
		wb.writeBool(v)
	case writable:
		v.writeTo(wb)
	default:
		panic(fmt.Sprintf("unsupported type: %T", a))
	}
}

func (wb *writeBuffer) Write(b []byte) (int, error) {
	return wb.w.Write(b)
}

func (wb *writeBuffer) WriteString(s string) (int, error) {
	return io.WriteString(wb.w, s)
}

func (wb *writeBuffer) Flush() error {
	if x, ok := wb.w.(interface{ Flush() error }); ok {
		return x.Flush()
	}
	return nil
}

type writable interface {
	writeTo(*writeBuffer)
}

func varIntLen(i int64) int {
	u := uint64((i << 1) ^ (i >> 63)) // zig-zag encoding
	n := 0

	for u >= 0x80 {
		u >>= 7
		n++
	}

	return n + 1
}

func varBytesLen(b []byte) int {
	return varIntLen(int64(len(b))) + len(b)
}

func varStringLen(s string) int {
	return varIntLen(int64(len(s))) + len(s)
}

func varArrayLen(n int, f func(int) int) int {
	size := varIntLen(int64(n))
	for i := 0; i < n; i++ {
		size += f(i)
	}
	return size
}

func messageSize(key, value []byte) int32 {
	return 4 + // crc
		1 + // magic byte
		1 + // attributes
		8 + // timestamp
		sizeofBytes(key) +
		sizeofBytes(value)
}

func messageSetSize(msgs ...Message) (size int32) {
	for _, msg := range msgs {
		size += 8 + // offset
			4 + // message size
			4 + // crc
			1 + // magic byte
			1 + // attributes
			8 + // timestamp
			sizeofBytes(msg.Key) +
			sizeofBytes(msg.Value)
	}
	return
}
