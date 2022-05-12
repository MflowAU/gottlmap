package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MFlowAU/gottlmap"
)

var f *os.File

func example_with_prehook() {
	// open file handler to txt file
	// this is used later in PreDeleteHook
	file := "./example.txt"
	f, _ = os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()

	t := time.NewTicker(1 * time.Second) // How often the cleanup routine is called
	ctx, cancel := context.WithCancel(context.Background())
	ttl_map, err := gottlmap.New(t, PreDeleteHook, ctx)
	if err != nil {
		log.Fatal(err)
	}
	ttl_map.Set("key1", "value1", 2*time.Second)                                                     // expire the key in 2 seconds
	log.Printf("Current Keys, Value, Expiry in TTLMap %s, %+v \n", ttl_map.Keys(), ttl_map.Values()) // Key1: value1, Expires: 2s
	log.Println("Sleeping for 5 seconds...")
	time.Sleep(3 * time.Second)
	log.Printf("Current Keys in TTLMap %+v\n", ttl_map.Keys()) // No Entries
	cancel()
	time.Sleep(2 * time.Second)
}

func PreDeleteHook(k string, v gottlmap.Element) error {
	log.Println("PreHook invoked...")
	log.Printf("Key: %s, Value: %v\n", k, v.Value)
	log.Println("Persisting to a file to before deleting from memory...")
	f.WriteString(fmt.Sprintf("Key: %s, Value: %v\n", k, v.Value))
	return nil
}
