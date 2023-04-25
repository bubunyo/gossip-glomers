package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	log.SetOutput(os.Stderr)
	// Update the message type to return back.
	g := NewIdGenerator()
	n.Handle("generate", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		nid := n.ID()
		ns := strings.ReplaceAll(nid, "n", "")
		if nid == "" {
			log.Fatalf("node id cannot be empty. NodeId=%s, nids=%v", n.ID(), n.NodeIDs())
		}
		nn, err := strconv.Atoi(ns)
		if err != nil {
			log.Fatal(err)
		}

		body["type"] = "generate_ok"
		body["id"] = g.GenerateId(nn)

		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
