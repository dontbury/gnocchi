package wskt

import (
	"fmt"

	"github.com/dontbury/gnocchi/bitbyte"
)

const (
	BYTE1_SIZE = 1
	BYTE2_SIZE = 2
	BYTE3_SIZE = 3
	BYTE4_SIZE = 4
	BYTE8_SIZE = 8

	WS_WORD_SIZE    = 16 // WebSocketで送受信する文字のビット数
	STR_SIZE_LENGTH = 16 // 文字列の長さを収めるためのビット数
	CHAR_SIZE       = 16 // 1文字あたりのビット数

	WSKTBUF_INCREMENT = 1024
)

type WSBuf struct {
	Br bitbyte.BitRow
}

func (s *WSBuf) Create(index, inc, size int) {
	s.Br = bitbyte.BitRow{Index: index, Inc: inc, Body: make([]uint64, (size+bitbyte.BITS_PER_VALUE-1)/bitbyte.BITS_PER_VALUE)}
}

func (s *WSBuf) CreateIdx(index, size int) {
	s.Create(index, WSKTBUF_INCREMENT, size)
}

func (s *WSBuf) CreateSize(size int) {
	s.Create(0, WSKTBUF_INCREMENT, size)
}

func (s *WSBuf) CreateWSBBytes(index, inc int, buf *[]byte) error {
	sz := (len(*buf) + bitbyte.BYTES_PER_VALUE - 1) / bitbyte.BYTES_PER_VALUE
	s.Br = bitbyte.BitRow{Index: index, Inc: inc, Body: make([]uint64, sz)}
	var err error
	for i, v := range *buf { // Indexは先頭に固定したまま、bufの内容をコピーする
		if err = s.Br.Set(i*bitbyte.BITS_PER_BYTE, bitbyte.BITS_PER_BYTE, uint64(v)); err != nil {
			return fmt.Errorf("wskt.WSBuf.CreateWSBBytes:bitbyte.BitRow.Append1Byte faile.\n\t%v", err)
		}
	}
	return nil
}

func (s *WSBuf) CreateWSB(index, headersize int, wsb *WSBuf) error {
	sz := (headersize + (len(wsb.Br.Body)+1)*bitbyte.BITS_PER_VALUE - 1 - wsb.Br.Index) / bitbyte.BITS_PER_VALUE
	s.Br = bitbyte.BitRow{Index: index, Inc: 0, Body: make([]uint64, sz)}
	total := len(wsb.Br.Body) * bitbyte.BITS_PER_VALUE
	var size int
	var value uint64
	var err error
	for idx := wsb.Br.Index; idx < total; {
		if total-idx < bitbyte.BITS_PER_VALUE {
			size = total - idx
		} else {
			size = bitbyte.BITS_PER_VALUE
		}
		if value, err = wsb.Br.Get(idx, size); err != nil {
			return fmt.Errorf("WSBuf.CreateWSB:Failed to get bit value.\n\t%+v", err)
		}
		if err = s.Br.Set(headersize+idx-wsb.Br.Index, size, value); err != nil {
			return fmt.Errorf("WSBuf.CreateWSB:Failed to set bit value.\n\t%+v", err)
		}
		idx += size
	}
	return nil
}

func (s *WSBuf) GetIndex() int {
	return s.Br.Index
}

func (s *WSBuf) SetIndexHead() {
	s.Br.Index = 0
}

func (s *WSBuf) SetIndexTail() {
	s.Br.Index = len(s.Br.Body)
}

// バッファをそのまま取得
func (s *WSBuf) Bytes() *[]byte {
	buf := make([]byte, len(s.Br.Body)*bitbyte.BYTES_PER_VALUE)
	for i, v := range s.Br.Body {
		for j := 0; j < bitbyte.BYTES_PER_VALUE; j++ {
			buf[i*bitbyte.BYTES_PER_VALUE+j] = byte((v >> (bitbyte.BYTES_PER_VALUE * j)) & 0xFF)
		}
	}
	return &buf
}

// 送信用に後ろの未使用部分を切り詰めたバッファを取得
func (s *WSBuf) GetSendBuf() (*[]byte, error) {
	buf := []byte{}
	if s.Br.Index < len(s.Br.Body)*bitbyte.BITS_PER_VALUE {
		sz := s.Br.Index
		buf = make([]byte, (sz+bitbyte.BITS_PER_BYTE-1)/bitbyte.BITS_PER_BYTE)
		var val uint64
		var err error
		for i := 0; len(buf) > i; i++ {
			if sz >= bitbyte.BITS_PER_BYTE {
				if val, err = s.Br.Get(i*bitbyte.BITS_PER_BYTE, bitbyte.BITS_PER_BYTE); err != nil {
					return nil, fmt.Errorf("WSBuf.GetSendBuf:Failed to get bit value.\n\t%+v", err)
				}
				buf[i] = byte(val)
				sz -= bitbyte.BITS_PER_BYTE
			} else {
				if val, err = s.Br.Get(i*bitbyte.BITS_PER_BYTE, sz); err != nil {
					return nil, fmt.Errorf("WSBuf.GetSendBuf:Failed to get bit value.\n\t%+v", err)
				}
				buf[i] = byte(val)
				sz = 0
			}
		}
	} else {
		buf = make([]byte, len(s.Br.Body)*bitbyte.BITS_PER_VALUE/bitbyte.BITS_PER_BYTE)
		for i, v := range s.Br.Body {
			for j := 0; j < bitbyte.BYTES_PER_VALUE; j++ {
				buf[i*bitbyte.BYTES_PER_VALUE+j] = byte((v >> (bitbyte.BYTES_PER_VALUE * j)) & 0xFF)
			}
		}
	}
	return &buf, nil
}

func (s *WSBuf) Append1Byte(val int) error {
	if err := s.Br.Write(bitbyte.BITS_PER_BYTE, uint64(val)); err != nil {
		return fmt.Errorf("WSBuf.Append1Byte:Failed to write bit value.\n\t%+v", err)
	}
	return nil
}

func (s *WSBuf) Append2Bytes(val int) error {
	if err := s.Br.Write(bitbyte.BITS_PER_BYTE*2, uint64(val)); err != nil {
		return fmt.Errorf("WSBuf.Append2Bytes:Failed to write bit value.\n\t%+v", err)
	}
	return nil
}
func (s *WSBuf) Append3Bytes(val int) error {
	if err := s.Br.Write(bitbyte.BITS_PER_BYTE*3, uint64(val)); err != nil {
		return fmt.Errorf("WSBuf.Append3Bytes:Failed to write bit value.\n\t%+v", err)
	}
	return nil
}

func (s *WSBuf) Append4Bytes(val int) error {
	if err := s.Br.Write(bitbyte.BITS_PER_BYTE*4, uint64(val)); err != nil {
		return fmt.Errorf("WSBuf.Append4Bytes:Failed to write bit value.\n\t%+v", err)
	}
	return nil
}

func (s *WSBuf) Append8Bytes(val int64) error {
	if err := s.Br.Write(bitbyte.BITS_PER_BYTE*8, uint64(val)); err != nil {
		return fmt.Errorf("WSBuf.Append8Bytes:Failed to write bit value.\n\t%+v", err)
	}
	return nil
}

func (s *WSBuf) AppendString(str string) error {
	if err := s.Br.WriteString(str); err != nil {
		return fmt.Errorf("wsbuf.WSBuf.AppendString:bitbyte.BitRow.WriteString failed.\n\t%v", err)
	}
	return nil
}

func (s *WSBuf) Get1Byte() (int, error) {
	var val uint64
	var err error
	if val, err = s.Br.Read(bitbyte.BITS_PER_BYTE); err != nil {
		return 0, fmt.Errorf("WSBuf.Get1Byte:Failed to read bit value s.Br:%v.\n\t%+v", s.Br, err)
	}
	return int(val), nil
}

func (s *WSBuf) Get2Bytes() (int, error) {
	var val uint64
	var err error
	if val, err = s.Br.Read(bitbyte.BITS_PER_BYTE * 2); err != nil {
		return 0, fmt.Errorf("WSBuf.Get2Bytes:Failed to read bit value s.Br:%v.\n\t%+v", s.Br, err)
	}
	return int(val), nil
}

func (s *WSBuf) Geta3Bytes() (int, error) {
	var val uint64
	var err error
	if val, err = s.Br.Read(bitbyte.BITS_PER_BYTE * 3); err != nil {
		return 0, fmt.Errorf("WSBuf.Get3Bytes:Failed to read bit value s.Br:%v.\n\t%+v", s.Br, err)
	}
	return int(val), nil
}

func (s *WSBuf) Get4Bytes() (int, error) {
	var val uint64
	var err error
	if val, err = s.Br.Read(bitbyte.BITS_PER_BYTE * 4); err != nil {
		return 0, fmt.Errorf("WSBuf.Get4Bytes:Failed to read bit value s.Br:%v.\n\t%+v", s.Br, err)
	}
	return int(val), nil
}

func (s *WSBuf) Get8Bytes() (int64, error) {
	var val uint64
	var err error
	if val, err = s.Br.Read(bitbyte.BITS_PER_BYTE * 8); err != nil {
		return 0, fmt.Errorf("WSBuf.Get8Bytes:Failed to read bit value s.Br:%v.\n\t%+v", s.Br, err)
	}
	return int64(val), nil
}

func (s *WSBuf) GetString() (string, error) {
	var err error
	var sz uint64
	if sz, err = s.Br.Read(bitbyte.BR_STR_SIZE); err != nil {
		return "", fmt.Errorf("WSBuf.GetString:Can't get size.\n\t%+v", err)
	}
	r := make([]rune, sz)
	for i := 0; int(sz) > i; i++ {
		if w, err := s.Br.Read(WS_WORD_SIZE); err != nil {
			return "", fmt.Errorf("WSBuf.GetString:Can't get rune index:%d.\n\t%+v", i, err)
		} else {
			r[i] = rune(w)
		}
	}
	return string(r), nil
}

func CalcStrSize(str string) int {
	return len([]rune(str))*CHAR_SIZE + STR_SIZE_LENGTH
}

func (s *WSBuf) CheckContinue() bool {
	return s.Br.Index+BYTE4_SIZE*bitbyte.BITS_PER_BYTE < len(s.Br.Body)*bitbyte.BITS_PER_VALUE
}
