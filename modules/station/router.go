package station

import "github.com/gin-gonic/gin"

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
	}
}
