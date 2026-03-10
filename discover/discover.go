package discover

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Emmanuel326/hermesh/announce"
	hederaclient "github.com/Emmanuel326/hermesh/hedera"
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
			var announcement announce.Message
			if err := json.Unmarshal(msg.Contents, &announcement); err != nil {
				// skip legacy or malformed messages
				return
			}
			d.handler(&announcement)
		})

	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	return nil
}
