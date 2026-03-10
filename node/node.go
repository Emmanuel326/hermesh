package node

import (
    "fmt"
    "net"
    "os"
    "time"

    "github.com/google/uuid"
)

type Status string

const (
    StatusAlive     Status = "alive"
    StatusSuspected Status = "suspected"
    StatusDead      Status = "dead"
)

type Node struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    IP        string    `json:"ip"`
    Port      string    `json:"port"`
    Status    Status    `json:"status"`
    Version   string    `json:"version"`
    StartedAt time.Time `json:"started_at"`
    LastSeen  time.Time `json:"last_seen"`
}

func New(name, port string) (*Node, error) {
    ip, err := getOutboundIP()
    if err != nil {
        // fallback to hostname
        ip, _ = os.Hostname()
    }

    return &Node{
        ID:        uuid.New().String(),
        Name:      name,
        IP:        ip,
        Port:      port,
        Status:    StatusAlive,
        Version:   "0.1.0",
        StartedAt: time.Now().UTC(),
        LastSeen:  time.Now().UTC(),
    }, nil
}

func (n *Node) String() string {
    return fmt.Sprintf("[%s] %s @ %s:%s (%s)",
        n.ID[:8], n.Name, n.IP, n.Port, n.Status)
}

func (n *Node) IsAlive() bool {
    return n.Status == StatusAlive
}

func (n *Node) MarkSuspected() {
    n.Status = StatusSuspected
}

func (n *Node) MarkDead() {
    n.Status = StatusDead
}

func (n *Node) MarkAlive() {
    n.Status = StatusAlive
    n.LastSeen = time.Now().UTC()
}

func getOutboundIP() (string, error) {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        return "", err
    }
    defer conn.Close()
    localAddr := conn.LocalAddr().(*net.UDPAddr)
    return localAddr.IP.String(), nil
}
