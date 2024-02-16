package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"database/sql"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/go-sql-driver/mysql"
)

type Link struct {
	ID          string `json:"id,omitempty"`
	URL         string `json:"url"`
	ShortKey    string `json:"shortKey"`
	RedirectURL string `json:"redirectUrl"`
}

type DB struct {
	*sql.DB
}

func customLogger() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339_nano} - ${uri} [${method} - ${status}] ${latency_human} - ${error}\n",
	})
}

func main() {
	slog.Info("OKA init...")

	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	db, err := NewDB()

	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}

	defer db.Close()

	e := echo.New()

	e.Use(customLogger())
	e.Use(middleware.Recover())

	// e.GET("/:shortKey")

	api := e.Group("/api")

	api.GET("/healthcheck", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	v1 := api.Group("/v1")

	links := v1.Group("/links")
	//
	links.GET("", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	links.GET("/:id", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Param(":id"))
	})

	links.POST("", func(c echo.Context) error {
		var link Link

		if err := c.Bind(&link); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, link)
	})

	links.DELETE("/:id", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start(":3333"))

}

func NewDB() (*DB, error) {

	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &DB{db}, err
}

func (db *DB) Close() {
	db.DB.Close()
}
