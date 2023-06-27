package middleware

import (
	"final-project-backend/domain"
	"final-project-backend/dto"
	"final-project-backend/util"

	"github.com/gin-gonic/gin"
)

func AuthorizeAdmin(c *gin.Context) {
	user := c.MustGet("user")
	role := user.(dto.UserResponse).Role
	if role != "admin" {
		util.ResponseErrorJSON(c, domain.ErrForbiddenAccess.Error(), "forbidden access", 403)
		c.Abort()
		return
	}
	c.Next()
}
