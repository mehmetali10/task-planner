package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	postgres_repository "github.com/mehmetali10/task-planner/internal/pkg/repository/postgres"
	"github.com/mehmetali10/task-planner/internal/task/config"
	"github.com/mehmetali10/task-planner/internal/task/handler"
	"github.com/mehmetali10/task-planner/internal/task/server"
	"github.com/mehmetali10/task-planner/internal/task/service"
)

// @title Task. Planner Store API
// @version 1.0
// @description This is a Task. Planner store API with CRUD operations.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@kvstore.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	repo := postgres_repository.NewPostgresRepo()
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	httpServer := server.NewServer(handler)
	httpServer.Start(config.GetApp().HTTPAddr)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("Received shutdown signal, stopping server...")

	httpServer.Stop()

	log.Println("Servers stopped")
}
