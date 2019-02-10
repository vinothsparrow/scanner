package controllers

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	h "github.com/vinothsparrow/scanner/helper"
)

type NotFoundController struct{}

func (n NotFoundController) Error404(c *gin.Context) {
	log.Info("404")
	h.AddError(c, h.ErrNotFound)
	c.Next()
}

func (n NotFoundController) Error405(c *gin.Context) {
	log.Info("405")
	h.AddError(c, h.ErrMethodNotFound)
	c.Next()
}
