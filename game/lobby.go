package game

import  (
	"net/http"
	"log"

	"github.com/gorilla/websocket"
)

const _clientSendBufferLen = 512
var _upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// The Lobby is a single instance of multiple players and a lobby
type Lobby struct {
	ID string // unique id per game lobby
	// Players

	// game state: created, start, choosing, guessing, finished
	//

	// Drawing
	// current drawing... how to represent this...

	// Connections
	// Registered clients.
	clients map[*Client]struct{}

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// Client is the glue for a player in a lobby.
type Client struct {
	Lobby *Lobby
	Player Player

	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

type Player struct {
	ID string // generated
	NickName string
}

func NewLobby(owner Player) Lobby {
	return Lobby{
		Creator: owner,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]struct{}),
	}
}

func (l *Lobby) Start()  {
	for {
		select {
		case client := <-l.register:
			l.clients[client] = struct{}{}
		case client := <-l.unregister:
			if _, ok := l.clients[client]; ok {
				delete(l.clients, client)
				close(client.send)
			}
		case message := <-l.broadcast:
			for client := range l.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(l.clients, client)
				}
			}
		}
	}
}

func  (l *Lobby) serveWS( w http.ResponseWriter, r *http.Request) error {
	conn, err := _upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	client := &Client{Lobby: l, conn: conn, send: make(chan []byte, _clientSendBufferLen)}
	l.registerClient(client)
	return nil
}

func (l *Lobby) registerClient(client *Client) error {
	return nil
}
