package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/0x5d/hash/core"
)

const (
	// Firefox allows URLs up to 64k chars, so 2x that should be enough.
	bodyLimit = 128000
	slash     = "/"
)

type errResponse struct {
	Msg string `json:"msg"`
}

type urlRouter struct {
	advertisedAddr string
	urlSvc         *core.URLService
	log            *zap.Logger
}

type urlBody struct {
	URL string `json:"url"`
}

func (b *urlBody) toCoreURL() core.ShortenedURL {
	return core.ShortenedURL{URL: b.URL}
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
	parseJSON(body, &u, res)

	shortened, err := r.urlSvc.ShortenAndSave(req.Context(), u.URL, u.Enabled)
	if err != nil {
		writeErrRes(res, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusCreated)
	bs, err := json.Marshal(&shortened)
	if err != nil {
		writeErrRes(res, "Failed to encode response", http.StatusInternalServerError)
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

func shiftPath(p string) (head, tail string) {
	p = path.Clean(slash + p)
	i := strings.Index(p[1:], slash) + 1
	if i <= 0 {
		return p[1:], slash
	}
	return p[1:i], p[i:]
}

func parseJSON(body io.ReadCloser, v any, res http.ResponseWriter) {
	err := json.NewDecoder(body).Decode(&v)
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, &http.MaxBytesError{}):
		writeErrRes(res, fmt.Sprintf("Request body should be under %dB", bodyLimit), http.StatusRequestEntityTooLarge)
	default:
		writeErrRes(res, "Invalid JSON", http.StatusBadRequest)
	}
}

func writeErrRes(res http.ResponseWriter, msg string, status int) {
	errRes := errResponse{Msg: msg}
	r, err := json.Marshal(&errRes)
	if err != nil {
		http.Error(res, "Failed to encode response.", http.StatusInternalServerError)
		return
	}
	http.Error(res, string(r), status)
}
