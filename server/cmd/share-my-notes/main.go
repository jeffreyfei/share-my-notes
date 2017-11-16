package main

import (
	"log"

	dbLib "github.com/jeffreyfei/share-my-notes/server/lib/db"
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

func initDB() {
	var err error
	if db, err = dbLib.InitDB(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()
}
