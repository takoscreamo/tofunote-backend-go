package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"tofunote-backend/domain/user"
	"tofunote-backend/infra"

	"tofunote-backend/usecases"

	"github.com/cmackenzie1/go-uuid"
	"github.com/gin-gonic/gin"
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
	id, err := uuid.NewV7()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "UUID生成に失敗しました"})
		return
	}
	newUUID := id.String()
	refreshToken, err := generateRefreshToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "リフレッシュトークン生成に失敗しました"})
		return
	}
	u := &user.User{
		ID:           newUUID,
		Nickname:     "ゲスト",
		IsGuest:      true,
		RefreshToken: refreshToken,
	}
	if err := c.repo.Create(ctx.Request.Context(), u); err != nil {
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
	u, err := c.repo.FindByRefreshToken(ctx.Request.Context(), req.RefreshToken)
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
	err := c.withdrawUsecase.Withdraw(ctx.Request.Context(), userIDStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "退会処理に失敗しました: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "ユーザーと日記データを全て削除しました"})
}

// GET /me: ユーザー情報取得API
func (c *UserController) GetMe(ctx *gin.Context) {
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
	u, err := c.repo.FindByID(ctx.Request.Context(), userIDStr)
	if err != nil || u == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーが見つかりません"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":       u.ID,
		"nickname": u.Nickname,
		// 必要に応じて他の項目も追加
	})
}

// PATCH /me: ユーザー情報部分更新API
func (c *UserController) PatchMe(ctx *gin.Context) {
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
	u, err := c.repo.FindByID(ctx.Request.Context(), userIDStr)
	if err != nil || u == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーが見つかりません"})
		return
	}
	var req map[string]interface{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストデータです"})
		return
	}
	updated := false
	if nickname, ok := req["nickname"].(string); ok {
		u.Nickname = nickname
		updated = true
	}
	// 他の項目もここで追加可能
	if !updated {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "更新可能な項目がありません"})
		return
	}
	if err := c.repo.Update(ctx.Request.Context(), u); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー情報の更新に失敗しました"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ユーザー情報を更新しました",
		"user": gin.H{
			"id":       u.ID,
			"nickname": u.Nickname,
			// 必要に応じて他の項目も追加
		},
	})
}
