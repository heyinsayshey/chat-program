package Controller

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

const socketBufferSize = 1024
const messageBufferSize = 256

var (
	upgrader = &websocket.Upgrader{
		ReadBufferSize:  socketBufferSize,
		WriteBufferSize: socketBufferSize,
	}
)

var clients []*Client

type Client struct {
	conn   *websocket.Conn // 웹소켓 커넥션
	send   chan *Message   // 메시지 전송용 채널
	roomId string          // 현재 접속한 채팅방 아이디
	user   *User           // 현재 접속한 사용자 정보
}

func CreateClient(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if nil != err {
		log.Println(err)
		return
	}

	newClient(socket, p.ByName("room_id"), GetCurrentUser(req))
}

func newClient(conn *websocket.Conn, roomId string, u *User) {

	c := &Client{
		conn:   conn,
		send:   make(chan *Message, messageBufferSize),
		roomId: roomId,
		user:   u,
	}

	clients = append(clients, c)

	go c.readLoop()
	go c.writeLoop()
}

func (c *Client) Close() {
	for i, client := range clients {
		if client == c {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	close(c.send)

	c.conn.Close()
	log.Printf("close connection. addr : %s", c.conn.RemoteAddr())
}

func (c *Client) readLoop() {
	for {
		m, err := c.read()
		if nil != err {
			log.Println(err)
			break
		}

		m.CreateMessage()
		broadcast(m)
	}
	c.Close()
}

func (c *Client) writeLoop() {
	for msg := range c.send {
		if c.roomId == msg.Room_id {
			c.write(msg)
		}
	}
}

func broadcast(m *Message) {
	for _, client := range clients {
		client.send <- m
	}
}

func (c *Client) read() (*Message, error) {
	var msg *Message

	if err := c.conn.ReadJSON(&msg); err != nil {
		return nil, err
	}

	now := time.Now()
	strTime := now.Format(time.RFC3339)
	idx := strings.Index(strTime, "+")
	strTime = strTime[:idx]
	strTime = strTime + ".000Z"

	msg.Send_date = strTime
	msg.User = *c.user

	log.Println("read from websocket : ", msg)

	return msg, nil
}

func (c *Client) write(m *Message) error {
	log.Println("write to websocket : ", m)

	return c.conn.WriteJSON(m)
}
