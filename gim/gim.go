package gim

import (
	"fmt"
	"os"
	"bufio"
)

type FileLine struct {
	Data	interface{}
	Next	*FileLine
}

func NewFileItem( text string ) ( interface{}, error ) {
	return &text, nil
}

func CreateFileLines( file string, New func( line string ) ( interface{}, error ) ) ( *FileLine, int, error ) {
	fp, err := os.Open( file )
	if err != nil {
		return nil, 0, fmt.Errorf( "gim.CreateFileLines:os.Open failure file:%q.\n\t%v.", file, err )
	}
	defer fp.Close();

	var first, last, line *FileLine = nil, nil, nil
	lineNum := 0
	var text string
	scanner := bufio.NewScanner( fp )
	for scanner.Scan() {
		text = scanner.Text()
		if len( text ) > 0 { // 空白行（末尾など）なら参照しない
			if text[ 0 ] != '#' {  // 行頭が#ならコメント行なので参照しない
				line = &FileLine{ Next:nil }
				if line.Data, err = New( text ); err != nil {
					return nil, 0, fmt.Errorf( "gim.CreateFileLines:New failure text:%q.\n\t%v.", text, err )
				}
				if first == nil {
					first = line
				}
				if last != nil {
					last.Next = line
				}
				last = line
				lineNum++
//	fmt.Printf( "lineNum:%d line:%v\n", lineNum, line.Data )
			}
		}
	}

	return first, lineNum, nil
}

func IntMin( a, b int ) int {
	if a < b {
		return a
	} else {
		return b
	}
}
