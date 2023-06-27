package middleware

import (
	"encoding/json"
	"final-project-backend/config"
	"final-project-backend/domain"
	"final-project-backend/dto"
	"final-project-backend/util"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func validateToken(encodedToken string) (*jwt.Token, error) {
	c := config.InitConfig().JWTConfig
	return jwt.Parse(encodedToken, func(t *jwt.Token) (interface{}, error) {
		if _, isValid := t.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, domain.ErrUnauthorized
		}
		return []byte(c.SecretString), nil
	})
}

func Authorize(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	s := strings.Split(authHeader, "Bearer ")

	if len(s) < 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	decodedToken := s[1]
	token, err := validateToken(decodedToken)

	if err != nil || !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	userJson, _ := json.Marshal(claims["user"])
	var user dto.UserResponse
	err = json.Unmarshal(userJson, &user)

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "unauthorized", 401)
		return
	}

	c.Set("user", user)
}
