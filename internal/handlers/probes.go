package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthProbe struct {
	Status string `json:"status"`
}

// ReadinessProbe godoc
//
//	@Summary		Readiness Probe
//	@Description	Check if server is ready to accept requests
//	@Id				readinessProbe
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	HealthProbe
//	@Router			/health/readiness [get]
func ReadinessProbe(c *gin.Context) {
	c.JSON(http.StatusOK, HealthProbe{
		Status: "ready",
	})
}

//  LivenessProbe godoc

// @Summary		Liveness Probe
// @Description	Check if server is up and running
// @Id				livenessProbe
// @Tags			health
// @Accept			json
// @Produce		json
// @Success		200	{object}	HealthProbe
// @Router			/health/liveness [get]
func LivenessProbe(c *gin.Context) {
	c.JSON(http.StatusOK, HealthProbe{
		Status: "ready",
	})
}
