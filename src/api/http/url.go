package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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

type urlUpdateBody struct {
	URL     *string `json:"url"`
	Enabled *bool   `json:"enabled"`
}

type urlResponse struct {
	Shortened string `json:"shortened"`
	Enabled   bool   `json:"enabled"`
}

func (r *urlRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	_, req.URL.Path = shiftPath(req.URL.Path)
	switch req.Method {
	case http.MethodGet:
		r.handleGet(res, req)
	case http.MethodPost:
		r.handlePost(res, req)
	case http.MethodPut:
		r.handlePut(res, req)
	default:
		writeErrRes(res, fmt.Sprintf("%q not allowed", req.Method), http.StatusMethodNotAllowed)
	}
}

func (r *urlRouter) handleGet(res http.ResponseWriter, req *http.Request) {
	key, _ := shiftPath(req.URL.Path)
	shortened, err := r.urlSvc.DecodeAndGet(req.Context(), key)
	if errors.Is(err, &core.ErrNotFound{}) {
		writeErrRes(res, fmt.Sprintf("URL with key %s not found", key), http.StatusNotFound)
		return
	}
	if err != nil {
		errMsg := "Failed to get URL"
		r.log.Error(errMsg, zap.Error(err))
		writeErrRes(res, errMsg, http.StatusInternalServerError)
		return
	}
	bs, err := json.Marshal(shortenedURLResponse(shortened, r.advertisedAddr))
	if err != nil {
		errMsg := "Failed to encode response"
		r.log.Error(errMsg, zap.Error(err))
		writeErrRes(res, errMsg, http.StatusInternalServerError)
		return
	}
	res.Write(bs)
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
	bs, err := json.Marshal(shortenedURLResponse(shortened, r.advertisedAddr))
	if err != nil {
		errMsg := "Failed to encode response"
		r.log.Error(errMsg, zap.Error(err), zap.String("url", u.URL))
		writeErrRes(res, errMsg, http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Write(bs)
}

func (r *urlRouter) handlePut(res http.ResponseWriter, req *http.Request) {
	key, _ := shiftPath(req.URL.Path)
	body := http.MaxBytesReader(res, req.Body, bodyLimit)
	var update urlUpdateBody
	err := parseJSON(r.log, res, body, &update)
	if err != nil {
		return
	}
	shortened, err := r.urlSvc.Update(req.Context(), key, update.URL, update.Enabled)
	if err != nil {
		errMsg := "Failed to update URL"
		r.log.Error(errMsg, zap.Error(err))
		writeErrRes(res, errMsg, http.StatusInternalServerError)
		return
	}
	bs, err := json.Marshal(shortenedURLResponse(shortened, r.advertisedAddr))
	if err != nil {
		errMsg := "Failed to encode response"
		r.log.Error(errMsg, zap.Error(err))
		writeErrRes(res, errMsg, http.StatusInternalServerError)
		return
	}
	res.Write(bs)
}

func shortenedURLResponse(url *core.ShortenedURL, advAddr string) *urlResponse {
	return &urlResponse{Shortened: fmt.Sprintf("%s/%s", advAddr, url.ShortKey), Enabled: url.Enabled}
}
