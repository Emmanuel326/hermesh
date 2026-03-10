package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Emmanuel326/hermesh/announce"
	"github.com/Emmanuel326/hermesh/cli"
	"github.com/Emmanuel326/hermesh/config"
	hederaclient "github.com/Emmanuel326/hermesh/hedera"
	"github.com/Emmanuel326/hermesh/discover"
	"github.com/Emmanuel326/hermesh/node"
	"github.com/Emmanuel326/hermesh/peer"
)

func main() {
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

	// 3. Create this node's identity
	self, err := node.New(cfg.NodeName, cfg.NodePort)
	if err != nil {
		fmt.Println("❌ Node error:", err)
		os.Exit(1)
	}

	// 4. Wire up modules
	store := peer.NewStore()
	announcer := announce.New(client)
	discoverer := discover.New(client, store.Handle)
	terminal := cli.New(store)

	// 5. Print banner and status
	terminal.PrintBanner()
	terminal.PrintStatus(self, cfg.TopicID)

	// 6. Start discovery — replay full history first
	terminal.PrintEvent("Subscribing to network topic...")
	if err := discoverer.Start(); err != nil {
		fmt.Println("❌ Discovery error:", err)
		os.Exit(1)
	}
	terminal.PrintEvent("Listening for peers...")

	// 7. Announce this node to the network
	time.Sleep(2 * time.Second) // let history replay first
	terminal.PrintEvent("Announcing node to network...")
	if err := announcer.Announce(self); err != nil {
		fmt.Println("❌ Announce error:", err)
		os.Exit(1)
	}
	terminal.PrintEvent(fmt.Sprintf("Node announced: %s", self))

	// 8. Heartbeat every 30 seconds
	heartbeatTicker := time.NewTicker(30 * time.Second)
	defer heartbeatTicker.Stop()

	// 9. Garbage collect every 30 seconds
	gcTicker := time.NewTicker(30 * time.Second)
	defer gcTicker.Stop()

	// 10. Print peers every 10 seconds
	peersTicker := time.NewTicker(10 * time.Second)
	defer peersTicker.Stop()

	// 11. Handle shutdown gracefully
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 12. Main loop
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
