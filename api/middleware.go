package api

import (
	"net/http"

	"git.digitus.me/library/prosper-kit/middleware"
	"git.digitus.me/prosperus/protocol/identity"
	"github.com/gin-gonic/gin"
)

func checkContributorNymID(
	hn func(*gin.Context, identity.PublicKey),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		pk, err := identity.PublicKeyFromString(c.Param("nymID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid NymID"})
			return
		}
		ok := middleware.IsWalletContributor(c, *pk)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "insufficient priviledge"})
			return
		}

		hn(c, *pk)
	}
}

func checkViewerNymID(
	hn func(*gin.Context, identity.PublicKey),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		pk, err := identity.PublicKeyFromString(c.Param("nymID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid NymID"})
			return
		}
		ok := middleware.IsWalletViewer(c, *pk)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "insufficient priviledge"})
			return
		}

		hn(c, *pk)
	}
}
