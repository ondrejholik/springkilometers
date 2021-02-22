package springkilometers

import (
	"github.com/gin-gonic/gin"
	models "github.com/ondrejholik/springkilometers/models"
)

// Chat --
func Chat(c *gin.Context) {
	tripID := c.Param("tripID")
	claims, err := ClaimsUser(c)
	if err == nil {
		models.ServeWs(c.Writer, c.Request, tripID, claims.UserID)
	}
}
