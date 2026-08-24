package bitbyte

import (
	"fmt"
)

const BYTES_PER_VALUE = 8
const BITS_PER_BYTE = 8
const BITS_PER_VALUE = BYTES_PER_VALUE * BITS_PER_BYTE
const BR_STR_SIZE = 16 // bitrowで送受信する文字列の一般的なサイズを規定するビット数

func StrBytes(bytes *[]byte) string {
	str := ""
	for i, v := range *bytes {
		if i > 0 {
			str += " "
		}
		str += fmt.Sprintf("%02x", v)
	}
	return str
}

type BitRow struct {
	Index int
	Inc   int
	Body  []uint64
}

func (b *BitRow) GetValues() []int64 {
	int64Body := make([]int64, len(b.Body))
	for i, v := range b.Body {
		int64Body[i] = int64(v)
	}
	return int64Body
}

func (b *BitRow) SetValues(int64Values *[]int64) {
	b.Body = make([]uint64, len(*int64Values))
	for i, v := range *int64Values {
		b.Body[i] = uint64(v)
	}
}

func (b *BitRow) SetIndexHead() {
	b.Index = 0
}

func (b *BitRow) Write(size int, value uint64) error {
	if err := b.Set(b.Index, size, value); err != nil {
		return fmt.Errorf("bitbyte.BitRow.Write:bitbyte.BitRow.Set failed b:%+v, size:%d value:%d\n\t%v", b, size, value, err)
	} else {
		b.Index += size
	}
	return nil
}

func (b *BitRow) WriteString(str string) error {
	runeStr := []rune(str)
	sz := len(runeStr)
	if sz >= 0x10000 { // 2倍して2バイトに収まる必要がある
		return fmt.Errorf("BitRow.WriteString:str too long:%d.", sz)
	} else {
		if err := b.Write(BR_STR_SIZE, uint64(sz)); err != nil {
			return fmt.Errorf("WSData.AppendString:Failed to write string size.\n\t%+v", err)
		} else {
			for _, r := range runeStr {
				if err = b.Write(BR_STR_SIZE, uint64(r)); err != nil {
					return fmt.Errorf("WSData.AppendString:Failed to write rune:%v, str:%s b.Body:%v.\n\t%+v", err, str, b.Body, err)
				}
			}
		}
	}
	return nil
}

func (b *BitRow) Set(shift, size int, value uint64) error {
	if shift < 0 {
		return fmt.Errorf("bitbyte.BitRow.Set shift(%d) is too small", shift)
	} else if size <= 0 {
		return fmt.Errorf("bitbyte.BitRow.Set size(%d) is too small", size)
	} else if size > BITS_PER_VALUE {
		return fmt.Errorf("bitbyte.BitRow.Set size(%d) is too large", size)
	}
	hIndex := (shift + size - 1) / BITS_PER_VALUE
	lIndex := shift / BITS_PER_VALUE
	if hIndex > lIndex { // ビットごとの更新が配列の区切りにまたがってしまう場合
		lShift := shift % BITS_PER_VALUE
		lSize := BITS_PER_VALUE - lShift
		hSize := size - lSize
		lmask := uint64(1<<lSize) - 1
		if err := b.set(hIndex, 0, hSize, value>>lSize); err != nil { // sliceが拡張される場合に備えて、上位インデックスからセットしていく
			return fmt.Errorf("bitbyte.BitRow.Set:bitbyte.BitRow.set failed b:%+v, shift:%d size:%d value:%d\n\t%v", b, shift, size, value, err)
		} else if err = b.set(lIndex, lShift, lSize, value&lmask); err != nil {
			return fmt.Errorf("bitbyte.BitRow.Set:bitbyte.BitRow.set failed b:%+v, shift:%d size:%d value:%d\n\t%v", b, shift, size, value, err)
		}
	} else if err := b.set(hIndex, shift%BITS_PER_VALUE, size, value); err != nil {
		return fmt.Errorf("bitbyte.BitRow.Set:bitbyte.BitRow.set failed b:%+v, shift:%d size:%d value:%d\n\t%v", b, shift, size, value, err)
	}
	return nil
}

func (b *BitRow) set(index, shift, size int, value uint64) error {
	if index < 0 {
		return fmt.Errorf("bitbyte.BitRow.set Index(%d) is too small", index)
	} else if shift < 0 {
		return fmt.Errorf("bitbyte.BitRow.set shift(%d) is too small", shift)
	} else if size <= 0 {
		return fmt.Errorf("bitbyte.BitRow.set size(%d) is too small", size)
	} else if shift+size > BITS_PER_VALUE {
		return fmt.Errorf("bitbyte.BitRow.set shift(%d) size(%d) is too large", shift, size)
	} else if index >= len(b.Body) {
		if b.Inc <= 0 {
			return fmt.Errorf("bitbyte.BitRow.set Index(%d) is out of bounds", index)
		} else {
			add := make([]uint64, b.Inc)
			b.Body = append(b.Body, add...)
		}
	}
	if err := b.setMask(index, shift, size, value); err != nil {
		return fmt.Errorf("bitbyte.BitRow.set failed b:%+v, shift:%d size:%d value:%d\n\t%v", b, shift, size, value, err)
	}
	return nil
}

func (b *BitRow) setMask(index, shift, size int, value uint64) error {
	if mask := uint64(1<<size) - 1; value > mask {
		return fmt.Errorf("bitbyte.BitRow.setMask value(%d) is too large for size(%d)", value, size)
	} else {
		// fmt.Printf("Index:%d shift:%d size:%d value:%d mask:%b ^mask:%b\n", Index, shift, size, value, mask << shift, ^(mask << shift))
		b.Body[index] = (b.Body[index] & (^(mask << shift))) | (value << shift)
	}
	return nil
}

func (b *BitRow) Read(size int) (uint64, error) {
	var val uint64
	var err error
	if val, err = b.Get(b.Index, size); err != nil {
		return 0, fmt.Errorf("bitbyte.BitRow.Read:bitbyte.BitRow.Get failed b:%+v, size:%d\n\t%v", b, size, err)
	} else {
		b.Index += size
	}
	return val, nil
}

func (b *BitRow) Get(shift, size int) (uint64, error) {
	if shift < 0 {
		return 0, fmt.Errorf("bitbyte.BitRow.Get shift(%d) is too small", shift)
	} else if size <= 0 {
		return 0, fmt.Errorf("bitbyte.BitRow.Get size(%d) is too small", size)
	} else if size > BITS_PER_VALUE {
		return 0, fmt.Errorf("bitbyte.BitRow.Get size(%d) is too large", size)
	}
	hIndex := (shift + size - 1) / BITS_PER_VALUE
	lIndex := shift / BITS_PER_VALUE
	if hIndex > lIndex { // ビットごとの更新が配列の区切りにまたがってしまう場合
		lShift := shift % BITS_PER_VALUE
		lSize := BITS_PER_VALUE - lShift
		hSize := size - lSize
		var hVal, lVal uint64
		var err error
		if hVal, err = b.get(hIndex, 0, hSize); err != nil {
			return 0, fmt.Errorf("bitbyte.BitRow.Get:bitbyte.BitRow.get failed b:%+v, shift:%d size:%d\n\t%v", b, shift, size, err)
		} else if lVal, err = b.get(lIndex, lShift, lSize); err != nil {
			return 0, fmt.Errorf("bitbyte.BitRow.Get:bitbyte.BitRow.get failed b:%+v, shift:%d size:%d\n\t%v", b, shift, size, err)
		} else {
			return (hVal << lSize) + lVal, nil
		}
	} else if val, err := b.get(hIndex, shift%BITS_PER_VALUE, size); err != nil {
		return 0, fmt.Errorf("bitbyte.BitRow.Get:bitbyte.BitRow.get failed b:%+v, shift:%d size:%d\n\t%v", b, shift, size, err)
	} else {
		return val, nil
	}
}

func (b *BitRow) get(index, shift, size int) (uint64, error) {
	if index < 0 {
		return 0, fmt.Errorf("bitbyte.BitRow.get index(%d) is too small", index)
	} else if shift < 0 {
		return 0, fmt.Errorf("bitbyte.BitRow.get shift(%d) is too small", shift)
	} else if size <= 0 {
		return 0, fmt.Errorf("bitbyte.BitRow.get size(%d) is too small", size)
	} else if shift+size > BITS_PER_VALUE {
		return 0, fmt.Errorf("bitbyte.BitRow.get shift(%d),size(%d) is too large", shift, size)
	} else if index >= len(b.Body) {
		return 0, fmt.Errorf("bitbyte.BitRow.get index(%d) is out of bounds", index)
	}
	mask := (uint64)(1<<size) - 1
	return (b.Body[index] >> shift) & mask, nil
}

func (b *BitRow) StrBody() string {
	str := ""
	for i, v := range b.Body {
		if i > 0 {
			str += " "
		}
		str += fmt.Sprintf("%016x", v)
	}
	return str
}

func (b *BitRow) Bytes() *[]byte {
	const MASK = (1 << BITS_PER_BYTE) - 1
	sz := (b.Index + BITS_PER_BYTE - 1) / BITS_PER_BYTE
	if sz > len(b.Body)*BYTES_PER_VALUE {
		sz = len(b.Body) * BYTES_PER_VALUE
	}
	bytes := make([]byte, sz)
	for i := 0; i < sz; i++ {
		bytes[i] = byte(b.Body[i>>3]>>((i&7)<<3)) & MASK
	}
	return &bytes
}
