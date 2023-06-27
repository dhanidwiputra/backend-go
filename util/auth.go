package util

import (
	"final-project-backend/config"
	"final-project-backend/dto"
	"final-project-backend/entity"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthUtil interface {
	GenerateAccessToken(user *entity.User) (string, error)
	ComparePassword(hashedPwd string, inputPwd string) bool
	HashPassword(password string) (string, error)
}

type authUtilImpl struct{}

func NewAuthUtil() AuthUtil {
	return &authUtilImpl{}
}

type customClaims struct {
	jwt.RegisteredClaims
	User *dto.UserResponse `json:"user"`
}

func (a *authUtilImpl) GenerateAccessToken(user *entity.User) (string, error) {
	c := config.InitConfig().JWTConfig
	userDTO := &dto.UserResponse{
		ID:   user.ID,
		Role: user.Role.Name,
	}

	expiredTime, _ := strconv.Atoi(c.ExpTimeMinutes)
	dur := time.Duration(expiredTime) * time.Minute

	claims := &customClaims{
		User: userDTO,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(dur)),
			Issuer:    c.JWTIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	hmacSampleSecret := c.SecretString

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *authUtilImpl) ComparePassword(hashedPwd string, inputPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(inputPwd))
	return err == nil
}

func (a *authUtilImpl) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}
