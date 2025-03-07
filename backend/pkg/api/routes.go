package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/s0und0fs1lence/ads-zero/pkg/api/controllers"
	"github.com/s0und0fs1lence/ads-zero/pkg/db"
)

func userRoute(router *gin.Engine, prefix string, dbSvc db.DbService) {
	userGroup := router.Group(fmt.Sprintf("%s/user", prefix))
	controllers.NewUserController(dbSvc, userGroup)
}

func providerRoute(router *gin.Engine, prefix string, dbSvc db.DbService) {
	providerGroup := router.Group(fmt.Sprintf("%s/provider", prefix))
	controllers.NewProviderController(dbSvc, providerGroup)
}
