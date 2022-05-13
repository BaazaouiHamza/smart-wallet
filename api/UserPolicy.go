package api

import (
	"net/http"

	prospercontext "git.digitus.me/library/prosper-kit/context"
	"github.com/gin-gonic/gin"
)

type UserPolicyUri struct {
	ID int `uri:"id" binding:"required,min=1"`
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
	if err := s.service.DeleteUserPolicy(prospercontext.JoinContexts(c), reqUri.ID); err != nil {
		return
	}
	c.Status(http.StatusOK)
}
