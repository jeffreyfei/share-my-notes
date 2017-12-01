package main

import (
	"fmt"
	"net/http"
	"os"

	dbLib "github.com/jeffreyfei/share-my-notes/server/lib/db"
	"github.com/jeffreyfei/share-my-notes/server/lib/md_note"
	"github.com/jeffreyfei/share-my-notes/server/lib/server"
	"github.com/jeffreyfei/share-my-notes/server/lib/user"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

var (
	db        *gorm.DB
	webServer *server.Server
)

func initDB() {
	var err error
	if db, err = dbLib.GetDB(); err != nil {
		log.Fatal(err)
	}
	if err := user.AutoMigrate(db); err != nil {
		log.WithField("err", err).Error("Failed to migrate user model")
		os.Exit(1)
	}
	if err := md_note.AutoMigrate(db); err != nil {
		log.WithField("err", err).Error("Failed to migrate md note model")
		os.Exit(1)
	}
}

func initServer() {
	clientID := os.Getenv("GOOGLEKEY")
	clientSecret := os.Getenv("GOOGLESECRET")
	baseURL := os.Getenv("BASE_URL")
	sessionKey := os.Getenv("SESSION_KEY")
	webServer = server.NewServer(db, baseURL, sessionKey, clientID, clientSecret)
}

func main() {
	initDB()
	initServer()
	port := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	log.Infof("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(port, webServer.Router))
}
