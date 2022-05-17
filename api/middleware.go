package api

import (
	"net/http"

	"git.digitus.me/library/prosper-kit/middleware"
	"git.digitus.me/prosperus/protocol/identity"
	"github.com/gin-gonic/gin"
)

func checkContributorNymID(
	permissionLevel string,
	hn func(*gin.Context, identity.PublicKey),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		pk, err := identity.PublicKeyFromString(c.Param("nymID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid NymID"})
			return
		}
		permission, ok := middleware.GetWalletPermission(c, *pk)
		if !ok ||
			permission != permissionLevel {
			c.JSON(http.StatusBadRequest, gin.H{"message": "insufficient priviledge"})
			return
		}

		hn(c, *pk)
	}
}

func checkViewerNymID(
	permissionLevel string,
	hn func(*gin.Context, identity.PublicKey),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		pk, err := identity.PublicKeyFromString(c.Param("nymID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid NymID"})
			return
		}
		_, ok := middleware.GetWalletPermission(c, *pk)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "insufficient priviledge"})
			return
		}

		hn(c, *pk)
	}
}
