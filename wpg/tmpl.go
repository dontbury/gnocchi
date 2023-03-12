package wpg

import (
	"fmt"
	"strings"
	"net/http"
	"text/template"
	"net/url"
//	"os"
//	"bufio"

//	"github.com/dontbury/gnocchi/bkt"
	"gnocchi/gim"
)

const (
    TMPL_PARAM_KEY = iota
    TMPL_PARAM_FILE
    NUM_TMPL_PARAM
)

const	EXECUTE_TEMPLATE_NAME = "page"

type TmplMap struct {
  tmpls  map[ string ]*template.Template
}

type TmplFuncArgElm struct {
	AtFirst bool
	Seq int
	Req *http.Request
	Path *Paths
	Arg *url.Values
	Conn *Conn
	Param01 string
	Param02 string
	Param03 string
	Param04 string
	Param05 string
	Param06 string
	Param07 string
	Param08 string
	Param09 string
	Param10 string
	Param11 string
	Param12 string
	Param13 string
	Param14 string
	Param15 string
	Param16 string
	Param17 string
	Param18 string
	Param19 string
	Param20 string
}

type TmplFuncArg struct {
	Index int
	Elm *TmplFuncArgElm
}

type tmpl struct {
	key		string
	page	string
	head	string
	header	string
	nav		string
	main	string
	footer	string
}

func newTemplate( text string ) ( interface{}, error ) {
	var t tmpl
	num, err := fmt.Sscanf( text, "%s %s %s %s %s %s %s", &t.key, &t.page, &t.head, &t.header, &t.nav, &t.main, &t.footer )
	if err != nil {
		return nil, fmt.Errorf( "stk.Create:fmt.Sscanf failure num num:%d, text:%q.\n\t%v", num, text, err )
	} else if num < NUM_TMPL_PARAM {
		return nil, fmt.Errorf( "stk.Create:number of item is too short. item num:%d, text:%q.\n\t%v", num, text, err )
	}
	t.key = strings.Trim( t.key, "\"" )
	return &t, nil
}

func ( mgr *TmplMap ) GetTemplate( key string ) ( *template.Template, bool ) {
	v, ok := mgr.tmpls[ key ]
	return v, ok
}

func Create( root, file, path, static string, funcMap template.FuncMap ) ( *TmplMap, error ) {
	line, num, err := gim.CreateFileLines( root + file, newTemplate )
	if err != nil {
		return nil, fmt.Errorf( "wpg.Create:gim.CreateFileLines failure file:%q.\n\t%v.", file, err )
	}

	mgr := &TmplMap{ tmpls:make( map[ string ]*template.Template, num ) }

	tmplPath, tmplStatic := root + path + "/", root + path + static
	var t *tmpl
//	fmt.Printf( "tmplStatic:%q\n", tmplStatic )
	for line != nil {
		t = (line.Data).( *tmpl )
//	fmt.Printf( "wpg.Create:key:%q line:%v\n", t.key, t )
		mgr.tmpls[ t.key ] = template.Must( template.New( tmplPath + t.page ).Funcs( funcMap ).ParseFiles( 
			tmplPath + t.page, tmplPath + t.head, tmplPath + t.header, 
			tmplPath + t.nav, tmplPath + t.main, tmplPath + t.footer, tmplStatic ) )
		line = line.Next
	}

	return mgr, nil
}

func TmplFunc( tmplFunc func( arg *TmplFuncArg ) error, arg *TmplFuncArg ) ( string, error ) {
	if arg.Elm.AtFirst {
		arg.Elm.AtFirst = false
		if err := tmplFunc( arg ); err != nil { return "", fmt.Errorf( "wpg.TmplFunc func failure err:%v", err ) }
	}

	var str string
	switch arg.Index {
		case  1: str = arg.Elm.Param01
		case  2: str = arg.Elm.Param02
		case  3: str = arg.Elm.Param03
		case  4: str = arg.Elm.Param04
		case  5: str = arg.Elm.Param05
		case  6: str = arg.Elm.Param06
		case  7: str = arg.Elm.Param07
		case  8: str = arg.Elm.Param08
		case  9: str = arg.Elm.Param09
		case 10: str = arg.Elm.Param10
		case 11: str = arg.Elm.Param11
		case 12: str = arg.Elm.Param12
		case 13: str = arg.Elm.Param13
		case 14: str = arg.Elm.Param14
		case 15: str = arg.Elm.Param15
		case 16: str = arg.Elm.Param16
		case 17: str = arg.Elm.Param17
		case 18: str = arg.Elm.Param18
		case 19: str = arg.Elm.Param19
		case 20: str = arg.Elm.Param20
		default: return "", fmt.Errorf( "wpg.TmplFunc Invalid arg.Index:%v arg:%#v", arg.Index, arg )
	}
	return str, nil
}

func ( elm *TmplFuncArgElm ) IncSet( p string ) error {
	switch elm.Seq {
		case  1: elm.Param01 = p
		case  2: elm.Param02 = p
		case  3: elm.Param03 = p
		case  4: elm.Param04 = p
		case  5: elm.Param05 = p
		case  6: elm.Param06 = p
		case  7: elm.Param07 = p
		case  8: elm.Param08 = p
		case  9: elm.Param09 = p
		case 10: elm.Param10 = p
		case 11: elm.Param11 = p
		case 12: elm.Param12 = p
		case 13: elm.Param13 = p
		case 14: elm.Param14 = p
		case 15: elm.Param15 = p
		case 16: elm.Param16 = p
		case 17: elm.Param17 = p
		case 18: elm.Param18 = p
		case 19: elm.Param19 = p
		case 20: elm.Param20 = p
		default: return fmt.Errorf( "wpg.TmplFuncArgElm.IncSet Invalid Seq:%v.", elm.Seq )
	}
	elm.Seq++

	return nil
}
