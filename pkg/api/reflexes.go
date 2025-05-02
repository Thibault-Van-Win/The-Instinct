package api

import (
	"fmt"
	"net/http"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
	"github.com/labstack/echo/v4"
)

type ReflexController struct {
	service *reflex.ReflexService
}

func NewReflexController(service *reflex.ReflexService) *ReflexController {
	return &ReflexController{
		service: service,
	}
}

func (rc *ReflexController) Register(e *echo.Echo) {
	g := e.Group("/api/reflexes")

	g.POST("", rc.CreateReflex)
	g.GET("", rc.ListReflexes)
	g.GET("/:id", rc.GetReflexByID)
	g.GET("/name/:name", rc.GetReflexByName)
	g.PUT("/:id", rc.UpdateReflex)
	g.DELETE("/:id", rc.DeleteReflex)
}

// CreateReflex handles the creation of a new reflex
func (rc *ReflexController) CreateReflex(c echo.Context) error {
	var request struct {
		reflex.ReflexConfig
	}

	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	id, err := rc.service.CreateReflex(c.Request().Context(), request.ReflexConfig)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"id": id,
	})
}

// ListReflexes handles listing all reflexes
func (rc *ReflexController) ListReflexes(c echo.Context) error {
	fmt.Println("Received request for list reflexes")
	reflexes, err := rc.service.ListReflexes(c.Request().Context())
	if err != nil {
		fmt.Printf("Error while serving list reflexes: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, reflexes)
}

// GetReflexByID handles getting a reflex by ID
func (rc *ReflexController) GetReflexByID(c echo.Context) error {
	id := c.Param("id")

	reflex, err := rc.service.GetReflexByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Reflex not found")
	}

	return c.JSON(http.StatusOK, reflex)
}

// GetReflexByName handles getting a reflex by name
func (rc *ReflexController) GetReflexByName(c echo.Context) error {
	name := c.Param("name")

	reflex, err := rc.service.GetReflexByName(c.Request().Context(), name)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Reflex not found")
	}

	if reflex == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Reflex not found")
	}

	return c.JSON(http.StatusOK, reflex)
}

// UpdateReflex handles updating an existing reflex
func (rc *ReflexController) UpdateReflex(c echo.Context) error {
	id := c.Param("id")

	var request struct {
		reflex.ReflexConfig
	}

	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	err := rc.service.UpdateReflex(c.Request().Context(), id, request.ReflexConfig)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// DeleteReflex handles deleting a reflex
func (rc *ReflexController) DeleteReflex(c echo.Context) error {
	id := c.Param("id")

	err := rc.service.DeleteReflex(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
