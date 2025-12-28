package gim

import "fmt"

type BTreeItem interface {
	Compare(BTreeItem) (bool, error)
}

type BTree struct {
	it BTreeItem
	p  *BTree
	r  *BTree
	l  *BTree
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
