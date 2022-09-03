package stk

import (
	"log"
	"strconv"
	"fmt"
	"io"
	"math"

	"gnocchi/bkt"
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

type ctgry struct {
  smp map[int]string
}

type StringStocker struct {
  ctry  map[string]*ctgry
}

func Create( bckt *bkt.Bucket, path string ) ( *StringStocker, error ) {
	log.Println( "Start stk.ReadData." )
	buf, err := bckt.ReadLine( path )
	if err != nil {
		return nil, fmt.Errorf( "stk.Create:bkt.Bucket.ReadLine failure bckt:%v.\n\t%v.", bckt, err )
	}
	lineNum := len( buf )
	if lineNum <= 0 {
		return nil, fmt.Errorf( "stk.Create: too short lineNum:%d, bckt:%v.", lineNum, bckt )
	}

	mgr := &StringStocker{ ctry:  make( map[string]*ctgry, lineNum ) }

	for i := 0 ; i < lineNum; i++ {
		var name, file string
		if num, err := fmt.Sscanf( buf[i], "%s %s", &name, &file ); num < NUM_STRMGR_PARAM || err != nil {
			return mgr, fmt.Errorf( "stk.Create: too short item num:%d, lineNum:%d, i:%d, buf:%q.\n\t%v", num, lineNum, i, buf[i], err )
		}
		if err := mgr.createStrMap( bckt, name, file ); err != nil {
			return mgr, fmt.Errorf( "stk.Create:stk.createStrMap failure lineNum:%d, i:%d, buf:%q, file:%q.\n\t%v", lineNum, i, buf[i], file, err )
		}
	}
	log.Println( "Complete stk.ReadData." )
	return mgr, nil
}

func ( strStk *StringStocker ) createStrMap( bckt *bkt.Bucket, name, path string ) error {
	buf, err := bckt.ReadLine( path )
	if err != nil {
		return fmt.Errorf( "stk.StringStocker.createStrMap:bkt.ReadLine failure path:%q.\n\t%v", path, err )
	}
	lineNum := len( buf )
	if lineNum <= 0 {
		return fmt.Errorf( "stk.StringStocker.createStrMap:too short lineNum:%d, path:%q.\n\t%v", lineNum, path, err )
	}

	p := new( ctgry )
	p.smp = make( map[ int ]string, lineNum )
	for i := 0 ; i < lineNum; i++ {
		src := []rune( buf[ i ] )
		sz := len( src )
		dst := make([]rune, sz)
		mode := STRMODE_FIRST_SEPARATE
		id, index := 0, 0
		for j := 0; j < sz; j++ {
			r := src[j]
			switch mode {
			case STRMODE_FIRST_SEPARATE:
				if r >= 0x30 && r <= 0x39 {
					mode = STRMODE_ID
					n, err := strconv.Atoi( string( r ) )
					if err != nil {
						return fmt.Errorf( "stk.StringStocker.createStrMap:Atoi failure r:%x, j:%d, sz:%d, err:%v, mode:%d, lineNum:%d, i:%d, buf:%q, path:%q", r, j, sz, err, mode, lineNum, i, buf[i], path )
					}
					id = int( n )
				} else if r != 0x20 {
					return fmt.Errorf( "stk.StringStocker.createStrMap:Invalid rune r:%x, j:%d, sz:%d, mode:%d, lineNum:%d, i:%d, buf:%q, path:%q", r, j, sz, mode, lineNum, i, buf[i], path )
				}
			case STRMODE_ID:
				if r >= 0x30 && r <= 0x39 {
					n, err := strconv.Atoi( string( r ) )
					if err != nil {
						return fmt.Errorf( "stk.StringStocker.createStrMap:Atoi failure r:%x, j:%d, sz:%d, err:%v, mode:%d, lineNum:%d, i:%d, buf:%q, path:%q", r, j, sz, err, mode, lineNum, i, buf[i], path )
					}
					if id > math.MaxInt64 / 10 {
						return fmt.Errorf( "stk.StringStocker.createStrMap:id is too large id:%d, r:%x, j:%d, sz:%d, mode:%d, lineNum:%d, i:%d, buf:%q, rpath:%q", id, r, j, sz, mode, lineNum, i, buf[i], path )
					}
					id = id * 10 + int( n )
				} else if r == 0x20 {
					mode = STRMODE_SECOUND_SEPARATE
				} else {
					return fmt.Errorf( "stk.StringStocker.createStrMap:invalid rune r:%x, id:%d, j:%d, sz:%d, mode:%d, lineNum:%d, i:%d, buf:%q, path:%q", r, id, j, sz, mode, lineNum, i, buf[i], path )
				}
			case STRMODE_SECOUND_SEPARATE:
				if r == '"' {
					mode = STRMODE_BODY
				} else if r != 0x20 {
					return fmt.Errorf( "stk.StringStocker.createStrMap:invalid rune r:%x, j:%d, sz:%d, mode:%d, lineNum:%d, i:%d, buf:%q, path:%q", r, j, sz, mode, lineNum, i, buf[i], path )
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
					return fmt.Errorf( "stk.StringStocker.createStrMap:invalid escape r:%x, j:%d, sz:%d, mode:%d, lineNum:%d, i:%d, buf:%q, path:%q", r, j, sz, mode, lineNum, i, buf[i], path )
				}
			case STRMODE_END_SEPARATE:
				if r != ' ' {
					return fmt.Errorf( "stk.StringStocker.createStrMap:invalid end separate r:%x, j:%d, sz:%d, mode:%d, lineNum:%d, i:%d, buf:%q, path:%q", r, j, sz, mode, lineNum, i, buf[i], path )
				}
			}
		}
		var s string
		if mode == STRMODE_END_SEPARATE {
			s = string( dst[ :index ] )
		}
		p.smp[id] = s
	}
	strStk.ctry[ name ] = p
	return nil
}

func ( strStk *StringStocker ) String( category string, index int ) ( string, error ) {
	if c, ok := strStk.ctry[ category ]; !ok {
		return "", fmt.Errorf( "stk.GetString:invalid category:%q, index:%d.", category, index  )
	} else if s, ok := c.smp[ index ]; !ok {
		return "", fmt.Errorf( "stk.GetString:invalid index:%d, category:%q.", index, category )
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
		log.Printf( "stk.GetSize:invalid category:%q", category )
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
