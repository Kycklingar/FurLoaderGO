package data

import (
	"github.com/dgraph-io/badger"
	"log"
)

type DB struct{
	*badger.DB
}

func OpenDB()*DB{
	opts := badger.DefaultOptions
	opts.Dir = "./badger"
	opts.ValueDir = "./badger/value"

	db, err := badger.Open(opts)
	if err != nil{
		log.Fatal(err)
	}
	return &DB{db}
}


func (db *DB) Store(key, value string)error{
	return db.Update(func(txn *badger.Txn)error{
		err := txn.Set([]byte(key), []byte(value))
		return err
	})
}

func (db *DB) Get(key string)string{
	var value string
	err := db.View(func(txn *badger.Txn)error{
		item, err := txn.Get([]byte(key))
		if err != nil{
			log.Println(err)
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil{
			log.Println(err)
			return err
		}
		value = string(val)
		return nil
	})
	if err != nil{
		log.Println(err)
		return ""
	}
	return value
}
