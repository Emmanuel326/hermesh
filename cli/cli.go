package cli

import (
	"fmt"
        "encoding/json"
	"os"
	"text/tabwriter"
	"time"

	"github.com/Emmanuel326/hermesh/node"
	"github.com/Emmanuel326/hermesh/peer"
)

type CLI struct {
	store *peer.Store
}

func New(store *peer.Store) *CLI {
	return &CLI{store: store}
}

func (c *CLI) PrintBanner() {
	fmt.Println(`
██╗  ██╗███████╗██████╗ ███╗   ███╗███████╗███████╗██╗  ██╗
██║  ██║██╔════╝██╔══██╗████╗ ████║██╔════╝██╔════╝██║  ██║
███████║█████╗  ██████╔╝██╔████╔██║█████╗  ███████╗███████║
██╔══██║██╔══╝  ██╔══██╗██║╚██╔╝██║██╔══╝  ╚════██║██╔══██║
██║  ██║███████╗██║  ██║██║ ╚═╝ ██║███████╗███████║██║  ██║
╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝╚══════╝╚═╝  ╚═╝
	`)
	fmt.Println("  Decentralized service mesh coordination over Hedera HCS")
	fmt.Println("  Version 0.1.0")
	fmt.Println()
}

func (c *CLI) PrintStatus(self *node.Node, topicID string) {
	fmt.Printf("  Node     : %s\n", self.Name)
	fmt.Printf("  ID       : %s\n", self.ID[:8])
	fmt.Printf("  Address  : %s:%s\n", self.IP, self.Port)
	fmt.Printf("  Topic    : %s\n", topicID)
	fmt.Printf("  Network  : Hedera Testnet\n")
	fmt.Println()
}

func (c *CLI) PrintPeers() {
	peers := c.store.List()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "  NAME\tID\tADDRESS\tSTATUS\tLAST SEEN")
	fmt.Fprintln(w, "  ----\t--\t-------\t------\t---------")

	for _, n := range peers {
		fmt.Fprintf(w, "  %s\t%s\t%s:%s\t%s\t%s\n",
			n.Name,
			n.ID[:8],
			n.IP,
			n.Port,
			colorStatus(n.Status),
			humanize(n.LastSeen),
		)
	}
	w.Flush()

	fmt.Printf("\n  Total: %d peers (%d alive)\n\n",
		c.store.Count(),
		len(c.store.Alive()),
	)
}

func (c *CLI) PrintEvent(event string) {
	fmt.Printf("  [%s] %s\n", time.Now().Format("15:04:05"), event)
}

func colorStatus(s node.Status) string {
	switch s {
	case node.StatusAlive:
		return "🟢 alive"
	case node.StatusSuspected:
		return "🟡 suspected"
	case node.StatusDead:
		return "🔴 dead"
	default:
		return "❓ unknown"
	}
}

func humanize(t time.Time) string {
	diff := time.Since(t)
	switch {
	case diff < time.Minute:
		return fmt.Sprintf("%ds ago", int(diff.Seconds()))
	case diff < time.Hour:
		return fmt.Sprintf("%dm ago", int(diff.Minutes()))
	default:
		return fmt.Sprintf("%dh ago", int(diff.Hours()))
	}
}

func (c *CLI) Query(service string, statusFilter string) {
peers := c.store.List()
result := make([]*node.Node, 0)

for _, n := range peers {
if service != "" && n.Name != service {
continue
}
if statusFilter != "" && string(n.Status) != statusFilter {
continue
}
result = append(result, n)
}

out, _ := json.Marshal(map[string]any{
"peers": result,
"total": len(result),
})
fmt.Println(string(out))
}
