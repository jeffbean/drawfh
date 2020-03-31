package lobby

import "github.com/jeffbean/drawfh/game/client"

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
