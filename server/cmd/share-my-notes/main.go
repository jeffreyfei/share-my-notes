package main

import (
	"fmt"
	"net/http"
	"os"

	dbLib "github.com/jeffreyfei/share-my-notes/server/lib/db"
	"github.com/jeffreyfei/share-my-notes/server/lib/server"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

var (
	db        *gorm.DB
	webServer *server.Server
)

func initDB() {
	var err error
	if db, err = dbLib.InitDB(); err != nil {
		log.Fatal(err)
	}
}

func initServer() {
	webServer = server.NewServer(db)
}

func main() {
	initDB()
	initServer()
	port := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	log.Infof("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(port, webServer.Router))
}
