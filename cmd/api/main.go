package main

import (
	"log"
	"net/http"

	"github.com/NITHISH-2006/taskflow-go/internal/handler"
	"github.com/NITHISH-2006/taskflow-go/internal/middleware"
	"github.com/NITHISH-2006/taskflow-go/internal/repository"
	"github.com/NITHISH-2006/taskflow-go/internal/service"
)

func main() {
	repo := repository.NewInMemoryTaskRepository()
	taskService := service.NewTaskService(repo)
	taskHandler := handler.NewTaskHandler(taskService)

	mux := http.NewServeMux()
	taskHandler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: middleware.RequestLogger(mux),
	}

	log.Println("starting TaskFlow API on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
