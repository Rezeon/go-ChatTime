package main

import (
	"gotry/database"
	"gotry/models"
	"gotry/routes"
	"gotry/ws"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	logger.Info("Starting application...")

	database.Connect()
	database.DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Follow{},
		&models.Message{},
	)
	logger.Info("Database connected and migrated",
		zap.Strings("tables", []string{"users", "posts", "follows", "messages"}))

	r := routes.SetupRouter()
	go ws.HandleMessage()
	logger.Info("Server running at http://localhost:8080")

	if err := r.Run(":8080"); err != nil {
		logger.Fatal("Failed to run server", zap.Error(err))
	}
}
