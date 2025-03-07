package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/s0und0fs1lence/ads-zero/pkg/db"
)

func NewUserController(dbSvc db.DbService, group *gin.RouterGroup) {
	group.GET("/", handleGetUserById(dbSvc))
	group.POST("/create", handleCreateUser(dbSvc))
	group.PUT("/update", handleUpdateUser(dbSvc))
	//account spend
	group.GET("/accounts/spend", handleGetAccountSpend(dbSvc))
	group.GET("/accounts/spend/grouped", handleGetAccountSpendGrouped(dbSvc))
	//campaigns
	group.GET("/campaigns/spend", handleGetCampaignSpend(dbSvc))
	group.GET("/campaigns/spend/grouped", handleGetCampaignSpendGrouped(dbSvc))

	//rules
	group.GET("/rules", handleGetRules(dbSvc))
	group.POST("/rules/create", handleCreateRule(dbSvc))
	group.PUT("/rules/update", handleUpdateRule(dbSvc))
	group.DELETE("/rules/delete", handleDeleteRule(dbSvc))
}

func handleDeleteRule(dbSvc db.DbService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		panic("unimplemented")
	}
}

func handleUpdateRule(dbSvc db.DbService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		panic("unimplemented")
	}
}

func handleCreateRule(dbSvc db.DbService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		panic("unimplemented")
	}
}

func handleGetRules(dbSvc db.DbService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Query("uid")
		user, err := dbSvc.GetRulesByClientID(userId)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data": user,
		})
	}
}

func handleGetCampaignSpend(dbSvc db.DbService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req db.ClientSpendRequest
		if err := ctx.BindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := dbSvc.GetCampaignSpend(req.ClientID, req.Start, req.End)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data": res,
		})
	}
}

func handleGetCampaignSpendGrouped(dbSvc db.DbService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req db.ClientSpendRequest
		if err := ctx.BindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := dbSvc.GetCampaignSpendGrouped(req.ClientID, req.Start, req.End)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data": res,
		})
	}
}

func handleGetAccountSpend(dbSvc db.DbService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req db.ClientSpendRequest
		if err := ctx.BindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := dbSvc.GetAccountSpend(req.ClientID, req.Start, req.End)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data": res,
		})
	}
}

func handleGetAccountSpendGrouped(dbSvc db.DbService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req db.ClientSpendRequest
		if err := ctx.BindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := dbSvc.GetAccountSpendGrouped(req.ClientID, req.Start, req.End)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data": res,
		})
	}
}

// TODO:
// define if the user is valid, and return only the necessary information to the frontend
// if it doesn't exist, we should pick if it's responsability of the frontend to send the information for the insert inside our system, otherwise we should proceed to the insert here

func handleGetUserById(db_svc db.DbService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Query("uid")
		user, err := db_svc.GetClientByID(userId)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"data": user,
		})

	}
}

func handleCreateUser(db_svc db.DbService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//TODO: crete the user
		var createReq db.ClientCreate
		if err := ctx.BindJSON(&createReq); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if !createReq.IsValid() {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}
		_, err := db_svc.GetClientByID(createReq.ClientID)
		if err == nil {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "user already exists",
			})
			return
		}

		user := db.DbClient{
			ClientID:   createReq.ClientID,
			UserEmail:  createReq.Email,
			InsertedAt: time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		}
		if err := db_svc.InsertClient(&user); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"data": user,
		})

	}
}

func handleUpdateUser(dbSvc db.DbService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var update db.ClientUpdate
		if err := ctx.ShouldBindJSON(&update); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload",
			})
			return
		}
		client, err := dbSvc.UpdateClient(&update)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data": client,
		})
	}
}
