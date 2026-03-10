package announce

import (
	"encoding/json"
	"fmt"
	"time"

	hederaclient "github.com/Emmanuel326/hermesh/hedera"
	"github.com/Emmanuel326/hermesh/identity"
	"github.com/Emmanuel326/hermesh/node"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

type Message struct {
	Type      string     `json:"type"`
	Node      *node.Node `json:"node"`
	Timestamp time.Time  `json:"timestamp"`
}

type Announcer struct {
	client   *hederaclient.Client
	topicID  hiero.TopicID
	identity *identity.Identity
}

func New(client *hederaclient.Client, id *identity.Identity) *Announcer {
	return &Announcer{
		client:   client,
		topicID:  client.TopicID,
		identity: id,
	}
}

func (a *Announcer) Announce(n *node.Node) error {
	return a.publish("join", n)
}

func (a *Announcer) Heartbeat(n *node.Node) error {
	n.MarkAlive()
	return a.publish("heartbeat", n)
}

func (a *Announcer) Leave(n *node.Node) error {
	n.MarkDead()
	return a.publish("leave", n)
}

func (a *Announcer) publish(msgType string, n *node.Node) error {
	msg := Message{
		Type:      msgType,
		Node:      n,
		Timestamp: time.Now().UTC(),
	}

	// Sign the message
	envelope, err := a.identity.Sign(n.ID, msg)
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}

	// Marshal the signed envelope
	payload, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("failed to marshal envelope: %w", err)
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
