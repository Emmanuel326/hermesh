package hedera

import (
    "fmt"

    hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"

    "github.com/Emmanuel326/hermesh/config"
)

type Client struct {
    inner     *hiero.Client
    AccountID hiero.AccountID
    TopicID   hiero.TopicID
}

func New(cfg *config.Config) (*Client, error) {
    accountID, err := hiero.AccountIDFromString(cfg.AccountID)
    if err != nil {
        return nil, fmt.Errorf("invalid account ID: %w", err)
    }

    privateKey, err := hiero.PrivateKeyFromString(cfg.PrivateKey)
    if err != nil {
        return nil, fmt.Errorf("invalid private key: %w", err)
    }

    var inner *hiero.Client
    switch cfg.Network {
    case "testnet":
        inner = hiero.ClientForTestnet()
    case "mainnet":
        inner = hiero.ClientForMainnet()
    default:
        inner = hiero.ClientForTestnet()
    }

    inner.SetOperator(accountID, privateKey)

    var topicID hiero.TopicID
    if cfg.TopicID != "" {
        topicID, err = hiero.TopicIDFromString(cfg.TopicID)
        if err != nil {
            return nil, fmt.Errorf("invalid topic ID: %w", err)
        }
    }

    return &Client{
        inner:     inner,
        AccountID: accountID,
        TopicID:   topicID,
    }, nil
}

func (c *Client) Inner() *hiero.Client {
    return c.inner
}

func (c *Client) Close() {
    c.inner.Close()
}
