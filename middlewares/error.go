package middlewares

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	h "github.com/vinothsparrow/scanner/helper"
	model "github.com/vinothsparrow/scanner/model"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		errors := h.GetErrors(c)
		if errors.HasError() {
			c.JSON(errors.StatusCode(), errors)
			c.AbortWithStatus(errors.StatusCode())
			return
		}
		c.Next()
		errors = h.GetErrors(c)
		if errors.HasError() {
			c.JSON(errors.StatusCode(), errors)
			c.AbortWithStatus(errors.StatusCode())
			return
		}
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("panic: %+v", err)
				errors := &model.Errors{Errors: []*model.Error{h.ErrInternalServer}}
				c.JSON(errors.StatusCode(), errors)
				c.AbortWithStatus(errors.StatusCode())
				return
			}
		}()
		c.Next()
	}
}
