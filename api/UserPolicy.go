package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserPolicyUri struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUserPolicyById(ctx *gin.Context) {
	var reqUri UserPolicyUri
	if err := ctx.BindUri(&reqUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userPolicy, err := server.store.GetUserPolicy(ctx, reqUri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, userPolicy)
}
func (server *Server) deleteUserPolicy(ctx *gin.Context) {
	var reqUri UserPolicyUri
	if err := ctx.BindUri(&reqUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := server.store.DeleteUserPolicy(ctx, reqUri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Policy deleted succesfully",
	})
}
