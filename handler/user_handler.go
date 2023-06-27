package handler

import (
	"final-project-backend/domain"
	"final-project-backend/dto"
	"final-project-backend/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (h *Handler) GetUserDetails(c *gin.Context) {
	user, ok := c.Get("user")

	if !ok {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", 401)
		return
	}

	userID := user.(dto.UserResponse).ID
	userDetail, err := h.userUsecase.GetUserByID(userID)

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrUserNotFound.Error(), "USER_NOT_FOUND", 404)
		return
	}

	util.ResponseSuccesJSON(c, userDetail, 200)
}

func (h *Handler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	intUserId, err := strconv.Atoi(userID)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", 400)
		return
	}

	userDetail, err := h.userUsecase.GetUserByID(uint(intUserId))

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrUserNotFound.Error(), "USER_NOT_FOUND", 404)
		return
	}

	util.ResponseSuccesJSON(c, userDetail, 200)
}

func (h *Handler) UpdateUserDetails(c *gin.Context) {
	user, ok := c.Get("user")

	if !ok {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", 401)
		return
	}

	userID := user.(dto.UserResponse).ID
	userData, err := h.userUsecase.GetUserByID(userID)

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrUserNotFound.Error(), "USER_NOT_FOUND", 404)
		return
	}

	updateUserRequestBody := dto.UpdateUserRequest{}
	err = c.MustBindWith(&updateUserRequestBody, binding.FormMultipart)

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", 400)
		return
	}

	isPhone, _ := util.IsPhone(updateUserRequestBody.Phone)
	if !isPhone {
		util.ResponseErrorJSON(c, domain.ErrInvalidPhoneFormat.Error(), "INVALID_PHONE_FORMAT", 400)
		return
	}

	if (userData.PictureUrl != "") && (updateUserRequestBody.ProfilePicture.Size != 0) {
		err = h.mediaUsecase.FileDelete(userData.PicturePublicId)
	}

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrDeleteImage.Error(), "DELETE_IMAGE_FAILED", 400)
		return
	}

	userData.Phone = updateUserRequestBody.Phone
	userData.FullName = updateUserRequestBody.FullName

	var picUrl, publicId string
	if updateUserRequestBody.ProfilePicture.Size != 0 {
		picUrl, publicId, err = h.mediaUsecase.FileUpload(updateUserRequestBody.ProfilePicture)
		userData.PictureUrl = picUrl
		userData.PicturePublicId = publicId
	}

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrUploadImage.Error(), "UPLOAD_IMAGE_FAILED", 400)
		return
	}

	userRes, err := h.userUsecase.UpdateUser(*userData)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrUserNotFound.Error(), "USER_NOT_FOUND", 404)
		return
	}

	userDtoResponse := dto.UserResponse{
		ID:         userRes.ID,
		FullName:   userRes.FullName,
		Email:      userRes.Email,
		Phone:      userRes.Phone,
		PictureUrl: userRes.PictureUrl,
		Role:       userRes.Role,
	}

	util.ResponseSuccesJSON(c, userDtoResponse, 200)
}

func (h *Handler) DeleteUserPhoto(c *gin.Context) {
	user, ok := c.Get("user")

	if !ok {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", 401)
		return
	}

	userID := user.(dto.UserResponse).ID

	err := h.userUsecase.DeleteUserPhoto(userID)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrUserNotFound.Error(), "USER_NOT_FOUND", 404)
		return
	}

	util.ResponseSuccesJSON(c, nil, http.StatusNoContent)
}

func (h *Handler) ResetGamesAttempt(c *gin.Context) {
	err := h.userUsecase.ResetGamesAttempt()
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", 500)
		return
	}

	util.ResponseSuccesJSON(c, "Games Attempt Reset", http.StatusAccepted)
}
