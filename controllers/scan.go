package controllers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vinothsparrow/scanner/helper"
	"github.com/vinothsparrow/scanner/model"
)

type ScanController struct{}

func (u ScanController) GitScan(c *gin.Context) {
	urlParam := strings.TrimSpace(c.Query("url"))
	if len(urlParam) != 0 {
		u, err := url.Parse(urlParam)
		if err != nil {
			c.JSON(500, gin.H{"message": "Not a valid URL"})
			c.Abort()
			return
		}
		scheme := strings.TrimSpace(u.Scheme)
		if len(scheme) == 0 {
			u, err = url.Parse(fmt.Sprintf("http://%s", urlParam))
			if err != nil {
				c.JSON(500, gin.H{"message": "Not a valid URL"})
				c.Abort()
				return
			}
		}
		uuidObj, _ := uuid.NewRandom()
		id := uuidObj.String()
		req := model.NewScanRequest(*u, id, "git")
		helper.WorkQueue <- req
		c.JSON(200, gin.H{"message": "Scan queued", "Id": id})
		return
	}
	c.JSON(500, gin.H{"message": "Not a valid URL"})
	c.Abort()
	return
}

func (u ScanController) Status(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(500, gin.H{"message": "Not a valid scan ID"})
		c.Abort()
		return
	}
	req, err := helper.GetScanRequest(id)
	if err != nil {
		helper.AddError(c, helper.ErrInternalServer)
		return
	}
	c.JSON(200, gin.H{"Id": req.Id, "Status": req.Status, "Result": req.Result})
	return
}
