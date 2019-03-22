package main

import (
	"fmt"
	"github.com/kycklingar/FurLoaderGO/data"
	"log"
)

func main() {
	db := data.OpenDB()
	defer db.Close()

	for i := 0; i < 10; i++ {
		if err := db.Store(fmt.Sprintf("key %d"), fmt.Sprintf("value %d", i)); err != nil {
			log.Fatal(err)
		}
	}

	for i := 0; i < 10; i++{
		fmt.Println(db.Get(fmt.Sprintf("key %d", i)))
	}
}
