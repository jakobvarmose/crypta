package server

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jakobvarmose/crypta/server/rpc"
	"github.com/jakobvarmose/crypta/transaction"
	"github.com/jakobvarmose/crypta/userstore"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewWsServer(us *userstore.Userstore, db transaction.Database) (http.Handler, error) {
	handler := securityCheck(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		state := &State{
			US: us,
			DB: db,
		}
		state.RPC = rpc.New(conn, state, methods)
	}))
	return handler, nil
}
