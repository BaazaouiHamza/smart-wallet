package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	db "git.digitus.me/pfe/smart-wallet/db/sqlc"
	"github.com/gin-gonic/gin"
)

type TTPRequest struct {
	Name            string          `json:"name" binding:"required"`
	Description     string          `json:"description" binding:"required"`
	NymID           string          `json:"nym_id" binding:"required"`
	TargetedBalance json.RawMessage `json:"targeted_balance" binding:"required"`
	Amount          int32           `json:"amount" binding:"required,min=1"`
}

func (server *Server) createTransactionTriggerPolicy(ctx *gin.Context) {
	var req TTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateTransactionTriggerPolicyParams{
		Name:            req.Name,
		Description:     req.Description,
		NymID:           req.NymID,
		TargetedBalance: req.TargetedBalance,
		Amount:          req.Amount,
	}
	ttp, err := server.store.CreateTransactionTriggerPolicy(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, ttp)
}

type TTPRequestUri struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) updateTransactionTriggerPolicy(ctx *gin.Context) {
	var req TTPRequest
	var reqUri TTPRequestUri
	if err := ctx.BindUri(&reqUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.UpdateTransactionTriggerPolicyParams{
		ID:              reqUri.ID,
		Name:            req.Name,
		Description:     req.Description,
		NymID:           req.NymID,
		Amount:          req.Amount,
		TargetedBalance: req.TargetedBalance,
	}

	ttp, err := server.store.UpdateTransactionTriggerPolicy(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ttp)
}

func (server *Server) deleteTransactionTriggerPolicy(ctx *gin.Context) {
	var reqUri TTPRequestUri
	if err := ctx.BindUri(&reqUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := server.store.DeleteTransactionTriggerPolicy(ctx, reqUri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "policy deleted succesfully",
	})
}

func (server *Server) getTransactionTriggerPolicyById(ctx *gin.Context) {
	var reqUri TTPRequestUri

	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ttp, err := server.store.GetTransactionTriggerPolicy(ctx, reqUri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ttp)
}

type TTPNymUri struct {
	NymID string `uri:"nym_id" binding:"required"`
}

type listTransactionTriggerPolicies struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listTransactionTriggerPolicies(ctx *gin.Context) {
	var reqUri TTPNymUri
	var reqForm listTransactionTriggerPolicies
	if err := ctx.ShouldBindQuery(&reqForm); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.BindUri(&reqUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.ListTransactionTriggerPoliciesParams{
		NymID:  reqUri.NymID,
		Limit:  reqForm.PageSize,
		Offset: (reqForm.PageId - 1) * reqForm.PageSize,
	}
	rtps, err := server.store.ListTransactionTriggerPolicies(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, rtps)
}
