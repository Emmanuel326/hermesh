package config

import (
    "fmt"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    AccountID  string
    PrivateKey string
    Network    string
    TopicID    string
    NodeName   string
    NodePort   string
}

func Load() (*Config, error) {
    godotenv.Load()

    cfg := &Config{
        AccountID:  os.Getenv("HEDERA_ACCOUNT_ID"),
        PrivateKey: os.Getenv("HEDERA_PRIVATE_KEY"),
        Network:    os.Getenv("HEDERA_NETWORK"),
        TopicID:    os.Getenv("HERMESH_TOPIC_ID"),
        NodeName:   os.Getenv("HERMESH_NODE_NAME"),
        NodePort:   os.Getenv("HERMESH_NODE_PORT"),
    }

    if err := cfg.validate(); err != nil {
        return nil, err
    }

    return cfg, nil
}

func (c *Config) validate() error {
    if c.AccountID == "" {
        return fmt.Errorf("HEDERA_ACCOUNT_ID is required")
    }
    if c.PrivateKey == "" {
        return fmt.Errorf("HEDERA_PRIVATE_KEY is required")
    }
    if c.NodeName == "" {
        return fmt.Errorf("HERMESH_NODE_NAME is required")
    }
    return nil
}
