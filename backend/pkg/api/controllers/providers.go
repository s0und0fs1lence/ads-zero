package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/s0und0fs1lence/ads-zero/pkg/db"
)

func NewProviderController(dbSvc db.DbService, group *gin.RouterGroup) {
	group.GET("/", handleGetProviders(dbSvc))
	group.POST("/create", handleCreateProvider(dbSvc))
	group.PUT("/update", handleUpdateProvider(dbSvc))
}

// TODO:
// define if the user is valid, and return only the necessary information to the frontend
// if it doesn't exist, we should pick if it's responsability of the frontend to send the information for the insert inside our system, otherwise we should proceed to the insert here

func handleGetProviders(db_svc db.DbService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Query("uid")
		if userId == "" {
			pid := ctx.Param("pid")
			if pid == "" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "missing user id or provider id",
				})
				return
			}
			provider, err := db_svc.GetProviderByID(pid)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"data": []db.DbProvider{*provider},
			})
			return
			// get by provider id

		}
		// get by user id
		providers, err := db_svc.GetProvidersByClientID(userId)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data": providers,
		})

	}
}

func handleCreateProvider(db_svc db.DbService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//TODO: crete the user
		var createReq db.ProviderCreate
		if err := ctx.BindJSON(&createReq); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// verify if the payload is valid
		if !createReq.IsValid() {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}
		// check if the user exists
		_, err := db_svc.GetClientByID(createReq.ClientID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "invalid client_id",
			})
			return
		}
		provider := createReq.AsDbProvider()
		provider.ProviderID = ulid.Make().String()
		provider.InsertedAt = time.Now().UTC()

		if err := db_svc.InsertProvider(&provider); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"data": provider,
		})

	}
}

func handleUpdateProvider(dbSvc db.DbService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var update db.ProviderUpdate
		if err := ctx.ShouldBindJSON(&update); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}
		provider, err := dbSvc.UpdateProvider(&update)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data": provider,
		})
	}
}
