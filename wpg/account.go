package wpg

import (
	"fmt"
	"log"
	"time"
	"net/http"

	"github.com/satori/go.uuid"
)

const	COOKIE_NAME = "sid"

type Paths []string

type Account struct {
	ID		int64
	Name	string
}

type Conn struct {
    Sid   string
    Acc   *Account
    Date  time.Time
}

type Accounts struct {
	size int
	conns map[string]*Conn
}

func ( accs *Accounts ) Create( size int ) {
	accs.size = size
	accs.conns = make( map[ string ]*Conn, size )
}

func ( accs *Accounts ) Connection( w http.ResponseWriter, r *http.Request ) *Conn {
	var sid string
	cookie, cerr := r.Cookie( COOKIE_NAME )
	if cerr != nil {
		log.Printf( "wgp.Accounts.Connection r.Cookie cerr:%#v.", cerr )	// 取得できなかったら作るので、抜けない
		u := uuid.NewV4()
/*	返り値の数が変更になったらしい 2022.05.03
		u, err := uuid.NewV4()
		if err != nil {
			log.Printf( "wgp.Accounts.Connection uuid.NewV4() failure u:%v, err:%v.", u, err )
			return nil
		}
*/
		sid = u.String()
		c := accs.conns[ sid ]
		if c != nil {
			log.Printf( "wgp.Accounts.Connection uuid.NewV4() sid:%q, already existed connection:%v.", sid, c )
			return nil
		}
		v := &http.Cookie{ Name:COOKIE_NAME, Value:sid, Path:"/" }
		log.Printf( "wgp.Accounts.Connection new cookie:%v.", v )
		http.SetCookie( w, v )
	} else if cookie != nil {
		log.Printf( "wgp.Accounts.Connection get cookie:%v.", cookie )
		con := accs.conns[ cookie.Value ]
		if con != nil { // すでにコネクションリストには存在している場合は、最終アクセス時間を更新して抜ける
			log.Printf( "wgp.Accounts.Connection match cookie:%v.", cookie )
			con.Date = time.Now()
			return con
		} else {
			log.Printf( "wgp.Accounts.Connection old but not exist cookie:%v.", cookie )
			sid = cookie.Value	// すでに存在していた場合は再利用する
		}
	}


	if len( accs.conns ) >= accs.size {// 接続数が上限に達していれば、最も古いものを一つだけリストから外す。
		log.Printf( "conns size:%d.", len( accs.conns ) )
		var lastCon *Conn = nil
		for _, c := range accs.conns {
			if lastCon == nil {
				lastCon = c
			} else if lastCon.Date.After( c.Date ) {
				lastCon = c
			}
		}
		if lastCon != nil {
			delete( accs.conns, lastCon.Sid )
			log.Printf( "lastCon:%v ejected conns size:%d.", lastCon, len( accs.conns ) )
		}
	}

	con := &Conn{ Sid:sid, Acc:nil, Date:time.Now() }
	log.Printf( "Conn.Connection:Set con:%+v into conns sid;%q.", con, sid )
	accs.conns[ sid ] = con

	return con
}

func ( accs *Accounts ) GetConn( sid string ) *Conn {
  return accs.conns[ sid ]
}

func ( accs *Accounts ) GetConnListHTML() string {
	buf := ""
	for _, v := range accs.conns {
		if v.Acc != nil {
			buf += fmt.Sprintf("<tr><td scope=\"row\">%s</td><td>%d</td><td>%s</td><td>%v</td></tr>", v.Sid, v.Acc.ID, v.Acc.Name, v.Date )
		} else {
			buf += fmt.Sprintf("<tr><td scope=\"row\">%s</td><td>nil</td><td>nil</td><td>%v</td></tr>", v.Sid, v.Date.Format( time.ANSIC ) )
		}
	}
	return buf
}

func LinkAddress(link string) string {
  return "\"/" + link + "\""
}

func LinkHTML( link, id, class, caption string, enable bool ) string {
	strID, strClass, strDisabled := "", "", ""
	if id != "" {
		strID = " id=\"" + id + "\""
	}
	if class != "" {
		strClass = " class=\"" + class + "\""
	}
	if !enable {
		strDisabled = " disabled"
	}
	return "<a href=" + LinkAddress( link ) + strID + strClass + strDisabled + ">" + caption + "</a>"
}

func ParagraphHTML(id, caption string) string {
  return "<p id=\"" + id + "\">" + caption + "</p>"
}

func AccountHTML(con *Conn) string {
	login := true
	if con == nil {
		login = false
	} else if con.Acc == nil {
		login = false
	}
	if login {
		return ParagraphHTML("account_name", "User:" + con.Acc.Name)
	}
	return ParagraphHTML("guest", "You are currently not logged in.")
}
