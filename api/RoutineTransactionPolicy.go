package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	prospercontext "git.digitus.me/library/prosper-kit/context"
	"git.digitus.me/library/prosper-kit/middleware"
	"git.digitus.me/pfe/smart-wallet/types"
	"git.digitus.me/prosperus/protocol/identity"
	ptclTypes "git.digitus.me/prosperus/protocol/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RTPRequest struct {
	Name              string             `json:"name" binding:"required"`
	Description       string             `json:"description" binding:"required"`
	NymID             identity.PublicKey `json:"nymId" binding:"required"`
	Recipient         identity.PublicKey `json:"recipient" binding:"required"`
	ScheduleStartDate time.Time          `json:"scheduleStartDate" binding:"required"`
	ScheduleEndDate   time.Time          `json:"scheduleEndDate" binding:"required"`
	Frequency         string             `json:"frequency" binding:"required,oneof= daily monthly weekly"`
	Amount            ptclTypes.Balance  `json:"amount" binding:"required"`
}

func (s *Server) addRoutineTransactionPolicy(c *gin.Context, pk identity.PublicKey) {
	logger := prospercontext.GetLogger(c)

	var req RTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// scheduleStartDate, _ := time.Parse("2006-01-02", req.ScheduleStartDate)
	// scheduleEndDate, _ := time.Parse("2006-01-02", req.ScheduleEndDate)
	if req.ScheduleEndDate.Before(req.ScheduleStartDate) {
		err := errors.New("scheduelEndDate must be after scheduelStartDate")
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO: insert
	err := s.service.CreateRoutineTransactionPolicy(prospercontext.JoinContexts(c), types.RoutineTransactionPolicy{
		Name:              req.Name,
		Description:       req.Description,
		ScheduleStartDate: req.ScheduleStartDate,
		ScheduleEndDate:   req.ScheduleEndDate,
		Amount:            req.Amount,
		Frequency:         req.Frequency,
		NymID:             pk,
		Recipient:         req.Recipient,
	})
	switch {
	case s.service.IsUserError(err):
		logger.Debug("invalid routine transaction policy", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid routine transaction policy"})
		return
	case err != nil:
		logger.Error("could not create routine transaction policy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	default:
		c.Status(http.StatusCreated)
		return
	}
}

type RTPRequestUri struct {
	ID int `uri:"id" binding:"required"`
}

// @ID delete-routine-transaction-policy
// @Tags routine-transaction-policy
// @Description Delete a routine transaction policy
// @Param nym-id path string true "NymID"
// @Param id path int true "ID"
// @Success 204
// @Router /api/:nym-id/routine-transaction-policy/:id [DELETE]
func (s *Server) deleteRoutineTransactionPolicy(c *gin.Context, pk identity.PublicKey) {
	var reqUri RTPRequestUri
	if err := c.BindUri(&reqUri); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := s.service.DeletePolicy(prospercontext.JoinContexts(c), pk, reqUri.ID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}

// @ID update-routine-transaction-policy
// @Tags routine-transaction-policy
// @Description Update a routine transaction policy
// @Param nym-id path string true "NymID"
// @Param id path int true "ID"
// @Param RTPRequest body RTPRequest true "Routine Transation Policy"
// @Success 200
// @Router /api/:nym-id/routine-transaction-policy/:id [PUT]
func (s *Server) updateRoutineTransactionPolicy(c *gin.Context, pk identity.PublicKey) {

	logger := prospercontext.GetLogger(c)
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
	// scheduleStartDate, _ := time.Parse("2006-01-02", req.ScheduleStartDate)
	// scheduleEndDate, _ := time.Parse("2006-01-02", req.ScheduleEndDate)
	if req.ScheduleEndDate.Before(req.ScheduleStartDate) {
		err := errors.New("scheduelEndDate must be after scheduelStartDate")
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.service.UpdateRoutineTransactionPolicy(prospercontext.JoinContexts(c), types.RoutineTransactionPolicy{
		ID:                int(reqUri.ID),
		Name:              req.Name,
		Description:       req.Description,
		Recipient:         req.Recipient,
		NymID:             pk,
		ScheduleStartDate: req.ScheduleStartDate,
		ScheduleEndDate:   req.ScheduleEndDate,
		Frequency:         req.Frequency,
		Amount:            req.Amount,
	})
	switch {
	case s.service.IsUserError(err):
		logger.Debug("invalid routine transaction policy", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid routine transaction policy"})
		return
	case err != nil:
		logger.Error("could not update routine transaction policy", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	default:
		c.Status(http.StatusOK)
		return
	}

}

// @ID get-routine-transaction-policy
// @Tags routine-transaction-policy
// @Description Get a routine transaction policy
// @Param nym-id path string true "NymID"
// @Param id path int true "ID"
// @Success 200 {object} types.RoutineTransactionPolicy
// @Router /api/:nym-id/routine-transaction-policy/:id [GET]
func (s *Server) getRoutineTransactionPolicyById(c *gin.Context, pk identity.PublicKey) {
	var reqUri RTPRequestUri

	if err := c.ShouldBindUri(&reqUri); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rtp, err := s.service.GetRoutineTransactionPolicy(prospercontext.JoinContexts(c), pk, int(reqUri.ID))
	if err != nil {
		if s.service.IsNotFoundError(err) {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, rtp)
}

type RTPNymUri struct {
	NymID string `uri:"nym_id" binding:"required"`
}

type listRoutineTransactionPolicies struct {
	Page         int `form:"page" binding:"min=0"`
	ItemsPerPage int `form:"itemsPerPage" binding:"required,min=5,max=10"`
}

type listRoutineTransactionPoliciesResponse struct {
	Data  []types.RoutineTransactionPolicy `json:"data"`
	Total int                              `json:"total"`
}

// @ID list-routine-transaction-policy
// @Tags routine-transaction-policy
// @Description Get all transaction trigger policies
// @Param nym-id path string true "NymID"
// @Param _ query listRoutineTransactionPoliciesResponse false "comment"
// @Success 200 {object} listRoutineTransactionPoliciesResponse
// @Router /api/:nym-id/routine-transaction-policy [GET]
func (s *Server) listRoutineTransactionPolicies(c *gin.Context, pk identity.PublicKey) {
	fmt.Println(c.Request.Header)
	logger := middleware.GetLogger(c)
	logger.Debug("Authorization", zap.String("", c.Request.Header.Get("Authorization")))
	var reqForm listRoutineTransactionPolicies
	if err := c.ShouldBindQuery(&reqForm); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rtps, total, err := s.service.ListRoutineTransactionPolicies(
		prospercontext.JoinContexts(c), pk, reqForm.Page, reqForm.ItemsPerPage,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, listRoutineTransactionPoliciesResponse{
		Data: rtps,
		// TODO: return total
		Total: total,
	})
}
