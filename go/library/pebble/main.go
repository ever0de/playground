package main

import (
	"fmt"
	"log"

	"github.com/cockroachdb/pebble"
)

func main() {
	db, err := pebble.Open("./tmp", &pebble.Options{
		Levels: []pebble.LevelOptions{
			{TargetFileSize: 1 << 20},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	key := []byte("hello")
	_, closer, err := db.Get(key)

	switch err {
	case pebble.ErrNotFound:
		fmt.Println("not found")

		if err := db.Set(key, []byte("world"), pebble.Sync); err != nil {
			log.Fatal(err)
		}
		value, closer, err := db.Get(key)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("key:%s, value:%s\n", key, value)

		if err := closer.Close(); err != nil {
			log.Fatal(err)
		}
	case nil:
		fmt.Println("found")

		if err := closer.Close(); err != nil {
			log.Fatal(err)
		}
	}

}
