package model

import (
	"net/url"
)

type ScanRequest struct {
	Url    url.URL
	Id     string
	Type   string
	Status string
	Result ScanResult
}

func NewScanRequest(url url.URL, id, scanType string) ScanRequest {
	req := ScanRequest{
		Id:     id,
		Url:    url,
		Type:   scanType,
		Status: "queue",
	}
	return req
}

type ScanResult struct {
	Url        string `json:"url"`
	Id         string `json:"id"`
	ResultCode int    `json:"code"`
	Message    string `json:"message"`
}
