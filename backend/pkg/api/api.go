package api

import (
	"github.com/gin-gonic/gin"
	"github.com/s0und0fs1lence/ads-zero/pkg/db"
)

func registerRoutes(r *gin.Engine, prefix string, dbSvc db.DbService) {
	userRoute(r, prefix, dbSvc)
	providerRoute(r, prefix, dbSvc)
}

func StartAPI(dbSvc db.DbService) error {
	r := gin.Default()
	r.UseH2C = true
	registerRoutes(r, "/api/v1", dbSvc)
	if err := r.Run(); err != nil {
		return err
	}
	return nil
}
