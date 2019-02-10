package helper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
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
	gitUrl := req.Url.String()
	if !strings.HasSuffix(gitUrl, "/") {
		gitUrl = gitUrl + "/"
	}
	resp, err := http.Get(fmt.Sprintf("%s.git", gitUrl))
	if err != nil {
		log.Error(err)
		return req, err
	}
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
			return req, err
		}
		if strings.ContainsAny(string(body), "Index of") {
			result := model.NewScanResult(req.Url, req.Id, 1, ".git folder found")
			req.Result = result
			return req, err
		}
	}
	resp, err = http.Get(fmt.Sprintf("%s.git/", gitUrl))
	if err != nil {
		log.Error(err)
		return req, err
	}
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
			return req, err
		}
		if strings.ContainsAny(string(body), "Index of") {
			result := model.NewScanResult(req.Url, req.Id, 1, ".git folder found")
			req.Result = result
			return req, err
		}
	}

	resp, err = http.Get(fmt.Sprintf("%s.git/HEAD", gitUrl))
	if err != nil {
		log.Error(err)
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
