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

	{
		rtpRouter := router.Group("/:nymID/routine-transaction-policy")
		rtpRouter.POST("", checkContributorNymID(server.addRoutineTransactionPolicy))
		rtpRouter.GET("", checkViewerNymID(server.listRoutineTransactionPolicies))
		rtpRouter.PUT("/:id", checkContributorNymID(server.updateRoutineTransactionPolicy))
		rtpRouter.GET("/:id", checkViewerNymID(server.getRoutineTransactionPolicyById))
		rtpRouter.DELETE("/:id", checkContributorNymID(server.deleteRoutineTransactionPolicy))
	}
	{
		ttRouter := router.Group("/:nymID/transaction-trigger-policy")

		//Transaction Trigger Policy ROUTER
		ttRouter.POST("", checkContributorNymID(server.createTransactionTriggerPolicy))
		ttRouter.PUT(
			"/:id",
			checkContributorNymID(server.updateTransactionTriggerPolicy),
		)
		ttRouter.GET(
			"/:id",
			checkViewerNymID(server.getTransactionTriggerPolicyById),
		)
		ttRouter.DELETE(
			"/:id",
			checkContributorNymID(server.deleteTransactionTriggerPolicy),
		)
		ttRouter.GET(
			"",
			checkViewerNymID(server.listTransactionTriggerPolicies),
		)
	}

}

func errorResponse(err error) gin.H {
	return gin.H{"message": err.Error()}
}
