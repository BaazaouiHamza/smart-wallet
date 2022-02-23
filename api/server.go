package api

import (
	db "git.digitus.me/pfe/smart-wallet/db/sqlc"
	"git.digitus.me/pfe/smart-wallet/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config util.Config
	router *gin.Engine
	store  db.Store
}

func NewServer(config util.Config, db db.Store) (*Server, error) {
	server := &Server{
		config: config,
		store:  db,
	}
	server.setUpRouter()
	return server, nil
}
func (server *Server) setUpRouter() {
	router := gin.Default()

	//User Policy Router
	router.GET("/userPolicy/:id", server.getUserPolicyById)
	router.DELETE("/userPolicy/:id", server.deleteUserPolicy)

	//Routine Transaction Policy ROUTER
	router.POST("/policy/routineTransactionPolicy", server.createRoutineTransactionPolicy)
	router.PATCH("/policy/routineTransactionPolicy/:id", server.updateRoutineTransactionPolicy)
	router.DELETE("/policy/routineTransactionPolicy/:id", server.deleteRoutineTransactionPolicy)
	router.GET("/policy/routineTransactionPolicy/wallet/:nym_id", server.listRoutineTransactionPolicies)
	router.GET("/policy/routineTransactionPolicy/:id", server.getRoutineTransactionPolicyById)

	//Transaction Trigger Policy ROUTER
	router.POST("/policy/transactionTriggerPolicy", server.createTransactionTriggerPolicy)
	router.PUT("/policy/transactionTriggerPolicy/:id", server.updateTransactionTriggerPolicy)
	router.GET("/policy/transactionTriggerPolicy/:id", server.getTransactionTriggerPolicyById)
	router.DELETE("/policy/transactionTriggerPolicy/:id", server.deleteTransactionTriggerPolicy)
	router.GET("/policy/transactionTriggerPolicy/wallet/:nym_id", server.listTransactionTriggerPolicies)

	server.router = router
}

//Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
