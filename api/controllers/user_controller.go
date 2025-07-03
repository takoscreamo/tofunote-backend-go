package controllers

import (
	"feelog-backend/domain/user"
	"feelog-backend/infra"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	repo user.Repository
}

func NewUserController(repo user.Repository) *UserController {
	return &UserController{repo: repo}
}

type GuestLoginRequest struct {
	ID string `json:"id" binding:"required,uuid4"`
}

type GuestLoginResponse struct {
	Token string `json:"token"`
}

// GuestLogin: UUIDのみでユーザー作成・トークン発行API
func (c *UserController) GuestLogin(ctx *gin.Context) {
	var req GuestLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	u, err := c.repo.FindByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー検索に失敗しました"})
		return
	}
	if u == nil {
		u = &user.User{
			ID:      req.ID,
			IsGuest: true,
		}
		if err := c.repo.Create(u); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー作成に失敗しました"})
			return
		}
	}
	token, err := infra.GenerateToken(u.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "トークン生成に失敗しました"})
		return
	}
	ctx.JSON(http.StatusOK, GuestLoginResponse{Token: token})
}
