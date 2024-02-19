package handlers

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tiagokriok/oka/internal/services"
)

type PublicHandler struct {
	publicSvc *services.PublicService
}

func NewPublicHandler(publicSvc *services.PublicService) *PublicHandler {
	return &PublicHandler{
		publicSvc,
	}
}

func (ph *PublicHandler) Redirect(c echo.Context) error {
	key := c.Param("key")

	if key == "" {
		slog.Error("Error missing params")
		return c.NoContent(http.StatusBadRequest)
	}

	link, err := ph.publicSvc.GetLinkByKey(key)
	if err != nil {
		slog.Error("Error get linking by key", err)
		return c.NoContent(http.StatusNotFound)
	}

	return c.Redirect(http.StatusPermanentRedirect, link.URL)
}
