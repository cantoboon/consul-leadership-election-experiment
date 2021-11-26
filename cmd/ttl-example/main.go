package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"time"
)

// You'll see that for each session we may, or may not, get the lock. The TTL on the session is
// 10s. Once the TTL ends, a new session can then acquire the lock. In real life, you would
// call session.Renew(...) to keep the session alive.
func main() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalln("Failed to create Consul client", err)
	}

	sessionClient := client.Session()
	kvClient := client.KV()

	i := 0
	for {
		i++
		sessionId, _, err := sessionClient.Create(
			&api.SessionEntry{
				Name: fmt.Sprintf("test-sessionClient-%d", i),
				TTL:  "10s",
			},
			nil,
		)
		if err != nil {
			log.Fatalln("Failed to create sessionClient", err)
		}

		fmt.Println("New sessionClient created. sessionClient id = ", sessionId)

		kvPair := &api.KVPair{Session: sessionId, Key: "my-key", Value: []byte("my-value")}

		res, _, err := kvClient.Acquire(kvPair, nil)
		if err != nil {
			log.Fatalln("Failed to acquire lock")
		}

		fmt.Println("Lock acquired: ", res)

		time.Sleep(5 * time.Second)
	}
}
