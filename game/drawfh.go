package game

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/jeffbean/drawfh/game/lobby"
)

type DrawGame struct {
	mu sync.RWMutex

	// Room ID -> room
	rooms map[string]*lobby.Room // protected by mu
}

// NewServer creates our main game server that manages the game
// rooms.
func NewGameServer() DrawGame {
	return DrawGame{}
}

func (d *DrawGame) HomeHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world v2!")
}

func (d *DrawGame) RoomCreate(w http.ResponseWriter, r *http.Request) {
	// get the player who created this room
	// make the room.
	room, err := lobby.NewRoom("todo")
	if err != nil {
		http.Error(w, "failed to create new room", http.StatusInternalServerError)
		return
	}
	d.mu.Lock()
	d.rooms[room.ID] = &room
	d.mu.Unlock()

	go room.Run()

	// return http created and send url to room page.
}

func (d *DrawGame) RoomDelete(w http.ResponseWriter, r *http.Request) {
	roomID := "todo"

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.rooms[roomID]; ok {
		delete(d.rooms, roomID)
	}
}

func (d *DrawGame) RoomList(w http.ResponseWriter, r *http.Request) {
	d.mu.RLock()
	p := struct {
		Rooms map[string]*lobby.Room `json:"rooms"`
	}{
		Rooms: d.rooms,
	}
	d.mu.RUnlock()

	json.NewEncoder(w).Encode(p)
}
