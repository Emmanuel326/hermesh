package announce

import (
    "encoding/json"
    "fmt"
    "time"

    hederaclient "github.com/Emmanuel326/hermesh/hedera"
    "github.com/Emmanuel326/hermesh/node"
    hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

type Message struct {
    Type      string      `json:"type"`
    Node      *node.Node  `json:"node"`
    Timestamp time.Time   `json:"timestamp"`
}

type Announcer struct {
    client  *hederaclient.Client
    topicID hiero.TopicID
}

func New(client *hederaclient.Client) *Announcer {
    return &Announcer{
        client:  client,
        topicID: client.TopicID,
    }
}

func (a *Announcer) Announce(n *node.Node) error {
    msg := Message{
        Type:      "join",
        Node:      n,
        Timestamp: time.Now().UTC(),
    }
    return a.publish(msg)
}

func (a *Announcer) Heartbeat(n *node.Node) error {
    n.MarkAlive()
    msg := Message{
        Type:      "heartbeat",
        Node:      n,
        Timestamp: time.Now().UTC(),
    }
    return a.publish(msg)
}

func (a *Announcer) Leave(n *node.Node) error {
    n.MarkDead()
    msg := Message{
        Type:      "leave",
        Node:      n,
        Timestamp: time.Now().UTC(),
    }
    return a.publish(msg)
}

func (a *Announcer) publish(msg Message) error {
    payload, err := json.Marshal(msg)
    if err != nil {
        return fmt.Errorf("failed to marshal message: %w", err)
    }

    _, err = hiero.NewTopicMessageSubmitTransaction().
        SetTopicID(a.topicID).
        SetMessage(payload).
        Execute(a.client.Inner())
    if err != nil {
        return fmt.Errorf("failed to publish to HCS: %w", err)
    }

    return nil
}
