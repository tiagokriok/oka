package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/tiagokriok/oka/internal/handlers"
	"github.com/tiagokriok/oka/internal/storages"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	slog.Info("OKA init...")
	slog.Info("Loading env")
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
		os.Exit(1)
	}
}

func main() {
	slog.Info("OKA starting server...")

	db, err := storages.NewMysqlDB()
	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to database")

	defer db.Close()

	router, err := handlers.NewHttp(db)
	if err != nil {
		slog.Error("Error initializing router", "error", err)
		os.Exit(1)
	}

	router.Serve(os.Getenv("PORT"))
}
