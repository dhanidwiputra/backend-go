package handler

import (
	"final-project-backend/dto"
	"final-project-backend/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdateDelivery(c *gin.Context) {
	deliveryId := c.Param("id")
	deliveryIdInt, err := strconv.ParseInt(deliveryId, 10, 64)

	var deliveryRequest dto.DeliveryRequest
	if err := c.ShouldBindJSON(&deliveryRequest); err != nil {
		util.ResponseErrorJSON(c, err.Error(), "BAD_REQUEST", http.StatusBadRequest)
		return
	}

	delivery, err := h.deliveryUsecase.GetDeliveryById(uint(deliveryIdInt))
	if err != nil {
		util.ResponseErrorJSON(c, err.Error(), "DELIVERY_NOT_FOUND", http.StatusNotFound)
		return
	}

	delivery.Status = deliveryRequest.Status
	delivery, err = h.deliveryUsecase.UpdateDeliveryStatus(*delivery)
	if err != nil {
		util.ResponseErrorJSON(c, err.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, delivery, http.StatusOK)

}
