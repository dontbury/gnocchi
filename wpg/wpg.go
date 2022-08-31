package wpg

import (
	"fmt"
	"strings"
//	"strconv"
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

// HTML ALIGN
const HTML_ALIGN = "align"
const	(
	HTML_ALIGN_NONE = 0
	HTML_ALIGN_TOP = 1
	HTML_ALIGN_BOTTOM = 2
	HTML_ALIGN_LEFT = 3
	HTML_ALIGN_RIGHT = 4
)

type SlctOpt struct {
	Val string
	Cap string
	Slct bool
	Dsbl bool
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

func ( s *HtmlTh )Str() string {
	return "<th>" + s.Hd + "</th>"
}

func ( s *HtmlTd )Str() string {
	str := "<td"
	if len( s.Cls ) > 0 {
		str += " class=\"" + s.Cls + "\""
	}
	str += ">"
	return str + s.Dt + "</td>"
}

func StrHtmlTblHdrs( hdrs *[]HtmlTh ) string {
	str := "<tr>"
	var v HtmlTh
	for _, v = range *hdrs {
		str += v.Str()
	}
	return str + "</tr>"
}

func ( row *HtmlTr )Str() string {
	str := "<tr>"
	if row.Hd != nil {
		str += row.Hd.Str()
	}
	var v HtmlTd
	for _, v = range row.Row {
		str += v.Str()
	}
	return str + "</tr>"
}

func ( cpt *HtmlCpt )Str() string {
	str := "<caption"
	switch ( cpt.Align ) {
		case HTML_ALIGN_NONE:	// なにもしない
		case HTML_ALIGN_TOP: str += " " + HTML_ALIGN + "=\"top\""
		case HTML_ALIGN_BOTTOM: str += " " + HTML_ALIGN + "=\"bottom\""
		case HTML_ALIGN_LEFT: str += " " + HTML_ALIGN + "=\"left\""
		case HTML_ALIGN_RIGHT: str += " " + HTML_ALIGN + "=\"right\""
		default: fmt.Errorf("wpg.HtmlCpt.Str:Invalid align:%d.", cpt.Align )
	}
	return  str + ">" + cpt.Title + "</caption>"
}

func ( tbl *HtmlTbl )Str() string {
	str := "<table"
	if len( tbl.Cls ) > 0 {
		str += " class=\"" + tbl.Cls + "\""
	}
	str += ">"
	if tbl.Cpt != nil {
		str += tbl.Cpt.Str()
	}
	str += StrHtmlTblHdrs( &( tbl.Hdrs ) )
	var v HtmlTr
	for _, v = range tbl.Bdy {
		str += v.Str()
	}
	return str + "</table>"
}
