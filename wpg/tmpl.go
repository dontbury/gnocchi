package wpg

import (
	"bufio"
	"embed"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"

	"github.com/dontbury/gnocchi/gim"
)

const (
	TMPL_PARAM_KEY = iota
	TMPL_PARAM_PAGE
	TMPL_PARAM_HEAD
	TMPL_PARAM_HEADER
	TMPL_PARAM_NAV
	TMPL_PARAM_MAIN
	TMPL_PARAM_FOOTER
	NUM_TMPL_PARAM
)

const EXECUTE_TEMPLATE_NAME = "page"

type TmplMap struct {
	tmpls map[string]*template.Template
}

type TmplFuncArgElm struct {
	AtFirst bool
	Seq     int
	Req     *http.Request
	Path    *Paths
	Arg     *url.Values
	Conn    *Conn
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
	Param21 string
	Param22 string
	Param23 string
	Param24 string
	Param25 string
	Param26 string
	Param27 string
	Param28 string
	Param29 string
	Param30 string
}

type TmplFuncArg struct {
	Index int
	Elm   *TmplFuncArgElm
}

type tmpl struct {
	key    string
	page   string
	head   string
	header string
	nav    string
	main   string
	footer string
}

func newTemplate(text string, index int, src interface{}) (interface{}, error) {
	var t tmpl
	num, err := fmt.Sscanf(text, "%s %s %s %s %s %s %s", &t.key, &t.page, &t.head, &t.header, &t.nav, &t.main, &t.footer)
	if err != nil {
		return nil, fmt.Errorf("stk.Create:fmt.Sscanf failure num num:%d, text:%q.\n\t%v", num, text, err)
	} else if num < NUM_TMPL_PARAM {
		return nil, fmt.Errorf("stk.Create:number of item is too short. item num:%d, text:%q.\n\t%v", num, text, err)
	}
	t.key = strings.Trim(t.key, "\"")
	return &t, nil
}

func (mgr *TmplMap) GetTemplate(key string) (*template.Template, bool) {
	v, ok := mgr.tmpls[key]
	return v, ok
}

func Create(root, file, path, static string, funcMap template.FuncMap) (*TmplMap, error) {
	line, num, err := gim.CreateFileLines(root+file, nil, newTemplate)
	if err != nil {
		return nil, fmt.Errorf("wpg.Create:gim.CreateFileLines failure file:%q.\n\t%v", file, err)
	}

	mgr := &TmplMap{tmpls: make(map[string]*template.Template, num)}

	tmplPath, tmplStatic := root+path+"/", root+path+static
	var t *tmpl
	// fmt.Printf("tmplStatic:%q\n", tmplStatic)
	for line != nil {
		t = (line.Data).(*tmpl)
		// fmt.Printf("wpg.Create:key:%q line:%v\n", t.key, t)
		mgr.tmpls[t.key] = template.Must(template.New(tmplPath+t.page).Funcs(funcMap).ParseFiles(
			tmplPath+t.page, tmplPath+t.head, tmplPath+t.header,
			tmplPath+t.nav, tmplPath+t.main, tmplPath+t.footer, tmplStatic))
		line = line.Next
	}

	return mgr, nil
}

func EmbdCreate(files *embed.FS, root, file, path, static string, funcMap template.FuncMap) (*TmplMap, error) {
	log.Printf("Start wpg.EmbdCreate root:%q file:%q path:%q static:%q.", root, file, path, static)
	defer log.Printf("End wpg.EmbdCreate root:%q file:%q path:%q static:%q.", root, file, path, static)
	byte, err := files.ReadFile(root + "/" + file)
	if err != nil {
		return nil, fmt.Errorf("wpg.EmbdCreate:ReadFile failure root:%q path:%q file:%q.\n\t%v", root, path, file, err)
	}
	var s, key, page, head, header, nav, main, footer string
	var num int
	mgr := &TmplMap{tmpls: make(map[string]*template.Template)}
	tmplPath, tmplStatic := root+"/"+path, root+"/"+path+"/"+static
	scanner := bufio.NewScanner(strings.NewReader(string(byte)))
	for scanner.Scan() {
		s = scanner.Text()
		if len(s) > 0 { // 空白行（末尾など）なら参照しない
			if s[0] != '#' { // 行頭が#ならコメント行なので参照しない
				num, err = fmt.Sscanf(s, "%s %s %s %s %s %s %s", &key, &page, &head, &header, &nav, &main, &footer)
				if err != nil {
					return nil, fmt.Errorf("wpg.EmbdCreate:fmt.Sscanf failure num num:%d, s:%q.\n\t%v", num, s, err)
				} else if num < NUM_TMPL_PARAM {
					return nil, fmt.Errorf("wpg.EmbdCreate:number of item is too short. item num:%d, s:%q.\n\t%v", num, s, err)
				}
				key = strings.Trim(key, "\"")
				mgr.tmpls[key] = template.Must(template.New(tmplPath+page).Funcs(funcMap).ParseFiles(
					tmplPath+"/"+page, tmplPath+"/"+head, tmplPath+"/"+header,
					tmplPath+"/"+nav, tmplPath+"/"+main, tmplPath+"/"+footer, tmplStatic))
			}
		}
	}
	return mgr, nil
}

func TmplFunc(tmplFunc func(arg *TmplFuncArg) error, arg *TmplFuncArg) (string, error) {
	if arg.Elm.AtFirst {
		arg.Elm.AtFirst = false
		if err := tmplFunc(arg); err != nil {
			return "", fmt.Errorf("wpg.TmplFunc func failure err:%v", err)
		}
	}

	var str string
	switch arg.Index {
	case 1:
		str = arg.Elm.Param01
	case 2:
		str = arg.Elm.Param02
	case 3:
		str = arg.Elm.Param03
	case 4:
		str = arg.Elm.Param04
	case 5:
		str = arg.Elm.Param05
	case 6:
		str = arg.Elm.Param06
	case 7:
		str = arg.Elm.Param07
	case 8:
		str = arg.Elm.Param08
	case 9:
		str = arg.Elm.Param09
	case 10:
		str = arg.Elm.Param10
	case 11:
		str = arg.Elm.Param11
	case 12:
		str = arg.Elm.Param12
	case 13:
		str = arg.Elm.Param13
	case 14:
		str = arg.Elm.Param14
	case 15:
		str = arg.Elm.Param15
	case 16:
		str = arg.Elm.Param16
	case 17:
		str = arg.Elm.Param17
	case 18:
		str = arg.Elm.Param18
	case 19:
		str = arg.Elm.Param19
	case 20:
		str = arg.Elm.Param20
	case 21:
		str = arg.Elm.Param21
	case 22:
		str = arg.Elm.Param22
	case 23:
		str = arg.Elm.Param23
	case 24:
		str = arg.Elm.Param24
	case 25:
		str = arg.Elm.Param25
	case 26:
		str = arg.Elm.Param26
	case 27:
		str = arg.Elm.Param27
	case 28:
		str = arg.Elm.Param28
	case 29:
		str = arg.Elm.Param29
	case 30:
		str = arg.Elm.Param30
	default:
		return "", fmt.Errorf("wpg.TmplFunc Invalid arg.Index:%v arg:%#v", arg.Index, arg)
	}
	return str, nil
}

func (elm *TmplFuncArgElm) IncSet(p string) error {
	switch elm.Seq {
	case 1:
		elm.Param01 = p
	case 2:
		elm.Param02 = p
	case 3:
		elm.Param03 = p
	case 4:
		elm.Param04 = p
	case 5:
		elm.Param05 = p
	case 6:
		elm.Param06 = p
	case 7:
		elm.Param07 = p
	case 8:
		elm.Param08 = p
	case 9:
		elm.Param09 = p
	case 10:
		elm.Param10 = p
	case 11:
		elm.Param11 = p
	case 12:
		elm.Param12 = p
	case 13:
		elm.Param13 = p
	case 14:
		elm.Param14 = p
	case 15:
		elm.Param15 = p
	case 16:
		elm.Param16 = p
	case 17:
		elm.Param17 = p
	case 18:
		elm.Param18 = p
	case 19:
		elm.Param19 = p
	case 20:
		elm.Param20 = p
	case 21:
		elm.Param11 = p
	case 22:
		elm.Param12 = p
	case 23:
		elm.Param13 = p
	case 24:
		elm.Param14 = p
	case 25:
		elm.Param15 = p
	case 26:
		elm.Param16 = p
	case 27:
		elm.Param17 = p
	case 28:
		elm.Param18 = p
	case 29:
		elm.Param19 = p
	case 30:
		elm.Param20 = p
	default:
		return fmt.Errorf("wpg.TmplFuncArgElm.IncSet Invalid Seq:%v", elm.Seq)
	}
	elm.Seq++

	return nil
}

func (elm *TmplFuncArgElm) FormValueInt(nm string) (int, error) {
	v := elm.Req.FormValue(nm)
	if v == "" {
		return 0, nil
	}
	if n, err := strconv.Atoi(v); err != nil {
		return 0, fmt.Errorf("wpg.TmplFuncArgElm.FormValueInt:strconv.Atoi failure nm:%q v:%q.\n\t%v", nm, v, err)
	} else {
		return n, nil
	}
}

func (elm *TmplFuncArgElm) FormValueInt64(nm string) (int64, error) {
	v := elm.Req.FormValue(nm)
	if v == "" {
		return 0, nil
	}
	if n, err := strconv.ParseInt(v, 10, 64); err != nil {
		return 0, fmt.Errorf("wpg.TmplFuncArgElm.FormValueInt64:strconv.Atoi failure nm:%q v:%q.\n\t%v", nm, v, err)
	} else {
		return n, nil
	}
}
