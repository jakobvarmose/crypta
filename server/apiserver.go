package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	cid "github.com/ipfs/go-cid"

	"github.com/ipfs/go-ipfs/core"

	"github.com/jakobvarmose/crypta/transaction"
	"github.com/jakobvarmose/crypta/userstore"
)

func securityCheck(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		originString := req.Header.Get("Origin")
		referrerString := req.Header.Get("Referer")
		if originString == "" && referrerString == "" {
			http.Error(resp, "403 forbidden - no origin or referrer", 403)
			return
		}
		if originString != "" {
			origin, err := url.ParseRequestURI(originString)
			if err != nil {
				http.Error(resp, "403 forbidden - invalid origin", 403)
				return
			}
			if origin.Host != req.Host {
				http.Error(resp, "403 forbidden - wrong origin", 403)
				return
			}
		}
		if referrerString != "" {
			referrer, err := url.ParseRequestURI(referrerString)
			if err != nil {
				http.Error(resp, "403 forbidden - invalid referrer", 403)
				return
			}
			if referrer.Host != req.Host {
				http.Error(resp, "403 forbidden - wrong referrer", 403)
				return
			}
		}
		handler.ServeHTTP(resp, req)
	})
}

func convert(val interface{}) interface{} {
	switch val := val.(type) {
	case map[interface{}]interface{}:
		res := make(map[string]interface{}, len(val))
		for k := range val {
			if k, ok := k.(string); ok {
				res[k] = convert(val[k])
			}
		}
		return res
	case map[string]interface{}:
		res := make(map[string]interface{}, len(val))
		for k := range val {
			res[k] = convert(val[k])
		}
		return res
	case []interface{}:
		res := make([]interface{}, len(val))
		for i := range val {
			res[i] = convert(val[i])
		}
		return res
	default:
		return val
	}
}

func NewApiServer(n *core.IpfsNode, us *userstore.Userstore, db transaction.Database) (http.Handler, error) {
	app := http.NewServeMux()
	app.Handle("/api/v0/upload", securityCheck(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		header := make([]byte, 512)
		m, err := io.ReadFull(req.Body, header)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hash, err := db.Put(io.MultiReader(bytes.NewReader(header[:m]), req.Body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		c, err := cid.Parse(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mime := http.DetectContentType(header[:m])
		if strings.HasPrefix(mime, "image/") {
			mime = "IMAGE"
		} else if strings.HasPrefix(mime, "video/") {
			mime = "VIDEO"
		} else {
			http.Error(w, "Unknown filetype", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)

		enc.Encode(map[string]interface{}{
			"t":    mime,
			"hash": c.String(),
		})
	})))
	return app, nil
}
