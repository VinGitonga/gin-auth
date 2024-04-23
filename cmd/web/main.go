package main

import (
	"context"
	"encoding/gob"
	"log"
	"os"

	"github.com/VinGitonga/gin-auth/driver"
	"github.com/VinGitonga/gin-auth/handlers"
	"github.com/VinGitonga/gin-auth/modules/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var app config.GoAppTools

func main() {
	// Register gob for streams
	gob.Register(map[string]interface{}{})
	gob.Register(primitive.NewObjectID())

	InfoLogger := log.New(os.Stdout, " ", log.LstdFlags|log.Lshortfile)
	ErrorLogger := log.New(os.Stdout, " ", log.LstdFlags|log.Lshortfile)

	app.InfoLogger = *InfoLogger
	app.ErrorLogger = *ErrorLogger

	err := godotenv.Load()

	if err != nil {
		app.ErrorLogger.Fatal("No .env file available")
	}

	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		app.ErrorLogger.Fatalln("Mongo URI string not found: ")
	}

	client := driver.Connection(uri)

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			app.ErrorLogger.Fatal(err)
			return
		}
	}()

	appRouter := gin.New()

	goApp := handlers.NewGoApp(&app, client)

	Routes(appRouter, goApp)

	err = appRouter.Run()

	if err != nil {
		log.Fatal(err)
	}
}
