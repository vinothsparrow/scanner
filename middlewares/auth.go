package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vinothsparrow/scanner/config"
	h "github.com/vinothsparrow/scanner/helper"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.SetErrors(c)
		config := config.GetConfig()
		reqKey := c.Request.Header.Get("X-Auth-Key")

		var key string
		if key = config.GetString("auth.key"); len(strings.TrimSpace(key)) == 0 {
			h.AddError(c, h.ErrForbidden)
		} else if key != reqKey {
			reqKey = c.Query("api_key")
			if key != reqKey {
				h.AddError(c, h.ErrForbidden)
			}
		}
		c.Next()
	}
}
