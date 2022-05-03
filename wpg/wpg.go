package wpg

import (
	"strings"
	"net/url"
)

const (
	URL_PATH_DEPTH0 = iota
	URL_PATH_DEPTH1
	URL_PATH_DEPTH2
)

const	(
	HTML_TAG_BR = "<br>"
	HTML_DISABLED = "disabled"
	HTML_CHECKED = "checked=\"checked\""
	HTML_CHECKED_VALUE = "1"
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
