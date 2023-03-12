package stk

import (
	"log"
	"strconv"
	"fmt"
	"io"
	"math"

	"gnocchi/gim"
)

// 文字列マネージャー設定ファイル列名
const (
    STRMGR_PARAM_NAME = iota
    STRMGR_PARAM_FILE
    NUM_STRMGR_PARAM
)

// 文字列読みとり遷移状態
const (
  STRMODE_FIRST_SEPARATE = iota
  STRMODE_ID
  STRMODE_SECOUND_SEPARATE
  STRMODE_BODY
  STRMODE_ESCAPE
  STRMODE_END_SEPARATE
)

type category struct {
  smp map[int]string
}

type ctgry struct {
	key			string
	filename	string
}

type item struct {
	id		int
	str		string
}

type StringStocker struct {
  ctry  map[ string ]*category
}

func newCategory( text string ) ( interface{}, error ) {
	var c ctgry
	if num, err := fmt.Sscanf( text, "%s %s", &c.key, &c.filename ); err != nil {
		return nil, fmt.Errorf( "stk.Create:fmt.Sscanf failure num num:%d, text:%q.\n\t%v", num, text, err )
	} else if num < NUM_STRMGR_PARAM {
		return nil, fmt.Errorf( "stk.Create:number of item is too short. item num:%d, buf:%q.\n\t%v", num, text, err )
	}
	return &c, nil
}

func newItem( text string ) ( interface{}, error ) {
	src := []rune( text )
	sz := len( text )
	dst := make( []rune, sz )
	mode := STRMODE_FIRST_SEPARATE
	index, n := 0, 0
	var it item
	var err error
	for i, r := range src {
		switch mode {
			case STRMODE_FIRST_SEPARATE:
				if r >= 0x30 && r <= 0x39 {
					mode = STRMODE_ID
					if n, err = strconv.Atoi( string( r ) ); err != nil {
						return nil, fmt.Errorf( "stk.newItem:Atoi failure r:%x, i:%d, sz:%d, err:%v, mode:%d, i:%d, src:%q", r, i, sz, err, mode, i, src )
					}
					it.id = int( n )
				} else if r != 0x20 {
					return nil, fmt.Errorf( "stk.newItem:Invalid rune r:%x, i:%d, sz:%d, mode:%d, i:%d, src:%q", r, i, sz, mode, i, src )
				}
			case STRMODE_ID:
				if r >= 0x30 && r <= 0x39 {
					if n, err = strconv.Atoi( string( r ) ); err != nil {
						return nil, fmt.Errorf( "stk.newItem:Atoi failure r:%x, i:%d, sz:%d, err:%v, mode:%d, i:%d, src:%q", r, i, sz, err, mode, i, src )
					} else if it.id > math.MaxInt64 / 10 {
						return nil, fmt.Errorf( "stk.newItem:id is too large id:%d, r:%x, i:%d, sz:%d, mode:%d, i:%d, src:%q", it.id, r, i, sz, mode, i, src )
					}
					it.id = it.id * 10 + int( n )
				} else if r == 0x20 {
					mode = STRMODE_SECOUND_SEPARATE
				} else {
					return nil, fmt.Errorf( "stk.newItem:Invalid rune r:%x, it.id:%d, i:%d, sz:%d, mode:%d, i:%d, src:%q", r, it.id, i, sz, mode, i, src )
				}
			case STRMODE_SECOUND_SEPARATE:
				if r == '"' {
					mode = STRMODE_BODY
				} else if r != 0x20 {
					return nil, fmt.Errorf( "stk.newItem:Invalid rune r:%x, i:%d, sz:%d, mode:%d, src:%q", r, i, sz, mode, src )
				}
			case STRMODE_BODY:
				if r == '"' {
					mode = STRMODE_END_SEPARATE
				} else if r == '\\' {
					mode = STRMODE_ESCAPE
				} else {
					dst[ index ] = r
					index++
				}
			case STRMODE_ESCAPE:
				if r == '"' || r == '\\' {
					dst[ index ] = r
					index++
					mode = STRMODE_BODY
				} else {
					return nil, fmt.Errorf( "stk.newItem:Invalid escape r:%x, i:%d, sz:%d, mode:%d, i:%d, src:%q", r, i, sz, mode, i, src )
				}
			case STRMODE_END_SEPARATE:
				if r != ' ' {
					return nil, fmt.Errorf( "stk.newItem:Invalid end separate r:%x, i:%d, sz:%d, mode:%d, i:%d, src:%q", r, i, sz, mode, i, src )
				}
		}
	}
	if mode == STRMODE_END_SEPARATE {
		it.str = string( dst[ :index ] )
//	fmt.Printf( "it:%v\n", it )
	} else {
		return nil, fmt.Errorf( "stk.newItem:Invalid end of line sz:%d, mode:%d, src:%q.", sz, mode, src )
	}

	return &it, nil
}

func Create( root, path, file string ) ( *StringStocker, error ) {
	log.Printf( "Start stk.Create root:%q paht:%q file:%q.", root, path, file )
	defer log.Printf( "End stk.Create root:%q path:%q file:%q.", root, path, file )

	line, num, err := gim.CreateFileLines( root + path + file, newCategory )
	if err != nil {
		return nil, fmt.Errorf( "stk.Create:gim.CreateFileLines failure root:%q path:%q file:%q.\n\t%v", root, path, file, err )
	}

	stkr := &StringStocker{ ctry:make( map[ string ]*category, num ) }

	var c *ctgry
	for line != nil {
		c = (line.Data).( *ctgry )
		if err = stkr.appendCategory( root, path, c.key, c.filename ); err != nil {
			return nil, fmt.Errorf( "stk.Create:stk.appendCategory failure num:%d, c:%v, file:%q.\n\t%v", num, c, file, err )
		}
		line = line.Next
	}
	return stkr, nil
}

func ( strStk *StringStocker ) appendCategory( root, path, name, file string ) error {
	line, num, err := gim.CreateFileLines( root + path + file, newItem )
	if err != nil {
		return fmt.Errorf( "stk.StringStocker.appendCategory:gim.CreateFileLines failure root:%q path:%q file:%q.\n\t%v", root, path, file, err )
	}

	var c category
	c.smp = make( map[ int ]string, num )

	var it *item
	for line != nil {
		it = (line.Data).( *item )
		c.smp[ it.id ] = it.str
		line = line.Next
	}

	strStk.ctry[ name ] = &c

	return nil
}

func ( strStk *StringStocker ) String( category string, index int ) ( string, error ) {
	if c, ok := strStk.ctry[ category ]; !ok {
		return "", fmt.Errorf( "stk.GetString:Invalid category:%q, index:%d.", category, index  )
	} else if s, ok := c.smp[ index ]; !ok {
		return "", fmt.Errorf( "stk.GetString:Invalid index:%d, category:%q.", index, category )
	} else {
		return s, nil
	}
}

func ( strStk *StringStocker ) GetString( category string, index int ) string {
	str, _ := strStk.String( category, index )
	return str
}

func ( strStk *StringStocker ) GetCategorySize( category string ) int {
	c, ok := strStk.ctry[ category ]
	if c == nil || !ok {  // 格納されている以上は、nilではないけど、念のため
		log.Printf( "stk.GetSize:Invalid category:%q", category )
		return -1
	}
	return len( c.smp )
}

func ( strStk *StringStocker ) Dump( w io.Writer ) {
	for key, v := range strStk.ctry {
		fmt.Fprintf( w, "CATEGORY key:%s, body%v\n", key, v )
		for i, s := range v.smp {
			fmt.Fprintf( w, "  STRING i:%3d, str:%s\n", i, s )
		}
	}
}

func ( strStk *StringStocker ) LogAll() {
	for key, v := range strStk.ctry {
		log.Printf( "CATEGORY key:%s, body%v\n", key, v )
		for i, s := range v.smp {
			log.Printf( "  STRING i:%3d, str:%s\n", i, s )
		}
	}
}
