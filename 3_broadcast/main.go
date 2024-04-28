package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type D map[string]any

func main() {
	b := newStore[float64]()
	t := newStore[string]()
	n := maelstrom.NewNode()

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		m := body["message"].(float64)
		b.save(m)

		// check if this is a propagation message
		if body["propagation"] == nil {
			propMsg := D{
				"type":        "broadcast",
				"message":     m,
				"propagation": true,
			}
			l := t.read()
			for _, tt := range l {
				if n.ID() != tt {
					n.RPC(tt, propMsg, func(msg maelstrom.Message) error {
						return nil
					})
				}
			}
		} else {
			delete(body, "propagation")
		}
		delete(body, "message")
		body["type"] = "broadcast_ok"
		// wg.Wait()
		return n.Reply(msg, body)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		body["type"] = "read_ok"
		body["messages"] = b.read()
		return n.Reply(msg, body)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		body["type"] = "topology_ok"
		t.save(n.NodeIDs()...)

		delete(body, "topology")
		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
