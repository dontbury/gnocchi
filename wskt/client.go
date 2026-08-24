package wskt

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"gnocchi/bitbyte"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

type Client interface {
	ReceiveWebsocket(id int, wsbuf *WSBuf) error
	CreateSendCliBuf(buf *[]byte) (*[]byte, error)

	OnClose(id int)
	GetHTMLText() string
}

type ChClient struct {
	id int

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	cli	Client
}

// The application runs receiveWebSocket in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *ChClient) receiveWebSocket() {
	log.Printf("Start wskt.receiveWebSocket.")
	defer func() {
		log.Printf("wskt.receiveWebSocket close id:%d.", c.id)
		//		c.conn.Close()		// この2つはOnCloseの呼び出し先のRemoveClientでリストから外した後にcloseされる
		//		close( c.send )
		c.cli.OnClose(c.id)
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for loop := true; loop; {
		if messageType, buf, err := c.conn.ReadMessage(); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				var ce *websocket.CloseError
				if errors.As(err, &ce) {
					fmt.Printf("wskt.ChClient.receiveWebSocket errorClose Code:%d, Text:%q\n", ce.Code, ce.Text)
				}
				// log.Printf("error: %v", err)	// これはクライアントが切断したときに発生するエラーなのでログに出す必要はない
			} else {
				log.Printf("wskt.receiveWebSocket c.Conn.ReadMessage() id:%d.\n\t%+v", c.id, err)
			}
			loop = false // ループから抜ける
		} else {
			switch messageType {
			case websocket.TextMessage:
				log.Printf("wskt.ChClient.receiveWebSocket TextMessage:%q", string(buf))
			case websocket.BinaryMessage:
				sz := len(buf)
				if sz > 0 {
					wsBuf := WSBuf{Br: bitbyte.BitRow{Index: 0, Inc: 0, Body: make([]uint64, (sz+bitbyte.BYTES_PER_VALUE-1)/bitbyte.BITS_PER_BYTE)}}
					for i, v := range buf {
						wsBuf.Br.Body[i>>3] |= uint64(v) << (bitbyte.BITS_PER_BYTE * (i & (bitbyte.BYTES_PER_VALUE - 1)))
					}
					log.Printf("wskt.receiveWebSocket c.conn.ReadMessage() buf:[%s].", bitbyte.StrBytes(&buf))
					log.Printf("wskt.receiveWebSocket c.conn.ReadMessage() wsBuf.Br.Body:[%s].", wsBuf.Br.StrBody())
					if err = c.cli.ReceiveWebsocket(c.id, &wsBuf); err != nil {
						log.Printf("wskt.receiveWebSocket.\n\t%+v", err)
						//				loop = false	// ループから抜ける
					}
				} else {
					log.Printf("wskt.receiveWebSocket c.Conn.ReadMessage() Empty Message id:%d.", c.id)
					loop = false // ループから抜ける
				}
			case websocket.CloseMessage:
				log.Printf("wskt.ChClient.receiveWebSocket CloseMessage buf:%v", buf)
				loop = false // 通信が閉じられた場合はerrとして扱われ、ここには来ないらしい
			case websocket.PingMessage:
				log.Printf("wskt.ChClient.receiveWebSocket PinMessage buf:%v", buf)
			case websocket.PongMessage:
				log.Printf("wskt.ChClient.receiveWebSocket PonMessage buf:%v", buf)
			default:
				log.Printf("wskt.receiveWebSocket c.Conn.ReadMessage() Invalid MessageType:%d id:%d.", messageType, c.id)
				loop = false // ループから抜ける
			}
		}
	}
	log.Printf("End wskt.receiveWebSocket.")
}

// A goroutine running receiveChannel is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *ChClient) receiveChannel() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case buf, ok := <-c.send:
			// log.Printf("wskt.client receiveChannel c.send buf:%+v.", buf)
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}

			if err = c.sendClient(&w, &buf); err != nil {
				log.Printf("wskt.ChClient.receiveChannel failed.\n\t%v", err)
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				buf = <-c.send
				if err = c.sendClient(&w, &buf); err != nil {
					log.Printf("wskt.ChClient.receiveChannel failed.\n\t%v", err)
					return
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			log.Printf("wskt.ChClient.receiveChannel ticker.C.")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *ChClient) sendClient(w *io.WriteCloser, buf *[]byte) error {
	// log.Printf("wskt.ChClient.sendClient:buf:%+v.", buf)
	if c.cli != nil { // Server.RemoveClientが呼ばれるとnilがセットされるので
		// 途中ずっとポインタで受け渡しをして最後にバイト列で送信
		if send, err := c.cli.CreateSendCliBuf(buf); err == nil {
			// log.Printf("wskt.ChClient.sendClient:size:%d send:%+v.", len(*send), send)
			(*w).Write(*send)
		} else {
			return fmt.Errorf("wskt.ChClient.sendClient:Client.CreateSendCliBuf failed.\n\t%v", err)
		}
	}
	return nil
}

func (c *ChClient) SendChannel(bytes *[]byte) {
	c.send <- *bytes
}

func (c *ChClient) GetHTMLText() string {
	return fmt.Sprintf("<tr><td scope=\"row\">%d</td>%s</tr>", c.id, c.cli.GetHTMLText())
}
