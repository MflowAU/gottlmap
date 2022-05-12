# gottlmap


gottlmap is a [Go](http://golang.org) package that provides an in-memory
key-value store for TTL-based expirable items.

This package also allows for userdefined action to be performed before deleting the key-value
pair from memory. Ex: Persist to DB, file or write to a network socket.


## Install


```go
go get -u github.com/MFlowAU/gottlmap
```


## Examples

```go
package main

import (
	"context"
	"log"
	"os"
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

```
For more examples check out the example folder

## Contribute

Please fork this project and send me a pull request!
