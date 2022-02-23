package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "git.digitus.me/pfe/smart-wallet/db/sqlc"
	"github.com/gin-gonic/gin"
)

type RTPRequest struct {
	Name              string `json:"name" binding:"required"`
	Description       string `json:"description" binding:"required"`
	NymID             string `json:"nym_id" binding:"required"`
	ScheduleStartDate string `json:"schedule_start_date" binding:"required"`
	ScheduleEndDate   string `json:"schedule_end_date" binding:"required"`
	Frequency         string `json:"frequency" binding:"required,oneof= DAILY MONTHLY WEEKLY"`
	Amount            int32  `json:"amount" binding:"required"`
}

func (server *Server) createRoutineTransactionPolicy(ctx *gin.Context) {
	var req RTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	scheduleStartDate, _ := time.Parse("2006-01-02", req.ScheduleStartDate)
	scheduleEndDate, _ := time.Parse("2006-01-02", req.ScheduleEndDate)
	if scheduleEndDate.Before(scheduleStartDate) {
		err := errors.New("scheduelEndDate must be after scheduelStartDate")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateRoutineTransactionPolicyParams{
		Name:              req.Name,
		Description:       req.Description,
		NymID:             req.NymID,
		ScheduleStartDate: scheduleStartDate,
		ScheduleEndDate:   scheduleEndDate,
		Frequency:         req.Frequency,
		Amount:            req.Amount,
	}
	rtp, err := server.store.AddRoutineTransactionPolicy(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, rtp)
}

type RTPRequestUri struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteRoutineTransactionPolicy(ctx *gin.Context) {
	var reqUri RTPRequestUri
	if err := ctx.BindUri(&reqUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := server.store.DeleteRoutineTransactionPolicy(ctx, reqUri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "routine transaction policy deleted successfully",
	})
}

func (server *Server) updateRoutineTransactionPolicy(ctx *gin.Context) {
	var req RTPRequest
	var reqUri RTPRequestUri
	if err := ctx.BindUri(&reqUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	scheduleStartDate, _ := time.Parse("2006-01-02", req.ScheduleStartDate)
	scheduleEndDate, _ := time.Parse("2006-01-02", req.ScheduleEndDate)
	if scheduleEndDate.Before(scheduleStartDate) {
		err := errors.New("scheduelEndDate must be after scheduelStartDate")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.UpdateRoutineTransactionPolicyParams{
		ID:                reqUri.ID,
		Name:              req.Name,
		Description:       req.Description,
		NymID:             req.NymID,
		ScheduleStartDate: scheduleStartDate,
		ScheduleEndDate:   scheduleEndDate,
		Frequency:         req.Frequency,
		Amount:            req.Amount,
	}
	rtp, err := server.store.UpdateRoutineTransactionPolicy(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rtp)
}

func (server *Server) getRoutineTransactionPolicyById(ctx *gin.Context) {
	var reqUri RTPRequestUri

	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rtp, err := server.store.GetRoutineTransactionPolicy(ctx, reqUri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rtp)
}

type RTPNymUri struct {
	NymID string `uri:"nym_id" binding:"required"`
}

type listRoutineTransactionPolicies struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listRoutineTransactionPolicies(ctx *gin.Context) {
	var reqUri RTPNymUri
	var reqForm listRoutineTransactionPolicies
	if err := ctx.ShouldBindQuery(&reqForm); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.BindUri(&reqUri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.ListRoutineTransactionPoliciesParams{
		NymID:  reqUri.NymID,
		Limit:  reqForm.PageSize,
		Offset: (reqForm.PageId - 1) * reqForm.PageSize,
	}
	rtps, err := server.store.ListRoutineTransactionPolicies(ctx, arg)
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
