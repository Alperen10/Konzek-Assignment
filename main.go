package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Alperen10/Konzek-App/database"
	"github.com/Alperen10/Konzek-App/handler"
	"github.com/Alperen10/Konzek-App/router"
)

// @title Task Service API
// @version 1.0
// @description This is a Task server.

// @host localhost:8080
// @BasePath /api
func main() {

	err := database.Initialize()
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}

	handler.InitLog()

	r := router.Router()

	fmt.Println("Starting server on the port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
