package infra

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndParseToken_OK(t *testing.T) {
	userID := int64(123)
	token, err := GenerateToken(userID)
	assert.NoError(t, err)
	parsedID, err := ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedID)
}

func TestParseToken_InvalidToken(t *testing.T) {
	_, err := ParseToken("invalid.token.value")
	assert.Error(t, err)
}

func TestParseToken_Expired(t *testing.T) {
	// 有効期限を過去にしたトークンを生成
	// ここではGenerateTokenのexpを一時的に過去にして生成する方法が必要
	// 省略: 実装例ではexpを直接書き換えたトークンを用意するのが一般的
}
