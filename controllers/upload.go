package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	cid "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-cid"
)

func Upload(p Provider, w http.ResponseWriter, req *http.Request) {
	header := make([]byte, 512)
	m, err := io.ReadFull(req.Body, header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	hash, err := p.DB().Put(io.MultiReader(bytes.NewReader(header[:m]), req.Body))
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
}
