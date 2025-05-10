package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/triggerconfig"
)

type TriggerConfigController struct {
	service *triggerconfig.TriggerConfigService
}

func NewTriggerConfigController(service *triggerconfig.TriggerConfigService) *TriggerConfigController {
	return &TriggerConfigController{
		service: service,
	}
}

func (tcc *TriggerConfigController) Register(e *echo.Echo) {
	g := e.Group("/api/triggerconfigs")

	g.POST("", tcc.CreateTriggerConfig)
	g.GET("", tcc.ListTriggerConfigs)
	g.GET("/:id", tcc.GetTriggerConfigByID)
	g.GET("/name/:name", tcc.GetTriggerConfigByName)
	g.PUT("/:id", tcc.UpdateTriggerConfig)
	g.DELETE("/:id", tcc.DeleteTriggerConfig)
}

// CreateTriggerConfig handles the creation of a new trigger config
func (rc *TriggerConfigController) CreateTriggerConfig(c echo.Context) error {
	var request struct {
		triggerconfig.TriggerConfig
	}

	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	id, err := rc.service.CreateTriggerConfig(c.Request().Context(), request.TriggerConfig)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"id": id,
	})
}

// ListTriggerConfigs handles listing all trigger configs
func (rc *TriggerConfigController) ListTriggerConfigs(c echo.Context) error {
	triggerconfigs, err := rc.service.ListTriggerConfigs(c.Request().Context())
	if err != nil {
		fmt.Printf("Error while serving list trigger configs: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, triggerconfigs)
}

// GetTriggerConfigByID handles getting a trigger config by ID
func (rc *TriggerConfigController) GetTriggerConfigByID(c echo.Context) error {
	id := c.Param("id")

	triggerconfig, err := rc.service.GetTriggerConfigByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Trigger config not found")
	}

	return c.JSON(http.StatusOK, triggerconfig)
}

// GetTriggerConfigByName handles getting a trigger config by name
func (rc *TriggerConfigController) GetTriggerConfigByName(c echo.Context) error {
	name := c.Param("name")

	triggerconfig, err := rc.service.GetTriggerConfigByName(c.Request().Context(), name)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Trigger config not found")
	}

	if triggerconfig == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Trigger config not found")
	}

	return c.JSON(http.StatusOK, triggerconfig)
}

// UpdateTriggerConfig handles updating an existing trigger config
func (rc *TriggerConfigController) UpdateTriggerConfig(c echo.Context) error {
	id := c.Param("id")

	var request struct {
		triggerconfig.TriggerConfig
	}

	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	err := rc.service.UpdateTriggerConfig(c.Request().Context(), id, request.TriggerConfig)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// DeleteTriggerConfig handles deleting a trigger config
func (rc *TriggerConfigController) DeleteTriggerConfig(c echo.Context) error {
	id := c.Param("id")

	err := rc.service.DeleteTriggerConfig(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
