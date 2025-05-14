package station

import (
	"github.com/aldngrha/mrt-schedules.git/modules/common/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Initiate(router *gin.RouterGroup) {
	stationService := NewService()
	station := router.Group("/station")
	station.GET("", func(c *gin.Context) {
		//code service
		GetAllStation(c, stationService)
	})
}

func GetAllStation(c *gin.Context, service Service) {
	datas, err := service.GetAllStations()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Get all stations success",
		Data:    datas,
	})
}
