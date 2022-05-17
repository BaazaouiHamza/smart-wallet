package api

import (
	"git.digitus.me/library/prosper-kit/middleware"
	"git.digitus.me/pfe/smart-wallet/service"
	"github.com/gin-gonic/gin"

	_ "git.digitus.me/pfe/smart-wallet/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @Title Smart Wallet
// @Version 0.0.0
// @Description ProsperUs Smart wallet

type Server struct {
	service   service.SmartWallet
	jwsGetter middleware.JWSGetter
}

func NewServer(
	svc service.SmartWallet,
	engine *gin.Engine,
	jwsGetter middleware.JWSGetter,
) *Server {
	server := &Server{service: svc, jwsGetter: jwsGetter}
	server.setUpRouter(engine)
	return server
}

func (server *Server) setUpRouter(engine *gin.Engine) {

	{
		url := ginSwagger.URL("swagger/doc.json")
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

	// TODO: get actual JWSGetter
	router := engine.Group("/api")
	router.Use(middleware.WithAuthentication(server.jwsGetter))

	//User Policy Router
	router.GET("/user-policy/:id", server.getUserPolicyById)
	router.DELETE("/user-policy/:id", server.deleteUserPolicy)

	//Routine Transaction Policy ROUTER
	router.POST("/policy/routineTransactionPolicy", server.createRoutineTransactionPolicy)
	router.PUT("/policy/routineTransactionPolicy/:id", server.updateRoutineTransactionPolicy)
	router.DELETE("/policy/routineTransactionPolicy/:id", server.deleteRoutineTransactionPolicy)
	router.GET("/policy/routineTransactionPolicy/wallet/:nym_id", server.listRoutineTransactionPolicies)
	router.GET("/policy/routineTransactionPolicy/:id", server.getRoutineTransactionPolicyById)
	{
		rtpRouter := router.Group("/:nymID/routine-transaction-policy")
		rtpRouter.POST("", checkContributorNymID(middleware.ContributorPermissionLevel, server.addRoutineTransactionPolicy))
	}
	{
		ttRouter := router.Group("/:nymID/transaction-trigger-policy")

		//Transaction Trigger Policy ROUTER
		ttRouter.POST("", checkContributorNymID(middleware.ContributorPermissionLevel, server.createTransactionTriggerPolicy))
		ttRouter.PUT(
			"/:id",
			checkContributorNymID(middleware.ContributorPermissionLevel, server.updateTransactionTriggerPolicy),
		)
		ttRouter.GET(
			"/:id",
			checkViewerNymID(middleware.ViewerPermissionLevel, server.getTransactionTriggerPolicyById),
		)
		ttRouter.DELETE(
			"/:id",
			checkContributorNymID(middleware.ContributorPermissionLevel, server.deleteTransactionTriggerPolicy),
		)
		ttRouter.GET(
			"",
			checkViewerNymID(middleware.ViewerPermissionLevel, server.listTransactionTriggerPolicies),
		)
	}

}

func errorResponse(err error) gin.H {
	return gin.H{"message": err.Error()}
}
