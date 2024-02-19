package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lucsky/cuid"
	"github.com/tiagokriok/oka/internal/repositories"
	"github.com/tiagokriok/oka/internal/services"
)

type LinkHandler struct {
	linkSvc *services.LinkService
}

func NewLinkHandler(linkSvc *services.LinkService) *LinkHandler {
	return &LinkHandler{
		linkSvc,
	}
}

func (lh *LinkHandler) Create(c echo.Context) error {
	var link repositories.Link

	if err := c.Bind(&link); err != nil {
		slog.Error("Error bad params link", err)
		return c.NoContent(http.StatusBadRequest)
	}

	link.ID = cuid.New()
	link.Key = cuid.Slug()

	newLink, err := lh.linkSvc.Create(&link)
	if err != nil {
		slog.Error("Error creating link", err)
		return c.NoContent(http.StatusBadRequest)
	}

	c.Response().Header().Set("Location", fmt.Sprintf("/%s", newLink.Key))

	return c.JSON(http.StatusCreated, link)
}
