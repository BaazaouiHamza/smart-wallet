package api

import (
	"errors"
	"net/http"
	"time"

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

func (s *Server) createRoutineTransactionPolicy(c *gin.Context) {
	var req RTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	scheduleStartDate, _ := time.Parse("2006-01-02", req.ScheduleStartDate)
	scheduleEndDate, _ := time.Parse("2006-01-02", req.ScheduleEndDate)
	if scheduleEndDate.Before(scheduleStartDate) {
		err := errors.New("scheduelEndDate must be after scheduelStartDate")
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO: insert
}

type RTPRequestUri struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) deleteRoutineTransactionPolicy(c *gin.Context) {
	var reqUri RTPRequestUri
	if err := c.BindUri(&reqUri); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO: delete
}

func (s *Server) updateRoutineTransactionPolicy(c *gin.Context) {
	var req RTPRequest
	var reqUri RTPRequestUri
	if err := c.BindUri(&reqUri); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	scheduleStartDate, _ := time.Parse("2006-01-02", req.ScheduleStartDate)
	scheduleEndDate, _ := time.Parse("2006-01-02", req.ScheduleEndDate)
	if scheduleEndDate.Before(scheduleStartDate) {
		err := errors.New("scheduelEndDate must be after scheduelStartDate")
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO: update
}

func (s *Server) getRoutineTransactionPolicyById(c *gin.Context) {
	var reqUri RTPRequestUri

	if err := c.ShouldBindUri(&reqUri); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO: Get
}

type RTPNymUri struct {
	NymID string `uri:"nym_id" binding:"required"`
}

type listRoutineTransactionPolicies struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) listRoutineTransactionPolicies(c *gin.Context) {
	var reqUri RTPNymUri
	var reqForm listRoutineTransactionPolicies
	if err := c.ShouldBindQuery(&reqForm); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := c.BindUri(&reqUri); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO: List
}
