package controllers

import (
	"feelog-backend/domain/user"
	"feelog-backend/infra"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	repo user.Repository
}

func NewUserController(repo user.Repository) *UserController {
	return &UserController{repo: repo}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname" binding:"required"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

// Login: ログインAPI
func (c *UserController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	u, err := c.repo.FindByEmail(req.Email)
	if err != nil || u == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	token, err := infra.GenerateToken(u.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	ctx.JSON(http.StatusOK, LoginResponse{Token: token})
}

// GuestLogin: ゲストユーザー作成・トークン発行API
func (c *UserController) GuestLogin(ctx *gin.Context) {
	guest := &user.User{
		IsGuest:  true,
		Nickname: "ゲスト",
		Email:    fmt.Sprintf("guest-%s@guest.local", uuid.New().String()),
	}
	if err := c.repo.Create(guest); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ゲストユーザー作成に失敗しました"})
		return
	}
	token, err := infra.GenerateToken(guest.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "トークン生成に失敗しました"})
		return
	}
	ctx.JSON(http.StatusOK, LoginResponse{Token: token})
}

// Register: ゲスト→正式ユーザー昇格対応
func (c *UserController) Register(ctx *gin.Context) {
	var req RegisterRequest
	var hash []byte
	var err error
	var newUser *user.User
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	// 既存ゲストユーザーの昇格を優先
	token := ctx.GetHeader("Authorization")
	var userID string
	if token != "" {
		token = strings.TrimPrefix(token, "Bearer ")
		uid, err := infra.ParseToken(token)
		if err == nil {
			userID = uid
		}
	}
	var u *user.User
	if userID != "" {
		u, err = c.repo.FindByID(userID)
		if err != nil || u == nil || !u.IsGuest {
			u = nil // ゲストでなければ新規作成
		}
	}
	if u != nil {
		// ゲスト昇格
		exists, _ := c.repo.FindByEmail(req.Email)
		if exists != nil {
			ctx.JSON(http.StatusConflict, gin.H{"error": "メールアドレスは既に登録されています"})
			return
		}
		hash, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードのハッシュ化に失敗しました"})
			return
		}
		u.Email = req.Email
		u.PasswordHash = string(hash)
		u.Nickname = req.Nickname
		u.IsGuest = false
		if err := c.repo.Update(u); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー昇格に失敗しました"})
			return
		}
		token, err := infra.GenerateToken(u.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "トークン生成に失敗しました"})
			return
		}
		ctx.JSON(http.StatusOK, RegisterResponse{Token: token})
		return
	}
	// 新規ユーザー作成
	exists, _ := c.repo.FindByEmail(req.Email)
	if exists != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "メールアドレスは既に登録されています"})
		return
	}
	hash, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードのハッシュ化に失敗しました"})
		return
	}
	newUser = &user.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Nickname:     req.Nickname,
		IsGuest:      false,
	}
	if err = c.repo.Create(newUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー登録に失敗しました"})
		return
	}
	token, err = infra.GenerateToken(newUser.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "トークン生成に失敗しました"})
		return
	}
	ctx.JSON(http.StatusOK, RegisterResponse{Token: token})
}
