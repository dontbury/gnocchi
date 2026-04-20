package wskt

import (
	"fmt"
	"log"

	"github.com/dontbury/gnocchi/bitbyte"
)

const (
	BYTE1_SIZE = 1
	BYTE2_SIZE = 2
	BYTE3_SIZE = 3
	BYTE4_SIZE = 4
	BYTE8_SIZE = 8

	WS_STR_SIZE = 16	// WebSocketで送受信する文字列の一般的なサイズを規定するビット数
	WS_WORD_SIZE = 16	// WebSocketで送受信する文字のビット数
	STR_SIZE_LENGTH = BYTE2_SIZE // 文字列の長さを収めるためのバイト数
	CHAR_SIZE       = BYTE2_SIZE // 1文字あたりのバイト数

	WSKTBUF_INCREMENT = 1024
)

type WSBuf struct {
	Br bitbyte.BitRow
}

func (s *WSBuf) Create(index, size int) {
	s.Br = bitbyte.BitRow{Index: 0, Inc: 0, Body: make([]uint64, size)}
}

func (s *WSBuf) CreateBuf(index, inc int, buf *[]byte) {
	s.Br = bitbyte.BitRow{Index: 0, Inc: inc, Body: []uint64{}}
	sz := (len(*buf) + bitbyte.BYTES_PER_UINT64 - 1) / bitbyte.BYTES_PER_UINT64
	if sz > 0 {
		s.Br.Body = make([]uint64, sz)
		for i, v := range *buf {
			s.Br.Body[i>>3] |= uint64(v) << (bitbyte.BYTES_PER_UINT64 * (i & (bitbyte.BYTES_PER_UINT64 - 1)))
		}
	}
}

func (s *WSBuf) CreateWSB(index, headersize int, wsb *WSBuf) error {
	sz := (headersize + (len(wsb.Br.Body)+1)*bitbyte.BYTES_PER_UINT64 - 1 - wsb.Br.Index) / bitbyte.BYTES_PER_UINT64
	s.Br = bitbyte.BitRow{Index: index, Inc: 0, Body: make([]uint64, sz)}
	for i := headersize; i < sz; i++ {
		s.Br.Body[i] |= wsb.Br.Body[wsb.Br.Index+i-headersize]
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
func (s *WSBuf) GetRawBuf() *[]byte {
	buf := make([]byte, len(s.Br.Body)*bitbyte.BYTES_PER_UINT64)
	for i, v := range s.Br.Body {
		for j := 0; j < bitbyte.BYTES_PER_UINT64; j++ {
			buf[i*bitbyte.BYTES_PER_UINT64+j] = byte((v >> (bitbyte.BYTES_PER_UINT64 * j)) & 0xFF)
		}
	}
	return &buf
}

// 送信用に後ろの未使用部分を切り詰めたバッファを取得
func (s *WSBuf) GetSendBuf() (*[]byte, error) {
	buf := []byte{}
	if s.Br.Index < len(s.Br.Body) {
		sz := (len(s.Br.Body)*bitbyte.BITS_PER_VALUE - s.Br.Index + bitbyte.BITS_PER_BYTE - 1) / bitbyte.BITS_PER_BYTE
		buf = make([]byte, sz)
		var val uint64
		var err error
		for i := 0; len(buf) > i; i++ {
			if sz > bitbyte.BITS_PER_VALUE {
				if val, err = s.Br.Get(i*bitbyte.BITS_PER_BYTE, bitbyte.BITS_PER_BYTE); err != nil {
					return nil, fmt.Errorf("WSBuf.GetSendBuf:Failed to get bit value.\t\n%+v", err)
				}
				buf[i] = byte(val)
				sz -= bitbyte.BITS_PER_VALUE
			} else {
				if val, err = s.Br.Get(i*bitbyte.BITS_PER_BYTE, sz); err != nil {
					return nil, fmt.Errorf("WSBuf.GetSendBuf:Failed to get bit value.\t\n%+v", err)
				}
				buf[i] = byte(val)
				sz = 0
			}
		}
		log.Printf("WSBuf.GetSendBuf: index:%d size:%d %v -> %v.", s.Br.Index, len(s.Br.Body), s.Br.Body, buf)
	} else {
		buf = make([]byte, len(s.Br.Body)*bitbyte.BITS_PER_VALUE/bitbyte.BITS_PER_BYTE)
		for i, v := range s.Br.Body {
			for j := 0; j < bitbyte.BYTES_PER_UINT64; j++ {
				buf[i*bitbyte.BYTES_PER_UINT64+j] = byte((v >> (bitbyte.BYTES_PER_UINT64 * j)) & 0xFF)
			}
		}
		log.Printf("WSBuf.GetSendBuf: index:%d size:%d %v -> %v.", s.Br.Index, len(s.Br.Body), s.Br.Body, buf)
	}
	return &buf, nil
}

func (s *WSBuf) Append1Byte(val int) error {
	if err := s.Br.Write(val, bitbyte.BITS_PER_BYTE); err != nil {
		return fmt.Errorf("WSBuf.Append1Byte:Failed to write bit value.\t\n%+v", err)
	}
	return nil
}

func (s *WSBuf) Append2Bytes(val int) error {
	if err := s.Append1Byte(val / 0x100); err != nil {
		return fmt.Errorf("WSData.Append2Bytes:Append1Byte 1 failure.\t\n%+v", err)
	}
	if err := s.Append1Byte(val % 0x100); err != nil {
		return fmt.Errorf("WSData.Append2Bytes:Append1Byte 2 failure.\t\n%+v", err)
	}
	return nil
}
func (s *WSBuf) Append3Bytes(val int) error {
	if err := s.Append1Byte(val / 0x10000); err != nil {
		return fmt.Errorf("WSData.Append3Bytes:Append1Byte failure.\t\n%+v", err)
	}
	if err := s.Append2Bytes(val % 0x10000); err != nil {
		return fmt.Errorf("WSData.Append3Bytes:Append2Bytes failure.\t\n%+v", err)
	}
	return nil
}

func (s *WSBuf) Append4Bytes(val int) error {
	if err := s.Append2Bytes(val / 0x10000); err != nil {
		return fmt.Errorf("WSData.Append4Bytes:Append2Bytes 1 failure.\t\n%+v", err)
	}
	if err := s.Append2Bytes(val % 0x10000); err != nil {
		return fmt.Errorf("WSData.Append4Bytes:Append2Bytes 2 failure.\t\n%+v", err)
	}
	return nil
}

func (s *WSBuf) Append8Bytes(val int64) error {
	if err := s.Append4Bytes((int)(val / 0x100000000)); err != nil {
		return fmt.Errorf("WSData.Append8Bytes:Append4Bytes 1 failure.\t\n%+v", err)
	}
	if err := s.Append4Bytes((int)(val % 0x100000000)); err != nil {
		return fmt.Errorf("WSData.Append8Bytes:Append4Bytes 2 failure.\t\n%+v", err)
	}
	return nil
}

func (s *WSBuf) AppendString(str string) error {
	runeStr := []rune(str)
	sz := len(runeStr)
	if sz >= 0x10000 { // 2倍して2バイトに収まる必要がある
		return fmt.Errorf("WSData.AppendString:val too long:%d.", sz)
	} else {
		if err := s.Br.Write(WS_STR_SIZE, uint64(sz)); err != nil {
			return fmt.Errorf("WSData.AppendString:Failed to write string size.\t\n%+v", err)
		} else {
			for _, r := range runeStr {
				if err = s.Br.Write(WS_WORD_SIZE, uint64(r)); err != nil {
					return fmt.Errorf("WSData.AppendString:Failed to write rune:%v, str:%s s.Br:%v.\t\n%+v", err, str, s.Br, err)
				}
			}
		}
	}
	return nil
}

func (s *WSBuf) Get1Byte() (int, error) {
	var val uint64
	var err error
	if val, err = s.Br.Read(bitbyte.BITS_PER_BYTE); err != nil {
		return 0, fmt.Errorf("WSBuf.Get1Byte:Failed to read bit value s.Br:%v.\t\n%+v", s.Br, err)
	}
	return int(val), nil
}

func (s *WSBuf) Get2Bytes() (int, error) {
	var val uint64
	var err error
	if val, err = s.Br.Read(bitbyte.BITS_PER_BYTE*2); err != nil {
		return 0, fmt.Errorf("WSBuf.Get2Bytes:Failed to read bit value s.Br:%v.\t\n%+v", s.Br, err)
	}
	return int(val), nil
}

func (s *WSBuf) Geta3Bytes() (int, error) {
	var val uint64
	var err error
	if val, err = s.Br.Read(bitbyte.BITS_PER_BYTE*3); err != nil {
		return 0, fmt.Errorf("WSBuf.Get3Bytes:Failed to read bit value s.Br:%v.\t\n%+v", s.Br, err)
	}
	return int(val), nil
}

func (s *WSBuf) Get4Bytes() (int, error) {
	var val uint64
	var err error
	if val, err = s.Br.Read(bitbyte.BITS_PER_BYTE*4); err != nil {
		return 0, fmt.Errorf("WSBuf.Get4Bytes:Failed to read bit value s.Br:%v.\t\n%+v", s.Br, err)
	}
	return int(val), nil
}

func (s *WSBuf) Get8Bytes() (int64, error) {
	var val uint64
	var err error
	if val, err = s.Br.Read(bitbyte.BITS_PER_BYTE*8); err != nil {
		return 0, fmt.Errorf("WSBuf.Get8Bytes:Failed to read bit value s.Br:%v.\t\n%+v", s.Br, err)
	}
	return int64(val), nil
}

func (s *WSBuf) GetString() (string, error) {
	var err error
	var sz uint64
	if sz, err = s.Br.Read(WS_STR_SIZE); err != nil {
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
	return s.Br.Index < len(s.Br.Body)*bitbyte.BITS_PER_VALUE
}
