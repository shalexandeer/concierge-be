package utils

import (
	"errors"
	"time"

	"concierge-be/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	RoleID   string `json:"roleId"`
	RoleName string `json:"roleName"`
	TenantID string `json:"tenantId"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID string, username string) (string, error) {
	expireTime := time.Duration(config.AppConfig.JWT.ExpireTime) * time.Hour
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

// GenerateTokenWithRole 生成包含角色信息的 JWT Token
func GenerateTokenWithRole(userID string, username string, roleID string, roleName string, tenantID string) (string, error) {
	expireTime := time.Duration(config.AppConfig.JWT.ExpireTime) * time.Hour
	claims := Claims{
		UserID:   userID,
		Username: username,
		RoleID:   roleID,
		RoleName: roleName,
		TenantID: tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

// ParseToken 解析 JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ParseJWT is an alias for ParseToken for consistency with middleware
func ParseJWT(tokenString string) (*Claims, error) {
	return ParseToken(tokenString)
}
