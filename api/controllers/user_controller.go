package controllers

import (
	"feelog-backend/domain/user"
	"feelog-backend/infra"
	"net/http"

	"github.com/gin-gonic/gin"
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

// Register: ユーザー登録API
func (c *UserController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	// メール重複チェック
	exists, _ := c.repo.FindByEmail(req.Email)
	if exists != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "メールアドレスは既に登録されています"})
		return
	}
	// パスワードハッシュ化
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードのハッシュ化に失敗しました"})
		return
	}
	user := &user.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Nickname:     req.Nickname,
	}
	if err := c.repo.Create(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー登録に失敗しました"})
		return
	}
	token, err := infra.GenerateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "トークン生成に失敗しました"})
		return
	}
	ctx.JSON(http.StatusOK, RegisterResponse{Token: token})
}
