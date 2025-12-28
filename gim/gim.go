package gim

import (
	"bufio"
	"fmt"
	"os"
)

type FileLine struct {
	Data interface{}
	Next *FileLine
}

func NewFileItem(text string, index int, src interface{}) (interface{}, error) {
	return &text, nil
}

func CreateFileLines(file string, src interface{}, New func(string, int, interface{}) (interface{}, error)) (*FileLine, int, error) {
	var first, last, line *FileLine = nil, nil, nil
	cnt := 0
	if fp, err := os.Open(file); err == nil {
		defer fp.Close()
		var text string
		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			text = scanner.Text()
			if len(text) > 0 { // 空白行（末尾など）なら参照しない
				if text[0] != '#' { // 行頭が#ならコメント行なので参照しない
					line = &FileLine{Next: nil}
					if line.Data, err = New(text, cnt, src); err != nil {
						return nil, 0, fmt.Errorf("gim.CreateFileLines:New failure text:%q.\n\t%v", text, err)
					}
					if first == nil {
						first = line
					}
					if last != nil {
						last.Next = line
					}
					last = line
					cnt++
					//	fmt.Printf( "lineNum:%d line:%v\n", lineNum, line.Data )
				}
			}
		}

	} else {
		return nil, 0, fmt.Errorf("gim.CreateFileLines:os.Open failure file:%q.\n\t%v", file, err)
	}
	return first, cnt, nil
}
