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

type BTreeItem interface {
	Compare(BTreeItem) (bool, error)
}

type BTree struct {
	it BTreeItem
	p  *BTree
	r  *BTree
	l  *BTree
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

func IntMin(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func Sort(src, dst []BTreeItem) error {
	var root, t *BTree = nil, nil
	for i, v := range src {
		if root == nil {
			root = &BTree{it: src[i], p: nil, l: nil, r: nil}
		} else {
			t = root
			for {
				if t.it == nil {
					t.it = src[i]
					break
				} else if more, err := v.Compare(t.it); err != nil {
					return fmt.Errorf("gnocchi.gim.Sort failure\n\t%v", err)
				} else if more {
					if t.r == nil {
						t.r = &BTree{it: src[i], p: t, l: nil, r: nil}
						break
					}
					t = t.r
				} else {
					if t.l == nil {
						t.l = &BTree{it: src[i], p: t, l: nil, r: nil}
						break
					}
					t = t.l
				}
			}
		}
	}
	index := 0
	t = root
/*
	for {
		if t.r != nil {
			t = t.r
			fmt.Println("gnocchi.gim.Sort move right")
		} else {
			dst[index] = t.it
			fmt.Printf("gnocchi.gim.Sort index:%d t.it:%+v\n", index, t.it)
			index++
			if t.l != nil {
				t = t.l
				fmt.Println("gnocchi.gim.Sort move left")
			} else if t.p != nil {
				t = t.p
				fmt.Println("gnocchi.gim.Sort move previous")
			} else {
				fmt.Println("gnocchi.gim.Sort move exit")
				break
			}
		}
	}
*/
	root.set(dst, &index)
	return nil
}

func (t *BTree) set(s []BTreeItem, index *int) {
	if t.r != nil {
		t.r.set(s, index)
	}
	s[*index] = t.it
	*index++
	if t.l != nil {
		t.l.set(s, index)
	}
}
