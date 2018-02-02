package server

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"

	"gx/ipfs/QmViBzgruNUoLNBnXcx8YWbDNwV8MNGEGKkLo6JGetygdw/go-ipfs/core"
	"gx/ipfs/QmViBzgruNUoLNBnXcx8YWbDNwV8MNGEGKkLo6JGetygdw/go-ipfs/core/coreunix"
)

func NewGuiServer(n *core.IpfsNode) (http.Handler, error) {
	app := http.NewServeMux()

	app.HandleFunc("/ipfs/", func(w http.ResponseWriter, req *http.Request) {
		// Find content type
		header := make([]byte, 512)
		r, err := coreunix.Cat(req.Context(), n, req.RequestURI)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		m, err := io.ReadFull(r, header)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mime := http.DetectContentType(header[:m])
		switch mime {
		case "text/plain; charset=utf-8",
			"image/jpeg", "image/png", "image/gif",
			"video/avi", "video/mp4",
			"audio/mpeg":
			w.Header().Set("Content-Type", mime)
			w.Header().Set("Content-Disposition", "inline")
		default:
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Content-Disposition", "attachment")
		}

		// Write content
		_, err = w.Write(header[:m])
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = io.Copy(w, r)
		if err != nil {
			fmt.Println(err)
			return
		}
	})

	app.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		uri := req.URL.RequestURI()
		_, err := AssetInfo(uri[1:])
		if err != nil {
			uri = "/index.html"
		}
		buf := MustAsset(uri[1:])
		contentType := mime.TypeByExtension(path.Ext(uri))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		resp.Header().Set("Content-Type", contentType)
		resp.Write(buf)
	})

	return app, nil
}
