package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	db "server-management/db/sqlc"
	"server-management/util"
)

type ScopeController struct {
	store db.Store
}

type createScopeRequest struct {
	Name string `json:"name" binding:"required"`
	Role string `json:"role" binding:"required"`
}

func (s *ScopeController) createScope(ctx *gin.Context) {
	log := util.NewLogger()
	var rq createScopeRequest

	if err := ctx.ShouldBindBodyWithJSON(&rq); err != nil {
		log.Error("error when bind json: " + err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateScopeParams{
		Name: rq.Name,
		Role: rq.Role,
	}

	scope, err := s.store.CreateScope(ctx, arg)

	if err != nil {
		log.Error("error when create scope: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	log.Info(fmt.Sprintf("created scope: %v", scope))
	ctx.JSON(http.StatusOK, scope)
}

type getScopeRequest struct {
	ID int64 `json:"id" binding:"required"`
}

func (s *ScopeController) getScope(ctx *gin.Context) {
	var rq getScopeRequest
	log := util.NewLogger()

	if err := ctx.ShouldBindBodyWithJSON(&rq); err != nil {
		log.Error("error when bind json: " + err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	scopes, err := s.store.GetScope(ctx, rq.ID)

	if err != nil {
		log.Error("error when get scope: " + err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Info(fmt.Sprintf("get scope: %v", scopes))
	ctx.JSON(http.StatusOK, scopes)
}

type updateScopeRequest struct {
	ID      int64  `json:"id" binding:"required"`
	SetName bool   `json:"setName"`
	Name    string `json:"name"`
	SetRole bool   `json:"setRole"`
	Role    string `json:"role"`
}

func (s *ScopeController) updateScope(ctx *gin.Context) {
	var rq updateScopeRequest
	log := util.NewLogger()

	if err := ctx.ShouldBindBodyWithJSON(&rq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		log.Error("error when bind json: " + err.Error())
		return
	}

	arg := db.UpdateScopeParams{
		ID:      rq.ID,
		SetName: rq.SetName,
		Name:    rq.Name,
		SetRole: rq.SetRole,
		Role:    rq.Role,
	}

	scope, err := s.store.UpdateScope(ctx, arg)

	if err != nil {
		log.Error("error when update scope: " + err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Info(fmt.Sprintf("updated scope: %v", scope))
	ctx.JSON(http.StatusOK, scope)
}

type deleteScopeRequest struct {
	ID int64 `json:"id" binding:"required"`
}

func (s *ScopeController) deleteScope(ctx *gin.Context) {
	var rq deleteScopeRequest
	log := util.NewLogger()

	if err := ctx.ShouldBindBodyWithJSON(&rq); err != nil {
		log.Error("error when bind json: " + err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.store.DeleteScope(ctx, rq.ID)

	if err != nil {
		log.Error("error when delete scope: " + err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Info("delete scope successfully")
	ctx.JSON(http.StatusOK, "delete successfully")
}
