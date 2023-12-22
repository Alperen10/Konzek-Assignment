package router

import (
	"net/http"

	_ "github.com/Alperen10/Konzek-App/docs"
	"github.com/Alperen10/Konzek-App/handler"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// Router is exported and used in main.go
func Router() http.Handler {
	router := http.NewServeMux()

	// Swagger documentation route
	router.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	router.HandleFunc("/api/singleTask/", handler.GetSingleTask)
	router.HandleFunc("/api/task", handler.GetAllTasks)
	router.HandleFunc("/api/newTask", handler.CreateTask)
	router.HandleFunc("/api/updateTask/", handler.UpdateTask)
	router.HandleFunc("/api/deleteTask/", handler.DeleteTaskByID)

	return router
}
