package lobby

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoomCreation(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	room, err := NewRoom("test-creator")
	require.NoError(t, err)

	socket, _, err := websocket.DefaultDialer.Dial("ws://"+s.Listener.Addr().String()+"/ws", nil)
	require.NoError(t, err)
	defer socket.Close()

	assert.NotEmpty(t, room.ID)
	go room.Start()

	testPlayer := NewPlayer("test-player-one")
	testPlayer.Start(&room, socket)

	testPlayer.send <- []byte("{}")

	room.Close()
}

// helpers
var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}
