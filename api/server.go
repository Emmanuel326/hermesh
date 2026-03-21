package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Emmanuel326/hermesh/node"
	"github.com/Emmanuel326/hermesh/peer"
)

type Server struct {
	store *peer.Store
	port  string
}

func New(store *peer.Store, port string) *Server {
	return &Server{store: store, port: port}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/peers", s.handlePeers)
	mux.HandleFunc("/health", s.handleHealth)

	fmt.Printf("  [API] Listening on http://localhost:%s\n", s.port)
	return http.ListenAndServe(":"+s.port, mux)
}

func (s *Server) handlePeers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	statusFilter := r.URL.Query().Get("status")
	service := r.URL.Query().Get("service")

	peers := s.store.List()
	result := make([]*node.Node, 0)

	for _, n := range peers {
		if statusFilter != "" && string(n.Status) != statusFilter {
			continue
		}
		if service != "" && n.Name != service {
			continue
		}
		result = append(result, n)
	}

	json.NewEncoder(w).Encode(map[string]any{
		"peers": result,
		"total": len(result),
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	alive := s.store.Alive()
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"peers":  s.store.Count(),
		"alive":  len(alive),
	})
}
