package middleware

import (
	"feelog-backend/infra"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthMiddleware_TableDriven(t *testing.T) {
	gin.SetMode(gin.TestMode)
	token, _ := infra.GenerateToken(42)

	tests := []struct {
		name       string
		header     string
		wantStatus int
		comment    string
	}{
		{"正常系: 有効なトークン", "Bearer " + token, 200, "有効なJWTで認証通過"},
		{"異常系: トークンなし", "", 401, "トークンなしで401"},
		{"異常系: 無効なトークン", "Bearer invalidtoken", 401, "不正なトークンで401"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(JWTAuthMiddleware())
			r.GET("/protected", func(c *gin.Context) {
				userID, exists := c.Get("userID")
				if !exists {
					c.JSON(500, gin.H{"error": "userID not set"})
					return
				}
				c.JSON(200, gin.H{"userID": userID})
			})
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/protected", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code, tt.comment)
		})
	}
}
