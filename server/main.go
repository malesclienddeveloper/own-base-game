package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type PlayerState struct {
	Resources int  `json:"resources"`
	Builders  int  `json:"builders"`
	Assigned  bool `json:"assigned"`
}

type BaseState struct {
	ID      string                 `json:"id"`
	OwnerID string                 `json:"ownerId"`
	THLevel int                    `json:"thLevel"`
	Players map[string]*PlayerState `json:"players"`
	mu      sync.Mutex
}

type RoomManager struct {
	rooms map[string]*BaseState
	mu    sync.Mutex
}

func NewRoomManager() *RoomManager { return &RoomManager{rooms: map[string]*BaseState{}} }

func (rm *RoomManager) GetOrCreate(id, owner string) *BaseState {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if r, ok := rm.rooms[id]; ok {
		return r
	}
	r := &BaseState{ID: id, OwnerID: owner, THLevel: 1, Players: map[string]*PlayerState{}}
	rm.rooms[id] = r
	return r
}

func (b *BaseState) AddPlayer(pid string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.Players) >= 4 {
		return false
	}
	if _, ok := b.Players[pid]; ok {
		return true
	}
	b.Players[pid] = &PlayerState{Resources: 1000, Builders: 2, Assigned: false}
	return true
}

func (b *BaseState) Assign(pid string, ok bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if p, exists := b.Players[pid]; exists {
		p.Assigned = ok
	}
}

func (b *BaseState) UpgradeTH() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.THLevel < 16 { // max placeholder
		b.THLevel++
	}
}

func handleWS(rm *RoomManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer conn.Close()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var cmd map[string]any
			if err := json.Unmarshal(msg, &cmd); err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"bad json"}`))
				continue
			}
			// Minimal command handling placeholder
			resp, _ := json.Marshal(map[string]any{"received": cmd})
			conn.WriteMessage(websocket.TextMessage, resp)
		}
	}
}

func main() {
	rm := NewRoomManager()
	_ = rm
	http.HandleFunc("/ws", handleWS(rm))
	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
