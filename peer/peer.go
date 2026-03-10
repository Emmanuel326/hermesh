package peer

import (
    "fmt"
    "sync"
    "time"

    "github.com/Emmanuel326/hermesh/announce"
    "github.com/Emmanuel326/hermesh/node"
)

const (
    SuspectAfter = 90 * time.Second
    DeadAfter    = 180 * time.Second
)

type Store struct {
    mu    sync.RWMutex
    peers map[string]*node.Node
}

func NewStore() *Store {
    return &Store{
        peers: make(map[string]*node.Node),
    }
}

func (s *Store) Handle(msg *announce.Message) {
    if msg.Node == nil {
        return
    }

    s.mu.Lock()
    defer s.mu.Unlock()

    switch msg.Type {
    case "join":
        msg.Node.LastSeen = msg.Timestamp
        msg.Node.Status = node.StatusAlive
        s.peers[msg.Node.ID] = msg.Node
        fmt.Printf("🟢 Node joined:     %s\n", msg.Node)

    case "heartbeat":
        if existing, ok := s.peers[msg.Node.ID]; ok {
            existing.LastSeen = msg.Timestamp
            existing.MarkAlive()
            fmt.Printf("💓 Heartbeat:       %s\n", existing)
        }

    case "leave":
        if existing, ok := s.peers[msg.Node.ID]; ok {
            existing.MarkDead()
            fmt.Printf("🔴 Node left:       %s\n", existing)
        }
    }
}

func (s *Store) GarbageCollect() {
    s.mu.Lock()
    defer s.mu.Unlock()

    now := time.Now().UTC()
    for _, n := range s.peers {
        age := now.Sub(n.LastSeen)
        if age > DeadAfter {
            n.MarkDead()
        } else if age > SuspectAfter {
            n.MarkSuspected()
        }
    }
}

func (s *Store) List() []*node.Node {
    s.mu.RLock()
    defer s.mu.RUnlock()

    nodes := make([]*node.Node, 0, len(s.peers))
    for _, n := range s.peers {
        nodes = append(nodes, n)
    }
    return nodes
}

func (s *Store) Count() int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return len(s.peers)
}

func (s *Store) Alive() []*node.Node {
    s.mu.RLock()
    defer s.mu.RUnlock()

    nodes := make([]*node.Node, 0)
    for _, n := range s.peers {
        if n.IsAlive() {
            nodes = append(nodes, n)
        }
    }
    return nodes
}
