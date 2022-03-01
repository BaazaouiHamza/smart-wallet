package api

import (
	"net/http"

	prospercontext "git.digitus.me/library/prosper-kit/context"
	"git.digitus.me/pfe/smart-wallet/types"
	"git.digitus.me/prosperus/protocol/identity"
	ptclTypes "git.digitus.me/prosperus/protocol/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TTPRequest struct {
	Name            string             `json:"name" binding:"required"`
	Description     string             `json:"description" binding:"required"`
	TargetedBalance ptclTypes.Balance  `json:"targetedBalance" binding:"required"`
	Recipient       identity.PublicKey `json:"recipient" binding:"required"`
	Amount          ptclTypes.Balance  `json:"amount" binding:"required,min=1"`
}

// @ID create-transaction-trigger-policy
// @Tags transaction-trigger-policy
// @Description Create a transaction trigger policy
// @Param nym-id path string true "NymID"
// @Param TTPRequest body TTPRequest true "Transation Trigger Policy"
// @Success 201
// @Router /api/:nym-id/transaction-trigger-policy [POST]
func (s *Server) createTransactionTriggerPolicy(c *gin.Context, pk identity.PublicKey) {
	logger := prospercontext.GetLogger(c)

	var req TTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}

	err := s.service.CreateTransactionTriggerPolicy(
		prospercontext.JoinContexts(c),
		types.TransactionTriggerPolicy{
			Name:            req.Name,
			Description:     req.Description,
			NymID:           pk,
			TargetedBalance: req.TargetedBalance,
			Recipient:       req.Recipient,
			Amount:          req.Amount,
		},
	)
	switch {
	case s.service.IsUserError(err):
		logger.Debug("invalid transaction trigger policy", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid transaction trigger policy"})
	case err != nil:
		logger.Error("could not create transaction trigger policy", zap.Error(err))
		c.Status(http.StatusInternalServerError)
	default:
		c.Status(http.StatusCreated)
	}
}

type TTPRequestUri struct {
	ID int `uri:"id" binding:"required,min=1"`
}

func (s *Server) updateTransactionTriggerPolicy(c *gin.Context, pk identity.PublicKey) {
	var req TTPRequest
	var reqUri TTPRequestUri

	if err := c.BindUri(&reqUri); err != nil {
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}

	arg := types.TransactionTriggerPolicy{
		ID:              reqUri.ID,
		Name:            req.Name,
		Description:     req.Description,
		NymID:           pk,
		Amount:          req.Amount,
		TargetedBalance: req.TargetedBalance,
	}

	err := s.service.UpdateTransactionTriggerPolicy(prospercontext.JoinContexts(c), arg)
	if err != nil {
		return
	}

	c.Status(http.StatusNoContent)
}

func (s *Server) deleteTransactionTriggerPolicy(c *gin.Context, pk identity.PublicKey) {
	var reqUri TTPRequestUri
	if err := c.BindUri(&reqUri); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := s.service.DeletePolicy(prospercontext.JoinContexts(c), pk, reqUri.ID); err != nil {
		return
	}

	c.Status(http.StatusNoContent)
}

func (s *Server) getTransactionTriggerPolicyById(c *gin.Context, pk identity.PublicKey) {
	var reqUri TTPRequestUri

	if err := c.ShouldBindUri(&reqUri); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ttp, err := s.service.GetTransactionTriggerPolicy(
		prospercontext.JoinContexts(c), pk, reqUri.ID,
	)
	if err != nil {
		if s.service.IsNotFoundError(err) {
			c.Status(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, ttp)
}

type TTPNymUri struct {
	NymID identity.PublicKey `uri:"nym_id" binding:"required"`
}

type listTransactionTriggerPolicies struct {
	Page         int `form:"page" binding:"required,min=1"`
	ItemsPerPage int `form:"itemsPerPage" binding:"required,min=5,max=10"`
}

func (s *Server) listTransactionTriggerPolicies(c *gin.Context, pk identity.PublicKey) {
	var reqForm listTransactionTriggerPolicies
	if err := c.ShouldBindQuery(&reqForm); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	rtps, count, err := s.service.ListTransactionTriggerPolicies(
		prospercontext.JoinContexts(c), pk, reqForm.Page, reqForm.ItemsPerPage,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": rtps, "total": count})
}
