package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tiagokriok/oka/internal/middlewares"
	"github.com/tiagokriok/oka/internal/repositories"
	"github.com/tiagokriok/oka/internal/storages"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lucsky/cuid"
)

func main() {
	slog.Info("OKA init...")

	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	db, err := storages.NewMysqlDB()
	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to database")

	defer db.Close()

	linkRepository := repositories.NewLinkRepository(db)

	e := echo.New()

	e.Use(middlewares.CustomLogger())
	e.Use(middleware.Recover())

	e.GET("/:key", func(c echo.Context) error {
		var params repositories.Link

		if err := c.Bind(&params); err != nil {
			slog.Error("Error params not found", err)
			return c.NoContent(http.StatusNotFound)
		}

		link, err := linkRepository.GetLinkByKey(params.Key)
		if err != nil {
			slog.Error("Error get linking by key", err)
			return c.NoContent(http.StatusNotFound)
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
		var link repositories.Link

		if err := c.Bind(&link); err != nil {
			slog.Error("Error bad params link", err)
			return c.NoContent(http.StatusBadRequest)
		}

		link.ID = cuid.New()
		link.Key = cuid.Slug()

		err := linkRepository.CreateLink(&link)
		if err != nil {
			slog.Error("Error creating link", err)
			return c.NoContent(http.StatusBadRequest)
		}

		c.Response().Header().Set("Location", fmt.Sprintf("/%s", link.Key))

		return c.JSON(http.StatusCreated, link)
	})

	e.Logger.Fatal(e.Start(":3333"))
}
