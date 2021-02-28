package springkilometers

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	cache "github.com/go-redis/cache/v8"
	models "github.com/ondrejholik/springkilometers/models"
)

// ShowAchievmentsPage --
func ShowAchievmentsPage(c *gin.Context) {

	claims, err := ClaimsUser(c)
	if err == nil {
		log.Println(err)
	}
	//achievments := models.GetAchievmentsByUserID(claims.UserID)
	var achievments []models.Achievment

	if err := models.MyCache.Get(models.Ctx, "achievments:"+strconv.Itoa(claims.UserID), &achievments); err != nil {
		achievments = models.GetAchievmentsByUserID(claims.UserID)
		if err := models.MyCache.Set(&cache.Item{
			Ctx:   models.Ctx,
			Key:   "achievments:" + strconv.Itoa(claims.UserID),
			Value: achievments,
			TTL:   10 * time.Hour,
		}); err != nil {
			panic(err)
		}

	}

	// Call the render function with the name of the template to render
	Render(c, gin.H{
		"title":   "Achievments",
		"payload": achievments,
	}, "achievments.html")
	//}
}
