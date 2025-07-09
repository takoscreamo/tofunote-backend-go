package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"feelog-backend/domain/user"
	"feelog-backend/infra"
	"net/http"

	"feelog-backend/usecases"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	repo            user.Repository
	withdrawUsecase *usecases.UserWithdrawUsecase // 退会用のみ残す
}

func NewUserController(repo user.Repository, withdrawUsecase *usecases.UserWithdrawUsecase) *UserController {
	return &UserController{repo: repo, withdrawUsecase: withdrawUsecase}
}

type GuestLoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ID           string `json:"id"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

var generateTokenForTest = func(id string) (string, error) {
	return infra.GenerateToken(id)
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GuestLogin: サーバー側でUUIDとリフレッシュトークンを生成し、ユーザー作成・トークン発行API
func (c *UserController) GuestLogin(ctx *gin.Context) {
	newUUID := uuid.New().String()
	refreshToken, err := generateRefreshToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "リフレッシュトークン生成に失敗しました"})
		return
	}
	u := &user.User{
		ID:           newUUID,
		IsGuest:      true,
		RefreshToken: refreshToken,
	}
	if err := c.repo.Create(u); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー作成に失敗しました"})
		return
	}
	token, err := generateTokenForTest(u.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "トークン生成に失敗しました"})
		return
	}
	ctx.JSON(http.StatusOK, GuestLoginResponse{Token: token, RefreshToken: refreshToken, ID: newUUID})
}

// RefreshToken: リフレッシュトークンで新しいJWTを発行
func (c *UserController) RefreshToken(ctx *gin.Context) {
	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	// リフレッシュトークンでユーザー検索
	u, err := c.repo.FindByRefreshToken(req.RefreshToken)
	if err != nil || u == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	token, err := infra.GenerateToken(u.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "トークン生成に失敗しました"})
		return
	}
	ctx.JSON(http.StatusOK, RefreshTokenResponse{Token: token})
}

// ユーザー自身によるアカウント削除API
func (c *UserController) DeleteMe(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証情報が見つかりません"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの形式が不正です"})
		return
	}
	// 退会ユースケースで一括削除
	err := c.withdrawUsecase.Withdraw(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "退会処理に失敗しました: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "ユーザーと日記データを全て削除しました"})
}
