package main

import (
	"context"
	"log"
	"time"

	"github.com/MFlowAU/gottlmap"
)

func example() {
	t := time.NewTicker(1 * time.Second) // How often the cleanup routine is called
	ctx, cancel := context.WithCancel(context.Background())
	ttl_map, err := gottlmap.New(t, nil, ctx)
	if err != nil {
		log.Fatal(err)
	}
	ttl_map.Set("key1", "value1", 2*time.Second)                                                     // expire the key in 2 seconds
	log.Printf("Current Keys, Value, Expiry in TTLMap %s, %+v \n", ttl_map.Keys(), ttl_map.Values()) // Key1: value1, Expires: 2s
	log.Println("Sleeping for 3 seconds...")
	time.Sleep(3 * time.Second)
	// The TTLMap should not have any entries after 3 seconds
	log.Printf("Current Keys in TTLMap %+v\n", ttl_map.Keys()) // No Entries
	cancel()                                                   // stop the ttlmap routine
	time.Sleep(2 * time.Second)
}
