package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/0x5d/hash/core"
	"go.uber.org/zap"
)

const (
	bodyLimit = 3000 // bytes
	slash     = "/"
)

type urlRouter struct {
	advertisedAddr string
	urlSvc         *core.URLService
	log            *zap.Logger
}

type urlBody struct {
	URL     string `json:"url"`
	Enabled bool   `json:"enabled"`
}

type urlResponse struct {
	Shortened string `json:"shortened"`
	Enabled   bool   `json:"enabled"`
}

func (r *urlRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.handleGet(res, req)
	case http.MethodPost:
		r.handlePost(res, req)
	case http.MethodPut:
		r.handlePut(res, req)
	default:
		http.Error(res, fmt.Sprintf("%q isn't allowed on URL resources", req.Method), http.StatusMethodNotAllowed)
	}
}

func (r *urlRouter) handleGet(res http.ResponseWriter, req *http.Request) {
	idStr, _ := shiftPath(req.URL.Path)
	_, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(res, "URL ID must be an integer", http.StatusBadRequest)
		return
	}
}

func (r *urlRouter) handlePost(res http.ResponseWriter, req *http.Request) {
	var u urlBody
	body := http.MaxBytesReader(res, req.Body, bodyLimit)
	err := parseJSON(r.log, res, body, &u)
	if err != nil {
		return
	}

	shortened, err := r.urlSvc.ShortenAndSave(req.Context(), u.URL, u.Enabled)
	if err != nil {
		errMsg := "Failed to create shortened URL"
		r.log.Error(errMsg, zap.Error(err), zap.String("url", u.URL))
		writeErrRes(res, errMsg, http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusCreated)
	bs, err := json.Marshal(shortenedURLResponse(shortened, r.advertisedAddr))
	if err != nil {
		errMsg := "Failed to encode response"
		r.log.Error(errMsg, zap.Error(err), zap.String("url", u.URL))
		writeErrRes(res, errMsg, http.StatusInternalServerError)
		return
	}
	res.Write(bs)
}

func (r *urlRouter) handlePut(res http.ResponseWriter, req *http.Request) {
	idStr, _ := shiftPath(req.URL.Path)
	_, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(res, "URL ID must be an integer", http.StatusBadRequest)
		return
	}
}

func shortenedURLResponse(url *core.ShortenedURL, advAddr string) *urlResponse {
	return &urlResponse{Shortened: fmt.Sprintf("%s/%s", advAddr, url.ShortKey), Enabled: url.Enabled}
}
