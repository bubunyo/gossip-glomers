package main

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	g := NewIdGenerator()
	n.Handle("generate", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		ns := strings.Replace(msg.Dest, "n", "", 1)
		nn, err := strconv.Atoi(ns)
		if err != nil {
			return err
		}
		// Update the message type to return back.
		body["type"] = "generate_ok"
		body["id"] = g.GenerateId(nn)

		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
