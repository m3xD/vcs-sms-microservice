package api

import (
	"database/sql"
	"net/http"
	db "server-management/db/sqlc"
	"server-management/token"
	"server-management/util"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	store db.Store
	token token.Maker
}

type createUserRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

func (u *UserController) createUser(ctx *gin.Context) {
	var request createUserRequest
	log := util.NewLogger()

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		log.Error("error when bind json: " + err.Error())
		return
	}

	hashPassword, err := util.HashPassword(request.Password)

	if err != nil {
		log.Error("error when hash password: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username: request.UserName,
		Password: hashPassword,
		Email:    request.Email,
		Role:     request.Role,
	}

	user, err := u.store.CreateUser(ctx, arg)

	if err != nil {
		log.Error("error when create user: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	log.Info("user created: " + user.Username)
	ctx.JSON(http.StatusOK, user)
}

type LoginUserRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *UserController) loginUser(ctx *gin.Context) {
	var request LoginUserRequest
	log := util.NewLogger()

	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Error("error when bind json: " + err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := u.store.GetUser(ctx, request.UserName)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error("user not found: " + request.UserName)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		log.Error("error when get user: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.ComparePassword(request.Password, user.Password)
	if err != nil {
		log.Error("error when compare password: " + err.Error())
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	scope, err := u.store.GetScope(ctx, user.ID)

	if err != nil {
		log.Error("error when get scope: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, err := u.token.CreateToken(user.Username, user.Role, scope, 10*time.Minute)
	if err != nil {
		log.Error("error when create token: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	log.Info("user login: " + user.Username)
	ctx.JSON(http.StatusOK, response)
}

type updateRoleRequest struct {
	ID   int    `json:"id" binding:"required"`
	Role string `json:"role" binding:"required"`
}

func (u *UserController) updateRole(ctx *gin.Context) {
	var rq updateRoleRequest
	log := util.NewLogger()

	if err := ctx.ShouldBindJSON(&rq); err != nil {
		log.Error("error when bind json: " + err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateRoleParams{
		ID:   int64(rq.ID),
		Role: rq.Role,
	}

	user, err := u.store.UpdateRole(ctx, arg)

	if err != nil {
		log.Error("error when update role: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type userResponse struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		UserName: user.Username,
		Email:    user.Email,
	}
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}
