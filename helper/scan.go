package helper

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/vinothsparrow/scanner/model"
)

type html struct {
	Head title `xml:"head"`
}

type title struct {
	Title string `xml:"title"`
}

func GitScan(req *model.ScanRequest) (*model.ScanRequest, error) {
	result := model.NewScanResult(req.Url, req.Id, -1, "error occured")
	req.Result = result
	resp, err := http.Get(fmt.Sprintf("%s/.git", req.Url.String()))
	if err != nil {
		return req, err
	}
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return req, err
		}
		h := html{}
		err = xml.NewDecoder(bytes.NewBuffer(body)).Decode(&h)
		if err != nil {
			return req, err
		}
		if len(strings.TrimSpace(h.Head.Title)) != 0 && strings.ContainsAny(h.Head.Title, "Index of") {
			result := model.NewScanResult(req.Url, req.Id, 1, ".git folder found")
			req.Result = result
			return req, err
		}
	}
	resp, err = http.Get(fmt.Sprintf("%s/.git/", req.Url.String()))
	if err != nil {
		return req, err
	}
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return req, err
		}
		h := html{}
		err = xml.NewDecoder(bytes.NewBuffer(body)).Decode(&h)
		if err != nil {
			return req, err
		}
		if len(strings.TrimSpace(h.Head.Title)) != 0 && strings.ContainsAny(h.Head.Title, "Index of") {
			result := model.NewScanResult(req.Url, req.Id, 1, ".git folder found")
			req.Result = result
			return req, err
		}
	}
	resp, err = http.Get(fmt.Sprintf("%s/.git/HEAD", req.Url.String()))
	if err != nil {
		return req, err
	}
	if resp.StatusCode == http.StatusOK {
		result := model.NewScanResult(req.Url, req.Id, 1, ".git folder found")
		req.Result = result
		return req, err
	}
	result = model.NewScanResult(req.Url, req.Id, 0, ".git folder not found")
	req.Result = result
	return req, nil
}
