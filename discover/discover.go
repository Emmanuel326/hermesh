package discover

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Emmanuel326/hermesh/announce"
	hederaclient "github.com/Emmanuel326/hermesh/hedera"
	"github.com/Emmanuel326/hermesh/identity"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

type Handler func(msg *announce.Message)

type Discoverer struct {
	client  *hederaclient.Client
	topicID hiero.TopicID
	handler Handler
}

func New(client *hederaclient.Client, handler Handler) *Discoverer {
	return &Discoverer{
		client:  client,
		topicID: client.TopicID,
		handler: handler,
	}
}

func (d *Discoverer) Start() error {
	_, err := hiero.NewTopicMessageQuery().
		SetTopicID(d.topicID).
		SetStartTime(time.Unix(0, 0)).
		Subscribe(d.client.Inner(), func(msg hiero.TopicMessage) {

			// Try to unmarshal as SignedEnvelope first
			var envelope identity.SignedEnvelope
			if err := json.Unmarshal(msg.Contents, &envelope); err != nil {
				// skip legacy or malformed messages
				return
			}

			// Skip unsigned messages — no public key
			if len(envelope.PublicKey) == 0 {
				return
			}

			// Verify signature — drop tampered messages silently
			payload, err := identity.Verify(&envelope)
			if err != nil {
				fmt.Printf("⚠️  Dropped unverified message from %s\n", envelope.NodeID[:8])
				return
			}

			// Unmarshal the verified payload
			var announcement announce.Message
			if err := json.Unmarshal(payload, &announcement); err != nil {
				return
			}

			d.handler(&announcement)
		})

	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	return nil
}
