package wskt

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/dontbury/gnocchi/bitbyte"
)

const BUFFER_SIZE = 1024

type CmdType byte

const (
	CmdNone CmdType = iota
	CmdJoin
	CmdDefect

	CmdCmnMax = 0x80
)

// ChClientへの通知を意味する構造体
type Cmd struct {
	cmd   CmdType
	bytes []byte
}

type Server interface {
	ReceiveChannel(buf *[]byte) error
}

type ChServer struct {
	clientID int
	clients  map[int]*ChClient

	// Inbound messages from the clients.
	broadcast chan []byte
	svr       Server
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  BUFFER_SIZE,
	WriteBufferSize: BUFFER_SIZE,
}

func StartServerProcess(svr Server, num int) *ChServer {
	s := &ChServer{
		clientID:  0,
		clients:   make(map[int]*ChClient, num),
		broadcast: make(chan []byte),
		svr:       svr,
	}
	go s.run()
	return s
}

func (s *ChServer) run() {
	for {
		buf := <-s.broadcast
		if err := s.svr.ReceiveChannel(&buf); err != nil {
			log.Printf("wskt.ChServer.run:Server.ReceiveChannel failure buf:[%s].\n\t%v", bitbyte.StrBytes(&buf), err)
		}
	}
}

func (s *ChServer) RegistClient(cli Client, w http.ResponseWriter, r *http.Request) *ChClient {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	s.clientID++
	client := &ChClient{
		id:   s.clientID,
		conn: ws,
		send: make(chan []byte, 256),
		cli:  cli,
	}

	s.clients[client.id] = client

	go client.receiveChannel()
	go client.receiveWebSocket()

	return client
}

func (s *ChServer) RemoveClient(id int) error {
	if c, ok := s.clients[id]; c != nil && ok {
		c.cli = nil
		delete(s.clients, id)
		c.conn.Close()
		close(c.send)
		log.Printf("wskt.ChServer.RemoveClient:id:%d c:%+v.", id, c)
	} else {
		return fmt.Errorf("wskt.ChServer.RemoveClient:Can't get Client id:%d, s.clients:%+v.", id, s.clients)
	}
	return nil
}

func (s *ChServer) CollectClients(f func(Client) (bool, error), send *WSBuf) (*[]*ChClient, error) {
	clientes := make([]*ChClient, 0, len(s.clients))
	var valid bool
	var err error
	bytes := (*[]byte)(nil)
	if send != nil {
		if bytes, err = send.GetSendBuf(); err != nil {
			return nil, fmt.Errorf("wskt.ChServer.CollectClients:WSBuf.GetSendBuf failure send:%v.\n\t%v", send, err)
		}
	}
	for _, c := range s.clients {
		if valid, err = f(c.cli); err != nil {
			return nil, fmt.Errorf("wskt.ChServer.CollectClients:Error occurred while checking client validity c:%v.\n\t%v", c, err)
		} else if valid {
			clientes = append(clientes, c)
			if bytes != nil {
				c.send <- *bytes
			}
		}
	}
	return &clientes, nil
}

func (s *ChServer) SendServerBroadcast(wsBuf *WSBuf) {
	// 途中ずっとポインタで受け渡しをして最後にバイト列で受信
	if buf, err := wsBuf.GetSendBuf(); err == nil {
		s.broadcast <- *buf
	} else {
		log.Printf("wskt.ChServer.SendServerBroadcast:WSBuf.GetSendBuf failure wsBuf:%v.\n\t%v", wsBuf, err)
	}
}

func (s *ChServer) SendWSBCli(clientID int, buf *[]byte) error {
	if c, ok := s.clients[clientID]; c != nil && ok {
		// 途中ずっとポインタで受け渡しをして最後にバイト列で送信
		c.send <- *buf
		/*	なぜ上のように書かないのかわからない。ChatGPTは下のようにコメントを補完したが、別にそんなことはないと思う。
			なぜ上のように途中ずっとポインタで受け渡しをして最後にバイト列で送信するのかというと、WSBuf.GetSendBuf()の中で、後ろの未使用部分を切り詰めたバッファを取得するために、WSBuf内部のBr.Bodyからビット単位で値を取得しているためです。
			もしWSBuf.GetSendBuf()の中で、Br.Bodyからビット単位で値を取得せずに、Br.Bodyをそのままバイト列として送信してしまうと、Br.Bodyの後ろの未使用部分も含めて送信されてしまい、クライアント側で正しくデータを受信できなくなってしまいます。
			そのため、WSBuf.GetSendBuf()の中で、Br.Bodyからビット単位で値を取得して、後ろの未使用部分を切り詰めたバッファを取得し、そのバッファを途中ずっとポインタで受け渡しをして最後にバイト列で送信することで、クライアント側で正しくデータを受信できるようにしています。
					send := buf
					c.send <- *send
		*/
	} else {
		return fmt.Errorf("wskt.ChServer.SendWSBCli:Can't get Client clientID:%d, s.clients:%+v.", clientID, s.clients)
	}
	return nil
}

func (s *ChServer) PrintHTMLClientList() string {
	buf := ""
	for _, v := range s.clients {
		buf += v.GetHTMLText()
	}
	return buf
}
