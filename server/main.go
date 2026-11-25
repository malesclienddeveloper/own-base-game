package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

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

func (b *BaseState) UpgradeTH() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.THLevel < 16 {
		b.THLevel++
		return true
	}
	return false
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func main() {
	rm := NewRoomManager()

	http.HandleFunc("/api/room/create", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		owner := r.URL.Query().Get("owner")
		if id == "" || owner == "" {
			http.Error(w, "missing id/owner", 400)
			return
		}
		base := rm.GetOrCreate(id, owner)
		base.AddPlayer(owner)
		base.Assign(owner, true)
		writeJSON(w, base)
	})

	http.HandleFunc("/api/room/join", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		pid := r.URL.Query().Get("player")
		base, ok := func() (*BaseState, bool) {
			rm.mu.Lock()
			defer rm.mu.Unlock()
			b, ok := rm.rooms[id]
			return b, ok
		}()
		if !ok {
			http.Error(w, "room not found", 404)
			return
		}
		if !base.AddPlayer(pid) {
			http.Error(w, "room full (max 4)", 400)
			return
		}
		writeJSON(w, base)
	})

	http.HandleFunc("/api/room/upgrade_th", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		rm.mu.Lock()
		base, ok := rm.rooms[id]
		rm.mu.Unlock()
		if !ok {
			http.Error(w, "room not found", 404)
			return
		}
		base.UpgradeTH()
		writeJSON(w, base)
	})

	http.Handle("/", http.FileServer(http.Dir("../web")))

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
