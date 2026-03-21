package main

import (
"flag"
"fmt"
"os"
"os/signal"
"syscall"
"time"

"github.com/Emmanuel326/hermesh/announce"
"github.com/Emmanuel326/hermesh/api"
"github.com/Emmanuel326/hermesh/cli"
"github.com/Emmanuel326/hermesh/config"
"github.com/Emmanuel326/hermesh/discover"
hederaclient "github.com/Emmanuel326/hermesh/hedera"
"github.com/Emmanuel326/hermesh/identity"
"github.com/Emmanuel326/hermesh/node"
"github.com/Emmanuel326/hermesh/peer"
)

func main() {
apiPort := flag.String("api-port", "9000", "HTTP API port")
pruneAfter := flag.Duration("prune", 10*time.Minute, "Prune dead nodes older than this duration")
flag.Parse()

// 1. Load config
cfg, err := config.Load()
if err != nil {
fmt.Println("❌ Config error:", err)
os.Exit(1)
}

// 2. Connect to Hedera
client, err := hederaclient.New(cfg)
if err != nil {
fmt.Println("❌ Hedera client error:", err)
os.Exit(1)
}
defer client.Close()

// 2b. Generate node identity
id, err := identity.New()
if err != nil {
fmt.Println("❌ Identity error:", err)
os.Exit(1)
}

// 3. Create this node's identity
self, err := node.New(cfg.NodeName, cfg.NodePort)
if err != nil {
fmt.Println("❌ Node error:", err)
os.Exit(1)
}

// 4. Wire up modules
store := peer.NewStore()
announcer := announce.New(client, id)
discoverer := discover.New(client, store.Handle)
terminal := cli.New(store)
apiServer := api.New(store, *apiPort)

// 5. Print banner and status
terminal.PrintBanner()
terminal.PrintStatus(self, cfg.TopicID)
fmt.Printf("  KeyPrint : %s\n\n", id.FingerPrint())

// 6. Start API server in background
go func() {
if err := apiServer.Start(); err != nil {
fmt.Println("❌ API server error:", err)
}
}()

// 7. Start discovery
terminal.PrintEvent("Subscribing to network topic...")
if err := discoverer.Start(); err != nil {
fmt.Println("❌ Discovery error:", err)
os.Exit(1)
}
terminal.PrintEvent("Listening for peers...")

// 8. Announce this node
time.Sleep(2 * time.Second)
terminal.PrintEvent("Announcing node to network...")
if err := announcer.Announce(self); err != nil {
fmt.Println("❌ Announce error:", err)
os.Exit(1)
}
terminal.PrintEvent(fmt.Sprintf("Node announced: %s", self))

// 9. Tickers
heartbeatTicker := time.NewTicker(30 * time.Second)
defer heartbeatTicker.Stop()

gcTicker := time.NewTicker(30 * time.Second)
defer gcTicker.Stop()

pruneTicker := time.NewTicker(1 * time.Minute)
defer pruneTicker.Stop()

peersTicker := time.NewTicker(10 * time.Second)
defer peersTicker.Stop()

// 10. Graceful shutdown
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

// 11. Main loop
for {
select {
case <-heartbeatTicker.C:
if err := announcer.Heartbeat(self); err != nil {
terminal.PrintEvent(fmt.Sprintf("❌ Heartbeat failed: %s", err))
} else {
terminal.PrintEvent("💓 Heartbeat sent")
}

case <-gcTicker.C:
store.GarbageCollect()

case <-pruneTicker.C:
n := store.Prune(*pruneAfter)
if n > 0 {
terminal.PrintEvent(fmt.Sprintf("🧹 Pruned %d dead nodes", n))
}

case <-peersTicker.C:
terminal.PrintPeers()

case <-quit:
terminal.PrintEvent("Shutting down gracefully...")
announcer.Leave(self)
terminal.PrintEvent("👋 Left the network. Goodbye.")
os.Exit(0)
}
}
}
