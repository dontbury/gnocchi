package wpg

import (
	"fmt"
	"strings"
	"strconv"
	"net/url"
)

const (
	URL_PATH_DEPTH0 = iota
	URL_PATH_DEPTH1
	URL_PATH_DEPTH2
)

const	(
	HTML_TAG_BR = "<br>"
	HTML_TAG_HR = "<hr>"
	HTML_DISABLED = "disabled"
	HTML_CHECKED = "checked=\"checked\""
	HTML_CHECKED_VALUE = "1"
	HTML_SPACE = "&nbsp;"
)

const	(
	TFAM_INDEX01 = iota + 1
	TFAM_INDEX02
	TFAM_INDEX03
	TFAM_INDEX04
	TFAM_INDEX05
	TFAM_INDEX06
	TFAM_INDEX07
	TFAM_INDEX08
	TFAM_INDEX09
	TFAM_INDEX10
	TFAM_INDEX11
	TFAM_INDEX12
	TFAM_INDEX13
	TFAM_INDEX14
	TFAM_INDEX15
	TFAM_INDEX16
	TFAM_INDEX17
	TFAM_INDEX18
	TFAM_INDEX19
	TFAM_INDEX20
)

// アカウント種別
const	(
	ACC_TYPE_ADMIN = iota + 1
	ACC_TYPE_02
	ACC_TYPE_03
	ACC_TYPE_04
	ACC_TYPE_05
	ACC_TYPE_06
	ACC_TYPE_07
	ACC_TYPE_08
	ACC_TYPE_09
	ACC_TYPE_10
	ACC_TYPE_11
	ACC_TYPE_12
	ACC_TYPE_13
	ACC_TYPE_14
	ACC_TYPE_15
	ACC_TYPE_USER
)

// HTML
// HTML ALIGN
const HTML_ALIGN = "align"
const	(
	HTML_ALIGN_NONE = 0
	HTML_ALIGN_TOP = 1
	HTML_ALIGN_BOTTOM = 2
	HTML_ALIGN_LEFT = 3
	HTML_ALIGN_RIGHT = 4
)
// HTML FORM METHOD
const (
	HTML_FORM_METHOD_NONE = 0
	HTML_FORM_METHOD_GET = 1
	HTML_FORM_METHOD_POST = 2
)
// HTML INPUT TYPE
const (
	HTML_INPUT_TYPE_NONE = 0
	HTML_INPUT_TYPE_TEXT = 1
	HTML_INPUT_TYPE_PASSWORD = 2
	HTML_INPUT_TYPE_RADIO = 3
	HTML_INPUT_TYPE_CHECKBOX = 4
	HTML_INPUT_TYPE_FILE = 5
	HTML_INPUT_TYPE_HIDDEN = 6
	HTML_INPUT_TYPE_SUBMIT = 7
	HTML_INPUT_TYPE_IMAGE = 8
	HTML_INPUT_TYPE_RESET = 9
	HTML_INPUT_TYPE_BUTTON = 10
)

type HtmlTag interface {
	Str() ( string, error )
}

type SlctOpt struct {
	Val string
	Cap string
	Slct bool
	Dsbl bool
}

type PlainText struct {
	Text string
}

type HtmlNbsp struct {
}

type HtmlBr struct {
}

type HtmlHr struct {
}

type HtmlForm struct {
	Action string
	Method int
	Id string
	Body []HtmlTag
}

type HtmlLabel struct {
	For string
	Body []HtmlTag
}

type HtmlAbbr struct {
	Title string
	Body string
}

type HtmlInput struct {
	Type int
	Name string
	Placeholder string
	Size int
	Minlength int
	Maxlength int
	Cls string
	Id string
	Value string
	Required bool
}

type HtmlTh struct {
	Hd string
}

type HtmlTd struct {
	Dt string
	Cls string
}

type HtmlTr struct {
	Hd *HtmlTh
	Row []HtmlTd
}

type HtmlCpt struct {
	Title string
	Align int
}

type HtmlTbl struct {
	Cls string
	Cpt *HtmlCpt
	Hdrs []HtmlTh
	Bdy []HtmlTr
}

type HtmlA struct {
	Path string
	Args url.Values
	Cpt string
	Id string
	Cls string
	Disable bool
}

func EscapeHTML( src string ) string {
	str := strings.Replace(src, "&", "&amp;", -1)
	str = strings.Replace(str, "<", "&lt;", -1)
	str = strings.Replace(str, ">", "&gt;", -1)
	str = strings.Replace(str, "\"", "&quot;", -1)
	return str
}

// HTMLフォーマットのパラグラフを出力
func PrintHTMLParagraph( p string ) string {
	return "<p>" + p + "</p>"
}

// HTMLフォーマットのヘッダー3を出力
func PrintHTMLH3( p string ) string {
	return "<h3>" + p + "</h3>"
}

// HTMLフォーマットのarticleを出力
func PrintHTMLArticle( class, p string ) string {
	str := "<article"
	if len( class ) > 0 {
		str += " class=\"" + class + "\""
	}
	str += ">" + p + "</article>"
	return str
}

// HTMLフォーマット:textタイプのinputを出力
func PrintHTMLInputText( name, value, size, maxsize string ) string {
	str := "<input type=\"text\" name=\"" + name + "\" value=\"" + value + "\""
	if len( size ) > 0 {
		str += " size=\"" + size + "\""
	}
	if len( maxsize ) > 0 {
		str += " maxsize=\"" + maxsize + "\""
	}
	str += ">"
	return str
}

// HTMLフォーマット:numberタイプのinputを出力
func PrintHTMLInputNumber( name, value, size, maxsize, min, max string ) string {
	str := "<input type=\"number\" name=\"" + name + "\""
	if len( value ) > 0 {
		str += " value=\"" + value + "\""
	}
	if len( size ) > 0 {
		str += " size=\"" + size + "\""
	}
	if len( maxsize ) > 0 {
		str += " maxsize=\"" + maxsize + "\""
	}
	if len( min ) > 0 {
		str += " min=\"" + min + "\""
	}
	if len( max ) > 0 {
		str += " max=\"" + max + "\""
	}
	str += ">"
	return str
}

// HTMLフォーマット:submitタイプのinputを出力
func PrintHTMLInputSubmit( value, cls string ) string {
	str := "<input type=\"submit\" value=\"" + value + "\""
	if len( cls ) > 0 { str += " class=\"" + cls + "\"" }
	return str + ">"
}

// HTMLフォーマット:hiddenタイプのinputを出力
func PrintHTMLInputHidden( name, value, cls string ) string {
	str := "<input type=\"hidden\" name=\"" + name + "\" value=\"" + value + "\""
	if len( cls ) > 0 { str += " class=\"" + cls + "\"" }
	return str + ">"
}

// HTMLフォーマット:postメソッドのformを出力
func PrintHTMLForm( action, content string ) string {
	return "<form method=\"GET\" action=\"" + action + "\">" + content + "<form>"
}

// HTMLフォーマット:selectを出力
func PrintHTMLSelect( name string, arOpt *[]SlctOpt ) string {
	str := "<select name=\"" + name + "\">"
	for _, v := range *arOpt {
		str += "<option value=\"" + v.Val + "\""
		if v.Slct { str+= " selected" }
		if v.Dsbl { str+= " disabled" }
		str += ">" + v.Cap + "</option>"
	}
	str += "</select>"
	return str
}

// 引数をMapで渡してリンクを作る
func LinkArgsMapHTML( path, caption string, mapArg url.Values ) string {
	first := true
	arg := ""
	for key, v := range mapArg {
		for _, s := range v {
			if first {
				arg += "?"
				first = false
			} else {
				arg += "&"
			}
			arg += key + "=" + s
		}
	}
	return LinkHTML( path + arg, "", "", caption, true )
}

func ( s *PlainText )Str() ( string, error ) {
	return s.Text, nil
}

func ( s *HtmlNbsp )Str() ( string, error ) {
	return "&nbsp", nil
}

func ( s *HtmlBr )Str() ( string, error ) {
	return "<br>", nil
}

func ( s *HtmlHr )Str() ( string, error ) {
	return "<hr>", nil
}

func ( s *HtmlForm )Str() ( string, error ) {
	str := "<form method=\""
	switch( s.Method ) {
		case HTML_FORM_METHOD_GET:
			str += "GET\""
		case HTML_FORM_METHOD_POST:
			str += "POST\""
		default:
			return "", fmt.Errorf( "wpg.HtmlForm.Str:Invalid Method:%d.", s.Method )
	}
	str += " action=\""+ s.Action + "\""
	if len( s.Id ) > 0 { str += " id=\"" + s.Id + "\"" }
	str += ">\n"
	var v HtmlTag
	var buf string
	var i int
	var err error
	for i, v = range s.Body {
		if v == nil { return "", fmt.Errorf( "wpg.HtmlForm.Str:v(%d) is nil.", i ) }
		if buf, err = v.Str(); err != nil { return "", fmt.Errorf( "wpg.HtmlForm.Str:wpg.HtmlTag.Str failure:%v.\t\n", err ) }
		str += buf + "\n"
	}
	return str + "</form>", nil
}

func ( s *HtmlLabel )Str() ( string, error ) {
	str, buf := "", ""
	var v HtmlTag
	var i int
	var err error
	for i, v = range s.Body {
		if v == nil { return "", fmt.Errorf( "wpg.HtmlLabel.Str:v(%d) is nil.", i ) }
		if buf, err = v.Str(); err != nil { return "", fmt.Errorf( "wpg.HtmlLabel.Str:wpg.HtmlTag.Str failure:%v.\t\n", err ) }
		str += buf
	}
	return  "<label for=\"" + s.For + "\">" + str + "</label>", nil
}

func ( s *HtmlAbbr )Str() ( string, error ) {
	return  "<abbr title=\"" + s.Title + "\">" + s.Body + "</abbr>", nil
}

func ( s *HtmlInput )Str() ( string, error ) {
	str := "type=\""
	switch( s.Type ) {
		case HTML_INPUT_TYPE_TEXT:
			str += "text"
		case HTML_INPUT_TYPE_PASSWORD:
			str += "password"
		case HTML_INPUT_TYPE_RADIO:
			str += "radio"
		case HTML_INPUT_TYPE_CHECKBOX:
			str += "checkbox"
		case HTML_INPUT_TYPE_FILE:
			str += "file"
		case HTML_INPUT_TYPE_HIDDEN:
			str += "hidden"
		case HTML_INPUT_TYPE_SUBMIT:
			str += "submit"
		case HTML_INPUT_TYPE_IMAGE:
			str += "image"
		case HTML_INPUT_TYPE_RESET:
			str += "reset"
		default:
			return "", fmt.Errorf( "wpg.HtmlInput.Str:Invalid Type:%d.", s.Type )
	}
	str += "\" name=\"" + s.Name + "\""
	if len( s.Placeholder ) > 0 { str += " placeholder=\"" + s.Placeholder + "\"" }
	if s.Size > 0 { str += " size=\"" + strconv.Itoa( s.Size ) + "\"" }
	if s.Minlength > 0 { str += " minlength=\"" + strconv.Itoa( s.Minlength ) + "\"" }
	if s.Maxlength > 0 { str += " maxlength=\"" + strconv.Itoa( s.Maxlength ) + "\"" }
	if len( s.Cls ) > 0 { str += " class=\"" + s.Cls + "\"" }
	if len( s.Id ) > 0 { str += " id=\"" + s.Id + "\"" }
	if len( s.Value ) > 0 { str += " value=\"" + s.Value + "\"" }
	if s.Required { str += " required" }
	return "<input " + str + ">", nil
}

func ( s *HtmlTh )Str() ( string, error ) {
	return  "<th>" + s.Hd + "</td>", nil
}

func ( s *HtmlTd )Str() ( string, error ) {
	str :=""
	if len( s.Cls ) > 0 {
		str += " class=\"" + s.Cls + "\""
	}
	return  "<td" + str + ">" + s.Dt + "</td>", nil
}

func ( row *HtmlTr )Str() ( string, error ) {
	str, buf := "", ""
	var err error
	if row.Hd != nil {
		if buf, err = row.Hd.Str(); err != nil { return "", fmt.Errorf( "wpg.HtmlTr.Str:wpg.HtmlTh.Str failure %v.\t\n", err ) }
		str += buf
	}
	var v HtmlTd
	for _, v = range row.Row {
		if buf, err = v.Str(); err != nil { return "", fmt.Errorf( "wpg.HtmlTr.Str:wpg.HtmlTd.Str failure %v.\t\n", err ) }
		str += buf
	}
	return "<tr>" + str + "</tr>", nil
}

func ( cpt *HtmlCpt )Str() ( string, error ) {
	str := ""
	switch ( cpt.Align ) {
		case HTML_ALIGN_NONE:	// なにもしない
		case HTML_ALIGN_TOP: str = " " + HTML_ALIGN + "=\"top\""
		case HTML_ALIGN_BOTTOM: str = " " + HTML_ALIGN + "=\"bottom\""
		case HTML_ALIGN_LEFT: str = " " + HTML_ALIGN + "=\"left\""
		case HTML_ALIGN_RIGHT: str = " " + HTML_ALIGN + "=\"right\""
		default: return "", fmt.Errorf("wpg.HtmlCpt.Str:Invalid align:%d.", cpt.Align )
	}
	return  "<caption" + str + ">" + cpt.Title + "</caption>", nil
}

func ( tbl *HtmlTbl )Str() ( string, error ) {
	str, buf := "<table", ""
	var err error
	if len( tbl.Cls ) > 0 {
		str += " class=\"" + tbl.Cls + "\""
	}
	str += ">"

	if tbl.Cpt != nil {
		if buf, err = tbl.Cpt.Str(); err != nil { return "", fmt.Errorf( "wpg.HtmlTbl.Str:wpg.HtmlCpt.Str failure %v.\t\n", err ) }
		str += buf
	}
	str += buf

	var h HtmlTh
	for _, h = range tbl.Hdrs {
		if buf, err = h.Str(); err != nil { return "", fmt.Errorf( "wpg.HtmlTbl.Str:wpg.HtmlTh.Str failure %v.\t\n", err ) }
		str += buf
	}

	var v HtmlTr
	for _, v = range tbl.Bdy {
		if buf, err = v.Str(); err != nil { return "", fmt.Errorf( "wpg.HtmlTbl.Str:wpg.HtmlTr.Str failure %v.\t\n", err ) }
		str += buf
	}
	return str + "</table>", nil
}

func ( a *HtmlA )Str() ( string, error ) {
	str := ""
	first := true
	for key, v := range a.Args {
		for _, s := range v {
			if first {
				str += "?"
				first = false
			} else {
				str += "&"
			}
			str += key + "=" + s
		}
	}
	str += "\""

	if len( a.Id  ) > 0 {
		str += " id=\"" + a.Id + "\""
	}
	if len( a.Cls  ) > 0 {
		str += " class=\"" + a.Cls + "\""
	}
	if a.Disable {
		str += " disabled"
	}

	return "<a href=\"/" + a.Path + str + ">" + a.Cpt + "</a>", nil
}
