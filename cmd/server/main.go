package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lucsky/cuid"
)

type KeyParam struct {
	Key string `param:"key"`
}

type Link struct {
	ID  string `json:"id,omitempty"`
	URL string `json:"url"`
	Key string `json:"key,omitempty"`
}

type DB struct {
	*sqlx.DB
}

type LinkRepository struct {
	db *DB
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
	slog.Info("Connected to database")

	defer db.Close()

	linkRepository := NewLinkRepository(db)

	e := echo.New()

	e.Use(customLogger())
	e.Use(middleware.Recover())

	e.GET("/:key", func(c echo.Context) error {
		var params KeyParam

		if err := c.Bind(&params); err != nil {
			return err
		}

		link, err := linkRepository.GetLinkByKey(params.Key)
		if err != nil {
			slog.Error("Error get linking by key", err)
		}

		return c.Redirect(http.StatusPermanentRedirect, link.URL)

	})

	api := e.Group("/api")

	api.GET("/healthcheck", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	v1 := api.Group("/v1")

	links := v1.Group("/links")

	links.POST("", func(c echo.Context) error {
		var link Link

		if err := c.Bind(&link); err != nil {
			return err
		}

		link.ID = cuid.New()
		link.Key = cuid.Slug()

		err := linkRepository.CreateLink(&link)
		if err != nil {
			slog.Error("Error creating link", err)
		}

		c.Response().Header().Set("Location", fmt.Sprintf("/%s", link.Key))

		return c.JSON(http.StatusOK, link)
	})

	e.Logger.Fatal(e.Start(":3333"))
}

func NewDB() (*DB, error) {
	db, err := sqlx.Connect("mysql", os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &DB{db}, err
}

func (db *DB) Close() {
	db.DB.Close()
}

func NewLinkRepository(db *DB) *LinkRepository {
	slog.Info("New Link Repository")
	return &LinkRepository{
		db,
	}
}

func (lr *LinkRepository) CreateLink(link *Link) error {
	_, err := lr.db.Exec("INSERT INTO links (id, `key`, url) VALUES (?, ?, ?)", link.ID, link.Key, link.URL)
	if err != nil {
		return err
	}

	return nil
}

func (lr *LinkRepository) GetLinkByKey(key string) (*Link, error) {
	var link Link

	err := lr.db.Get(&link, "SELECT * FROM links WHERE `key`=?", key)
	if err != nil {
		return &link, err
	}

	return &link, nil
}
