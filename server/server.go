package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/ipsn/go-ipfs/core"

	"github.com/jakobvarmose/crypta/transaction"
	"github.com/jakobvarmose/crypta/userstore"
	"github.com/skratchdot/open-golang/open"
)

func New(n *core.IpfsNode, us *userstore.Userstore, db transaction.Database) (http.Handler, error) {
	app := http.NewServeMux()
	api, err := NewApiServer(n, us, db)
	if err != nil {
		return nil, err
	}
	gui, err := NewGuiServer(n)
	if err != nil {
		return nil, err
	}
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
