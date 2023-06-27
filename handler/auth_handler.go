package handler

import (
	"final-project-backend/domain"
	"final-project-backend/dto"
	"strings"

	"errors"
	"final-project-backend/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (h *Handler) Login(c *gin.Context) {
	loginRequestBody := dto.LoginRequest{}
	err := c.BindJSON(&loginRequestBody)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", 400)
		return
	}
	token, err := h.authUsecase.Login(loginRequestBody)
	if errors.Is(err, domain.ErrInvalidEmail) {
		util.ResponseErrorJSON(c, domain.ErrInvalidEmail.Error(), "INVALID_EMAIL", http.StatusBadRequest)
		return
	}

	if errors.Is(err, domain.ErrInvalidPassword) {
		util.ResponseErrorJSON(c, domain.ErrInvalidPassword.Error(), "INVALID_PASSWORD", http.StatusBadRequest)
		return
	}

	if errors.Is(err, domain.ErrInvalidUsername) {
		util.ResponseErrorJSON(c, domain.ErrInvalidUsername.Error(), "INVALID_USERNAME", http.StatusBadRequest)
		return
	}

	if errors.Is(err, domain.ErrInvalidPhone) {
		util.ResponseErrorJSON(c, domain.ErrInvalidPhone.Error(), "INVALID_PHONE", http.StatusBadRequest)
		return
	}

	if errors.Is(err, domain.ErrInvalidIdentifier) {
		util.ResponseErrorJSON(c, domain.ErrInvalidIdentifier.Error(), "INVALID_IDENTIFIER", http.StatusBadRequest)
		return
	}

	if errors.Is(err, domain.ErrInternalServer) {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	token_json := dto.JWTAuthenticationResponse{
		Token: token,
	}

	util.ResponseSuccesJSON(c, token_json, http.StatusOK)
}

func (h *Handler) Register(c *gin.Context) {
	registerRequestBody := dto.RegisterRequest{}
	err := c.MustBindWith(&registerRequestBody, binding.FormMultipart)

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", 400)
		return
	}

	registerData := dto.RegisterData{
		FullName:    registerRequestBody.FullName,
		PhoneNumber: registerRequestBody.Phone,
		Username:    registerRequestBody.Username,
		Email:       registerRequestBody.Email,
		Password:    registerRequestBody.Password,
	}

	if registerRequestBody.ProfilePicture.Size != 0 {
		picUrl, publicId, _ := h.mediaUsecase.FileUpload(registerRequestBody.ProfilePicture)
		registerData.ProfilePicture = picUrl
		registerData.PublicId = publicId
	}

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	user, err := h.authUsecase.Register(registerData)

	if errors.Is(err, domain.ErrDuplicateEmail) {
		util.ResponseErrorJSON(c, domain.ErrDuplicateEmail.Error(), "DUPLICATE_EMAIL", http.StatusBadRequest)
		return
	}

	if errors.Is(err, domain.ErrDuplicateUsername) {
		util.ResponseErrorJSON(c, domain.ErrDuplicateUsername.Error(), "DUPLICATE_USERNAME", http.StatusBadRequest)
		return
	}

	if errors.Is(err, domain.ErrDuplicatePhone) {
		util.ResponseErrorJSON(c, domain.ErrDuplicatePhone.Error(), "DUPLICATE_PHONE", http.StatusBadRequest)
		return
	}

	if errors.Is(err, domain.ErrInvalidEmailFormat) {
		util.ResponseErrorJSON(c, domain.ErrInvalidEmailFormat.Error(), "INVALID_EMAIL_FORMAT", http.StatusBadRequest)
		return
	}

	if errors.Is(err, domain.ErrInvalidPhoneFormat) {
		util.ResponseErrorJSON(c, domain.ErrInvalidPhoneFormat.Error(), "INVALID_PHONE_FORMAT", http.StatusBadRequest)
		return
	}

	if errors.Is(err, domain.ErrInvalidUsernameFormat) {
		util.ResponseErrorJSON(c, domain.ErrInvalidUsernameFormat.Error(), "INVALID_USERNAME_FORMAT", http.StatusBadRequest)
		return
	}

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", 500)
		return
	}

	util.ResponseSuccesJSON(c, user, 201)
}

func (h *Handler) HasValidToken(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	s := strings.Split(token, "Bearer ")
	decodedToken := s[1]
	user := c.MustGet("user")
	userId := user.(dto.UserResponse).ID

	isValid := h.authUsecase.HasValidToken(userId, decodedToken)
	if !isValid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}
