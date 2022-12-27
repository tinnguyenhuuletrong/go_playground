package pebble_play

import (
	"fmt"
	"log"
	"strconv"

	"github.com/cockroachdb/pebble"
)

func Pebble_Db_Play() {
	db, err := pebble.Open("./tmp/pebble_demo", &pebble.Options{})
	if err != nil {
		log.Fatal(err)
	}

	key := []byte("one")
	db.Set(key, []byte("1"), pebble.NoSync)

	for i := 0; i < 100; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		db.Set(key, []byte(strconv.Itoa(i)), pebble.NoSync)
	}

	value, _, err := db.Get(key)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Get by key: one")
	log.Println("value:", string(value))

	log.Println("Scan: Begin")

	it := db.NewIter(&pebble.IterOptions{})
	it.First()
	log.Println("\tIter Begin: ", string(it.Key()))
	it.Last()
	log.Println("\tIter End: ", string(it.Key()))

	it.First()
	for {
		if !it.Valid() {
			break
		}

		log.Printf("\tkey: %s, val: %s", string(it.Key()), string(it.Value()))

		it.Next()
	}

	log.Println("Scan: End")

	log.Println("Metrics:", db.Metrics())

	log.Println("Compact: Begin")
	db.Compact([]byte("key_0"), []byte("key_9999"), true)
	log.Println("Compact: End")
	db.Close()
}
