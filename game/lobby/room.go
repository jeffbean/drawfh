package lobby

import (
	"github.com/google/uuid"
)

// The Room is a single instance of multiple players and a lobby
type Room struct {
	ID string `json:"id"` // unique id per game
	// Players

	// game state: created, start, choosing, guessing, finished

	// Drawing
	// current drawing... how to represent this...

	// Connections
	// Registered clients.
	clients map[*Player]struct{}

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Player

	// Unregister requests from clients.
	unregistered chan *Player

	done chan struct{}
}

// NewRoom creates a game room where players will join and leave a game.
// the room will remain running as long as at least one client is connected.
func NewRoom(creator string) (Room, error) {
	// TODO: should we do a url shortener to have these not expose the uuid we use as the ID?
	id, err := uuid.NewRandom()
	if err != nil {
		return Room{}, err
	}
	return Room{
		ID:           id.String(),
		clients:      make(map[*Player]struct{}),
		broadcast:    make(chan []byte),
		register:     make(chan *Player),
		unregistered: make(chan *Player),
		done:         make(chan struct{}, 1),
	}, nil
}

// PlayerJoin will join the player to the room.
func (r *Room) PlayerJoin(player *Player) {
	r.register <- player
}

// PlayerLeave will remove the player from the room.
func (r *Room) PlayerLeave(player *Player) {
	r.unregistered <- player
}

// Close will close the room and all tear down the websockets.
func (r *Room) Close() {
	for c := range r.clients {
		c.Stop()
	}
	close(r.broadcast)
	close(r.done)
}

// Start starts the room process. Each room will be its own loop.
func (r *Room) Start() {
	for {
		select {
		case <-r.done:
			return
		case c := <-r.register:
			r.clients[c] = struct{}{}
		case c := <-r.unregistered:
			if _, ok := r.clients[c]; ok {
				c.Stop()
				delete(r.clients, c)
			}
		case message := <-r.broadcast:
			for c := range r.clients {
				select {
				case c.send <- message:
				default:
					c.Stop()
					delete(r.clients, c)
				}
			}
		}
	}
}
