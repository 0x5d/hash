package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/0x5d/hash/core"
	"go.uber.org/zap"
)

type resolver struct {
	urlSvc *core.URLService
	log    *zap.Logger
}

func (r *resolver) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.handleGet(res, req)
	default:
		writeErrRes(res, fmt.Sprintf("%q not allowed", req.Method), http.StatusMethodNotAllowed)
	}
}

func (r *resolver) handleGet(res http.ResponseWriter, req *http.Request) {
	key, _ := shiftPath(req.URL.Path)

	url, err := r.urlSvc.DecodeAndGet(req.Context(), key)
	if err != nil {
		errMsg := "Failed to get shortened URL"
		r.log.Error(errMsg, zap.Error(err), zap.String("key", key))
		writeErrRes(res, errMsg, http.StatusInternalServerError)
		return
	}
	if !url.Enabled {
		writeErrRes(res, "Not found", http.StatusNotFound)
		return
	}

	target := url.Original
	if !strings.HasPrefix(target, "https://") && !strings.HasPrefix(target, "http://") {
		target = "http://" + target
	}
	http.Redirect(res, req, target, http.StatusTemporaryRedirect)
}
