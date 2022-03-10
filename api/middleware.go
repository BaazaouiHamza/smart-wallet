package api

import (
	"net/http"

	"git.digitus.me/library/prosper-kit/middleware"
	"git.digitus.me/prosperus/protocol/identity"
	"github.com/gin-gonic/gin"
)

func checkNymID(
	permissionLevel string,
	hn func(*gin.Context, identity.PublicKey),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		pk, err := identity.PublicKeyFromString(c.Param("nym-id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid NymID"})
			return
		}

		permission, ok := middleware.GetWalletPermission(c, *pk)
		if !ok ||
			permissionLevel == middleware.ViewerPermissionLevel ||
			permission != permissionLevel {
			c.JSON(http.StatusBadRequest, gin.H{"message": "insufficient priviledge"})
			return
		}

		hn(c, *pk)
	}
}
