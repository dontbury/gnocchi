package bitbyte

import (
	"fmt"
)

type Bytes *[]byte
type BytesList struct {
	bytes []Bytes
	next  *BytesList
}

func CreateBytesList(size int) *BytesList {
	b := BytesList{
		bytes: make([]Bytes, size),
		next:  nil,
	}
	return &b
}

func (b *BytesList) Size() int {
	return len(b.bytes)
}

func (b *BytesList) ByteSize() int {
	size := 0
	for _, v := range b.bytes {
		size += len(*v)
	}
	return size
}

func (b *BytesList) SetBytes(index int, bytes *[]byte) error {
	if index < len(b.bytes) {
		b.bytes[index] = Bytes(bytes)
	} else {
		return fmt.Errorf("bitbyte.BytesList.SetBytes:index too large index:%d len():%d", index, len(b.bytes))
	}
	return nil
}

func (b *BytesList) Bytes(maxDepth int) (*[]byte, error) {
	var bs Bytes
	t := b
	depth := 0
	sz := 0
	for t != nil {
		depth++
		for _, bs = range t.bytes {
			sz += len(*bs)
		}

		t = t.next
		if depth > maxDepth {
			return nil, fmt.Errorf("bitbyte.BytesList.Bytes:Loop depth over %d", maxDepth)
		}
	}

	index := 0
	buf := make([]byte, sz)

	t = b
	depth = 1
	for t != nil {
		for _, bs = range t.bytes {
			for _, bt := range *bs {
				buf[index] = bt
				index++
			}
		}

		t = t.next
		depth++
		if depth > maxDepth {
			return nil, fmt.Errorf("bitbyte.BytesList.Bytes:Loop depth over %d", maxDepth)
		}
	}

	return &buf, nil
}

func (b *BytesList) AppendTail(bl *BytesList, maxDepth int) error {
	t := b
	depth := 0
	for t != nil {
		if t.next == nil {
			t.next = bl
			break
		}
		t = t.next
		depth++
		if depth > maxDepth {
			return fmt.Errorf("bitbyte.BytesList.AppendTail:Loop depth over %d", maxDepth)
		}
	}
	return nil
}
