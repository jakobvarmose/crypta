package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/ipfs/go-ipfs/core"

	"github.com/jakobvarmose/crypta/transaction"
	"github.com/jakobvarmose/crypta/userstore"
	"github.com/skratchdot/open-golang/open"
)

func securityHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		fmt.Println(req.RequestURI)
		resp.Header().Set("X-Frame-Options", "DENY")
		resp.Header().Set("X-Content-Type-Options", "nosniff")
		handler.ServeHTTP(resp, req)
	})
}

func New(n *core.IpfsNode, us *userstore.Userstore, db transaction.Database) (http.Handler, error) {
	app := http.NewServeMux()
	ws, err := NewWsServer(us, db)
	if err != nil {
		return nil, err
	}
	api, err := NewApiServer(n, us, db)
	if err != nil {
		return nil, err
	}
	gui, err := NewGuiServer(n)
	if err != nil {
		return nil, err
	}
	app.Handle("/api/ws", securityHeaders(ws))
	app.Handle("/api/", securityHeaders(api))
	app.Handle("/", securityHeaders(gui))
	return app, nil
}

func Run(app http.Handler, dev bool) error {
	if dev {
		listener, err := net.Listen("tcp", "localhost:8701")
		if err != nil {
			return err
		}
		return http.Serve(listener, app)
	} else {
		listener, err := net.Listen("tcp", "localhost:8700")
		if err != nil {
			return err
		}
		err = open.Start("http://localhost:8700")
		if err != nil {
			fmt.Println(err)
		}
		return http.Serve(listener, app)
	}
	return nil
}
