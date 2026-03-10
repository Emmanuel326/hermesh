package identity

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
)

// Identity holds this node's keypair for the lifetime of the process
type Identity struct {
	PublicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

// SignedEnvelope wraps any message with a cryptographic signature
type SignedEnvelope struct {
	Payload   []byte `json:"payload"`    // raw JSON of Message
	Signature []byte `json:"signature"`  // ED25519 signature of Payload
	PublicKey []byte `json:"public_key"` // verifier's key
	NodeID    string `json:"node_id"`    // quick lookup
}

// New generates a fresh ED25519 keypair
// Called once on startup — new identity every run
func New() (*Identity, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate keypair: %w", err)
	}

	return &Identity{
		PublicKey:  pub,
		privateKey: priv,
	}, nil
}

// Sign wraps a message in a SignedEnvelope
func (id *Identity) Sign(nodeID string, v any) (*SignedEnvelope, error) {
	payload, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	signature := ed25519.Sign(id.privateKey, payload)

	return &SignedEnvelope{
		Payload:   payload,
		Signature: signature,
		PublicKey: id.PublicKey,
		NodeID:    nodeID,
	}, nil
}

// Verify checks the envelope signature
// Returns the raw payload if valid, error if tampered or invalid
func Verify(env *SignedEnvelope) ([]byte, error) {
	pub := ed25519.PublicKey(env.PublicKey)

	if !ed25519.Verify(pub, env.Payload, env.Signature) {
		return nil, fmt.Errorf("invalid signature from node %s", env.NodeID)
	}

	return env.Payload, nil
}

// FingerPrint returns a short human readable public key identifier
func (id *Identity) FingerPrint() string {
	return fmt.Sprintf("%x", id.PublicKey[:8])
}
