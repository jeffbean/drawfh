package game

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jeffbean/drawfh/game/lobby"
)

const (
	_socketBufferSize  = 1024
	_messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  _socketBufferSize,
	WriteBufferSize: _messageBufferSize,
}

type DrawGame struct {
	mu sync.RWMutex

	// Room ID -> room
	rooms map[string]*lobby.Room // protected by mu
}

// NewServer creates our main game server that manages the game
// rooms.
func NewGameServer() DrawGame {
	return DrawGame{
		rooms: make(map[string]*lobby.Room),
	}
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

	go room.Start()

	log.Printf("")
	// return http created and send url to room page.
	io.WriteString(w, room.ID)
}

func (d *DrawGame) RoomDetail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomID, ok := params["id"]
	if !ok {
		http.Error(w, "id is required variable", http.StatusBadRequest)
		return
	}

	room, ok := d.findRoom(roomID)
	if !ok {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(room); err != nil {
		http.Error(w, "encoding room json", http.StatusInternalServerError)
	}
}

func (d *DrawGame) RoomJoin(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	roomID, ok := params["id"]
	if !ok {
		http.Error(w, "id is required variable", http.StatusBadRequest)
		return
	}

	room, ok := d.findRoom(roomID)
	if !ok {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}

	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "connection upgrade failed", http.StatusInternalServerError)
		return
	}

	currentUser := lobby.NewPlayer("testing-user")
	currentUser.Start(room, socket)
	room.PlayerJoin(currentUser)
	io.WriteString(w, fmt.Sprintf("player %v joined room %v", currentUser, room.ID))
}

func (d *DrawGame) RoomDelete(w http.ResponseWriter, r *http.Request) {
	roomID := "todo"

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.rooms[roomID]; ok {
		delete(d.rooms, roomID)
	}
}

func (d *DrawGame) RoomList(w http.ResponseWriter, _ *http.Request) {
	d.mu.RLock()
	p := struct {
		Rooms map[string]*lobby.Room `json:"rooms"`
	}{
		Rooms: d.rooms,
	}
	d.mu.RUnlock()

	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, "encoding room list json", http.StatusInternalServerError)
	}
}

func (d *DrawGame) findRoom(roomID string) (*lobby.Room, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	r, ok := d.rooms[roomID]
	return r, ok
}
