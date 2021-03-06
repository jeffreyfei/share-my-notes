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

// Bootstraps and runs the provider server

var (
	db        *gorm.DB
	webServer *server.Server
)

// Initializes the database and automgirates all the models
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

// Initializes the server instance
func initServer() {
	clientID := os.Getenv("GOOGLEKEY")
	clientSecret := os.Getenv("GOOGLESECRET")
	baseURL := os.Getenv("BASE_URL")
	sessionKey := os.Getenv("SESSION_KEY")
	lbPubURL := os.Getenv("LB_PUB_URL")
	lbPriURL := os.Getenv("LB_PRI_URL")
	webServer = server.NewServer(db, baseURL, sessionKey, clientID, clientSecret, lbPriURL, lbPubURL)
	webServer.StartBufferProc()
}

// Bootstraping
func main() {
	initDB()
	initServer()
	port := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	log.Infof("Server running on port %s", port)
	go http.ListenAndServe(port, webServer.Router)
	go webServer.RegisterLoadBalancer()
	fmt.Scanln()
}
