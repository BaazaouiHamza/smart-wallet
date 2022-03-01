package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserPolicyUri struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getUserPolicyById(c *gin.Context) {
	var reqUri UserPolicyUri
	if err := c.BindUri(&reqUri); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO: insert
}

func (s *Server) deleteUserPolicy(c *gin.Context) {
	var reqUri UserPolicyUri
	if err := c.BindUri(&reqUri); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// TODO: delete
}
