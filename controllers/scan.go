package controllers

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vinothsparrow/scanner/helper"
	"github.com/vinothsparrow/scanner/model"
)

type ScanController struct{}

// GitScanController godoc
// @Summary GitScan api
// @Description submit url for scan
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param url query string true "url to scan"
// @Success 200 {object} model.ScanResultResponse
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.ScanResultError
// @Router /v1/scan/git [get]
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
		if checUrlInternal(u) {
			c.JSON(500, gin.H{"message": "Not a valid URL"})
			c.Abort()
			return
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

// ScanStatusController godoc
// @Summary Scan status api
// @Description get the status scan
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param id path string true "scan ID"
// @Success 200 {object} model.ScanStatusResultResponse
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Failure 500 {object} model.ScanResultError
// @Router /v1/scan/status/{id} [get]
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

// Check if the host in internal to avoid SSRF
func checUrlInternal(url *url.URL) bool {
	private := false
	scheme := url.Scheme
	host := url.Hostname()
	if scheme != "http" && scheme != "https" {
		return true
	}
	if host == "localhost" || host == "0.0.0.0" || host == "metadata.google.internal" {
		return true
	}
	IPS, _ := net.LookupIP(host)
	for _, IP := range IPS {
		if IP.String() == "169.254.169.254" {
			private = true
		}
		if !private && IP != nil {
			_, privateLoopBitBlock, _ := net.ParseCIDR("127.0.0.0/8")
			_, private24BitBlock, _ := net.ParseCIDR("10.0.0.0/8")
			_, private20BitBlock, _ := net.ParseCIDR("172.16.0.0/12")
			_, private16BitBlock, _ := net.ParseCIDR("192.168.0.0/16")
			_, privateIPV6LoopBlock, _ := net.ParseCIDR("::1/128")
			_, privateIPV6LinkBlock, _ := net.ParseCIDR("fe80::/10")
			_, privateIPV6Block, _ := net.ParseCIDR("fd00::/8")
			private = privateLoopBitBlock.Contains(IP) || private24BitBlock.Contains(IP) || private20BitBlock.Contains(IP) || private16BitBlock.Contains(IP)
			private = private || privateIPV6LoopBlock.Contains(IP) || privateIPV6LinkBlock.Contains(IP) || privateIPV6Block.Contains(IP)
		}
		if private {
			return private
		}
	}
	return private
}
