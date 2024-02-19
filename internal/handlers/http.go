package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tiagokriok/oka/internal/middlewares"
	"github.com/tiagokriok/oka/internal/repositories"
	"github.com/tiagokriok/oka/internal/services"
	"github.com/tiagokriok/oka/internal/storages"
)

type Router struct {
	*echo.Echo
}

func NewHttp(db *storages.MysqlDB) (*Router, error) {
	e := echo.New()

	e.Use(middlewares.CustomLogger())
	e.Use(middleware.Recover())

	api := e.Group("/api")

	api.GET("/healthcheck", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	v1 := api.Group("/v1")

	// Links
	linkRepo := repositories.NewLinkRepository(db)
	linkService := services.NewLinkService(linkRepo)
	linkHandler := NewLinkHandler(linkService)
	linksRouter := v1.Group("/links")

	linksRouter.POST("", linkHandler.Create)

	// Public
	publicService := services.NewPublicService(linkRepo)
	publicHandler := NewPublicHandler(publicService)

	e.GET("/:key", publicHandler.Redirect)

	return &Router{
		e,
	}, nil
}

func (router *Router) Serve(listenAddr string) {
	router.Logger.Fatal(router.Start(fmt.Sprintf(":%s", listenAddr)))
}
