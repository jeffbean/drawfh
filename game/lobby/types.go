package lobby

import (
	"github.com/google/uuid"
	"github.com/jeffbean/drawfh/game/client"
)

// The Room is a single instance of multiple players and a lobby
type Room struct {
	ID string // unique id per game lobby
	// Players

	// game state: created, start, choosing, guessing, finished
	//

	// Drawing
	// current drawing... how to represent this...

	// Connections
	// Registered clients.
	clients map[*client.Player]struct{}

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *client.Player

	// Unregister requests from clients.
	unregister chan *client.Player
}

func NewRoom(creator string) (Room, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Room{}, err
	}
	return Room{
		ID:         id.String(),
		broadcast:  make(chan []byte),
		register:   make(chan *client.Player),
		unregister: make(chan *client.Player),
	}, nil
}

// PlayerJoin will join the player to the room.
func (r *Room) PlayerJoin(player *client.Player) {
	r.register <- player
}

// PlayerLeave will remove the player from the room.
func (r *Room) PlayerLeave(player *client.Player) {
	r.unregister <- player
}

// Close will close the room and all tear down the websockets.
func (r *Room) Close() {
	for client := range r.clients {
		client.Stop()
	}
	close(r.broadcast)
}

// Run starts the room process. Each room will be its own loop.
func (r *Room) Run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = struct{}{}
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				client.Stop()
			}
		case message := <-r.broadcast:
			for client := range r.clients {
				client.Send([]byte(message))
			}
		}
	}
}
