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
	service service.SmartWallet
}

func NewServer(
	svc service.SmartWallet,
	engine *gin.Engine,
	jwsGetter middleware.JWSGetter,
) *Server {
	server := &Server{service: svc}
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
	// router.Use(middleware.WithAuthentication(nil))

	//User Policy Router
	router.GET("/user-policy/:id", server.getUserPolicyById)
	router.DELETE("/user-policy/:id", server.deleteUserPolicy)

	{
		rtRouter := router.Group("/:nym-id/routine-transaction-policy")

		//Routine Transaction Policy ROUTER
		rtRouter.POST("", server.createRoutineTransactionPolicy)
		rtRouter.PATCH("/:id", server.updateRoutineTransactionPolicy)
		rtRouter.DELETE("/:id", server.deleteRoutineTransactionPolicy)
		rtRouter.GET(
			"/wallet/:nym-id",
			server.listRoutineTransactionPolicies,
		)
		rtRouter.GET("/:id", server.getRoutineTransactionPolicyById)
	}

	{
		ttRouter := router.Group("/:nym-id/transaction-trigger-policy")

		//Transaction Trigger Policy ROUTER
		ttRouter.POST(
			"",
			checkNymID(middleware.ContributorPermissionLevel, server.createTransactionTriggerPolicy),
		)
		ttRouter.PUT(
			"/:id",
			checkNymID(middleware.ContributorPermissionLevel, server.updateTransactionTriggerPolicy),
		)
		ttRouter.GET(
			"/:id",
			checkNymID(middleware.ViewerPermissionLevel, server.getTransactionTriggerPolicyById),
		)
		ttRouter.DELETE(
			"/:id",
			checkNymID(middleware.ContributorPermissionLevel, server.deleteTransactionTriggerPolicy),
		)
		ttRouter.GET(
			"",
			checkNymID(middleware.ViewerPermissionLevel, server.listTransactionTriggerPolicies),
		)
	}

}

func errorResponse(err error) gin.H {
	return gin.H{"message": err.Error()}
}
